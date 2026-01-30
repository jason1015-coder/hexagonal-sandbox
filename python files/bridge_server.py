"""
TCP Bridge Server
Runs the Python game logic and communicates with Rust Bevy engine via TCP.
"""

import socket
import json
import threading
import pickle
import os
from main import Game, World

class BridgeServer:
    def __init__(self, host='127.0.0.1', port=9999):
        self.host = host
        self.port = port
        self.game = Game()
        self.running = False
        self.clients = []
        
    def start(self):
        self.running = True
        server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        server_socket.bind((self.host, self.port))
        server_socket.listen(5)
        
        print(f"Bridge server listening on {self.host}:{self.port}")
        
        while self.running:
            try:
                client_socket, address = server_socket.accept()
                print(f"Client connected: {address}")
                self.clients.append(client_socket)
                
                client_thread = threading.Thread(
                    target=self.handle_client,
                    args=(client_socket,)
                )
                client_thread.daemon = True
                client_thread.start()
            except Exception as e:
                print(f"Error accepting connection: {e}")
                break
        
        server_socket.close()
    
    def handle_client(self, client_socket):
        while self.running:
            try:
                data = client_socket.recv(4096)
                if not data:
                    break
                
                message = json.loads(data.decode('utf-8'))
                response = self.process_message(message)
                
                if response:
                    client_socket.sendall((json.dumps(response) + "\n").encode('utf-8'))
            except Exception as e:
                print(f"Error handling client: {e}")
                break
        
        client_socket.close()
        if client_socket in self.clients:
            self.clients.remove(client_socket)
    
    def process_message(self, message):
        msg_type = message.get('type')
        
        if msg_type == 'start_game':
            mode = message.get('mode', 'survival')
            world_name = message.get('world_name', 'world')
            self.game.start_game(mode, world_name)
            return {'type': 'start_game', 'status': 'success'}
        
        elif msg_type == 'update':
            return self.handle_update(message)
        
        elif msg_type == 'get_state':
            return self.get_game_state()
        
        elif msg_type == 'jump':
            if self.game.player:
                self.game.player.jump()
            return {'type': 'jump', 'status': 'success'}
        
        elif msg_type == 'select_slot':
            slot = message.get('slot', 0)
            if self.game.player:
                self.game.player.selected_slot = slot
            return {'type': 'select_slot', 'status': 'success'}
        
        elif msg_type == 'save':
            self.game.save_game()
            return {'type': 'save', 'status': 'success'}
        
        else:
            return {'type': 'error', 'message': f'Unknown message type: {msg_type}'}
    
    def handle_update(self, message):
        if self.game.state != "GAME":
            return {'type': 'update', 'status': 'not_in_game', 'state': self.game.state}
        
        keys = message.get('keys', {})
        mining = message.get('mining', False)
        place_block = message.get('place_block', False)
        mouse_pos = message.get('mouse_position', [0, 0])
        
        # Map key names to pygame key codes
        key_map = {
            'w': 119, 'a': 97, 's': 115, 'd': 100,
            'space': 32, 'escape': 27,
            '1': 49, '2': 50, '3': 51, '4': 52, '5': 53,
            '6': 54, '7': 55, '8': 56, '9': 57,
            'equals': 61, 'minus': 45, 'f': 102
        }
        
        pygame_keys = {}
        for key, pressed in keys.items():
            if key in key_map:
                pygame_keys[key_map[key]] = pressed
        
        class MockKeys:
            def __getitem__(self, key):
                return pygame_keys.get(key, False)
        
        # Update player
        nearby_hexagons = self.game.world.get_nearby_hexagons(
            self.game.player.x, self.game.player.y, radius=500
        )
        self.game.player.update(MockKeys(), nearby_hexagons)
        
        # Handle mining and placing
        self.game.mining = mining
        self.game.camera_x = self.game.player.x - 640
        self.game.camera_y = self.game.player.y - 360
        
        if mining:
            world_mx = mouse_pos[0] + self.game.camera_x
            world_my = mouse_pos[1] + self.game.camera_y
            self.game.mine_block()
        
        if place_block:
            world_mx = mouse_pos[0] + self.game.camera_x
            world_my = mouse_pos[1] + self.game.camera_y
            self.game.place_block()
        
        # Update particles
        self.game.particles = [p for p in self.game.particles if not (p.update() or p.is_dead())]
        
        # Unload distant chunks
        self.game.world.unload_distant_chunks(self.game.player.x, self.game.player.y)
        
        return {'type': 'update', 'status': 'success'}
    
    def get_game_state(self):
        if self.game.state != "GAME" or self.game.player is None:
            return {
                'type': 'state',
                'in_game': False,
                'state': self.game.state
            }
        
        # Get visible hexagons
        visible_hexagons = []
        player_chunk_x, player_chunk_y = self.game.world.get_chunk_coords(
            self.game.player.x, self.game.player.y
        )
        
        for chunk_x in range(player_chunk_x - self.game.render_distance,
                             player_chunk_x + self.game.render_distance + 1):
            for chunk_y in range(player_chunk_y - self.game.render_distance,
                                 player_chunk_y + self.game.render_distance + 1):
                if (chunk_x, chunk_y) in self.game.world.chunks:
                    chunk = self.game.world.chunks[(chunk_x, chunk_y)]
                    for hexagon in chunk.hexagons.values():
                        screen_x = hexagon.x - self.game.camera_x
                        screen_y = hexagon.y - self.game.camera_y
                        
                        if -30 < screen_x < 1310 and -30 < screen_y < 750:
                            visible_hexagons.append({
                                'x': float(hexagon.x),
                                'y': float(hexagon.y),
                                'block_type': hexagon.block_type,
                                'color': list(hexagon.color),
                                'health': float(hexagon.health),
                                'max_health': float(hexagon.max_health),
                                'transparent': hexagon.transparent
                            })
        
        # Get visible organisms
        visible_organisms = []
        for organism in self.game.world.get_nearby_organisms(
            self.game.player.x, self.game.player.y, radius=600
        ):
            screen_x = organism.x - self.game.camera_x
            screen_y = organism.y - self.game.camera_y
            if -100 < screen_x < 1380 and -100 < screen_y < 820:
                visible_organisms.append({
                    'x': float(organism.x),
                    'y': float(organism.y),
                    'type': organism.organism_type,
                    'health': float(organism.health),
                    'max_health': float(organism.max_health)
                })
        
        # Get particles
        particles = []
        for particle in self.game.particles:
            px = particle.x - self.game.camera_x
            py = particle.y - self.game.camera_y
            if -10 < px < 1290 and -10 < py < 730:
                particles.append({
                    'x': float(particle.x),
                    'y': float(particle.y),
                    'color': list(particle.color),
                    'age': float(particle.age),
                    'lifetime': float(particle.lifetime),
                    'size': float(particle.size)
                })
        
        # Get player state
        player_state = {
            'x': float(self.game.player.x),
            'y': float(self.game.player.y),
            'radius': float(self.game.player.radius),
            'color': list(self.game.player.color),
            'selected_slot': self.game.player.selected_slot,
            'flight_mode': self.game.player.flight_mode,
            'on_ground': self.game.player.on_ground
        }
        
        # Get camera
        camera = {
            'x': float(self.game.camera_x),
            'y': float(self.game.camera_y)
        }
        
        # Get inventory
        from item import DEFAULT_HOTBAR_ITEMS
        inventory = {}
        for item_id in DEFAULT_HOTBAR_ITEMS[:9]:
            inventory[item_id] = self.game.player.inventory.get(item_id, 0)
        
        return {
            'type': 'state',
            'in_game': True,
            'player': player_state,
            'camera': camera,
            'hexagons': visible_hexagons,
            'organisms': visible_organisms,
            'particles': particles,
            'inventory': inventory,
            'game_mode': self.game.game_mode,
            'render_distance': self.game.render_distance
        }

if __name__ == '__main__':
    server = BridgeServer()
    try:
        server.start()
    except KeyboardInterrupt:
        print("\nShutting down server...")
        server.running = False