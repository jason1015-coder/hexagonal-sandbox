"""
Python Game API Server
Provides HTTP endpoints for the Rust Bevy engine to access game state.
"""

from flask import Flask, jsonify, request
from flask_cors import CORS
import pickle
import os
import threading
import time

# Import game modules
from main import Game, World, Player
from blocks import COLOR_BY_TYPE, BLOCK_DEFINITIONS
from item import ITEM_DEFINITIONS, ITEM_NAME_BY_ID, ITEM_ICON_COLOR_BY_ID, DEFAULT_HOTBAR_ITEMS
from biomes import BIOME_CHARACTERISTICS
from organism import ORGANISM_DEFINITIONS

app = Flask(__name__)
CORS(app)

# Global game instance
game_instance = None
game_lock = threading.Lock()

def init_game():
    """Initialize the game instance."""
    global game_instance
    with game_lock:
        if game_instance is None:
            game_instance = Game()
        return game_instance

@app.route('/api/game/start', methods=['POST'])
def start_game():
    """Start a new game or load existing game."""
    data = request.json
    mode = data.get('mode', 'survival')
    world_name = data.get('world_name', 'world')
    
    game = init_game()
    game.start_game(mode, world_name)
    
    return jsonify({
        'status': 'success',
        'mode': mode,
        'world_name': world_name
    })

@app.route('/api/game/update', methods=['POST'])
def update_game():
    """Update game state with input from Rust engine."""
    game = init_game()
    
    if game.state != "GAME":
        return jsonify({
            'status': 'not_in_game',
            'state': game.state
        })
    
    data = request.json
    
    # Process keys
    keys = data.get('keys', {})
    pygame_keys = {}
    key_map = {
        'w': 119, 'a': 97, 's': 115, 'd': 100,
        'space': 32, 'escape': 27,
        '1': 49, '2': 50, '3': 51, '4': 52, '5': 53,
        '6': 54, '7': 55, '8': 56, '9': 57,
        'equals': 61, 'minus': 45, 'f': 102
    }
    
    for key, pressed in keys.items():
        if key in key_map:
            pygame_keys[key_map[key]] = pressed
    
    # Create mock pygame key state
    class MockKeys:
        def __getitem__(self, key):
            return pygame_keys.get(key, False)
    
    # Update player
    nearby_hexagons = game.world.get_nearby_hexagons(
        game.player.x, game.player.y, radius=500
    )
    game.player.update(MockKeys(), nearby_hexagons)
    
    # Handle mouse input
    if data.get('mining', False):
        game.mining = True
    else:
        game.mining = False
    
    if data.get('place_block', False):
        # Need mouse position
        mouse_pos = data.get('mouse_position', [0, 0])
        game.camera_x = game.player.x - 640  # SCREEN_WIDTH // 2
        game.camera_y = game.player.y - 360  # SCREEN_HEIGHT // 2
        world_mx = mouse_pos[0] + game.camera_x
        world_my = mouse_pos[1] + game.camera_y
        game.place_block()
    
    if game.mining:
        mouse_pos = data.get('mouse_position', [0, 0])
        game.camera_x = game.player.x - 640
        game.camera_y = game.player.y - 360
        game.mine_block()
    
    # Update particles
    game.particles = [p for p in game.particles if not (p.update() or p.is_dead())]
    
    # Unload distant chunks
    game.world.unload_distant_chunks(game.player.x, game.player.y)
    
    return jsonify({'status': 'success'})

@app.route('/api/game/state', methods=['GET'])
def get_game_state():
    """Get current game state for rendering."""
    game = init_game()
    
    if game.state != "GAME" or game.player is None:
        return jsonify({
            'state': game.state,
            'in_game': False
        })
    
    # Get visible chunks
    player_chunk_x, player_chunk_y = game.world.get_chunk_coords(
        game.player.x, game.player.y
    )
    
    visible_hexagons = []
    render_distance = game.render_distance
    
    for chunk_x in range(player_chunk_x - render_distance, 
                         player_chunk_x + render_distance + 1):
        for chunk_y in range(player_chunk_y - render_distance, 
                             player_chunk_y + render_distance + 1):
            if (chunk_x, chunk_y) in game.world.chunks:
                chunk = game.world.chunks[(chunk_x, chunk_y)]
                for hexagon in chunk.hexagons.values():
                    screen_x = hexagon.x - game.camera_x
                    screen_y = hexagon.y - game.camera_y
                    
                    if -30 < screen_x < 1310 and -30 < screen_y < 750:
                        visible_hexagons.append({
                            'x': hexagon.x,
                            'y': hexagon.y,
                            'block_type': hexagon.block_type,
                            'color': hexagon.color,
                            'health': hexagon.health,
                            'max_health': hexagon.max_health,
                            'transparent': hexagon.transparent
                        })
    
    # Get visible organisms
    visible_organisms = []
    for organism in game.world.get_nearby_organisms(
        game.player.x, game.player.y, radius=600
    ):
        screen_x = organism.x - game.camera_x
        screen_y = organism.y - game.camera_y
        if -100 < screen_x < 1380 and -100 < screen_y < 820:
            visible_organisms.append({
                'x': organism.x,
                'y': organism.y,
                'type': organism.organism_type,
                'health': organism.health,
                'max_health': organism.max_health
            })
    
    # Get particles
    particles = []
    for particle in game.particles:
        px = particle.x - game.camera_x
        py = particle.y - game.camera_y
        if -10 < px < 1290 and -10 < py < 730:
            particles.append({
                'x': particle.x,
                'y': particle.y,
                'color': particle.color,
                'age': particle.age,
                'lifetime': particle.lifetime,
                'size': particle.size
            })
    
    # Get player state
    player_state = {
        'x': game.player.x,
        'y': game.player.y,
        'radius': game.player.radius,
        'color': game.player.color,
        'selected_slot': game.player.selected_slot,
        'flight_mode': game.player.flight_mode,
        'on_ground': game.player.on_ground
    }
    
    # Get camera
    camera = {
        'x': game.camera_x,
        'y': game.camera_y
    }
    
    # Get inventory
    inventory = {}
    for item_id in DEFAULT_HOTBAR_ITEMS[:9]:
        inventory[item_id] = game.player.inventory.get(item_id, 0)
    
    return jsonify({
        'state': game.state,
        'in_game': True,
        'player': player_state,
        'camera': camera,
        'hexagons': visible_hexagons,
        'organisms': visible_organisms,
        'particles': particles,
        'inventory': inventory,
        'game_mode': game.game_mode,
        'render_distance': render_distance
    })

@app.route('/api/game/jump', methods=['POST'])
def jump():
    """Make player jump."""
    game = init_game()
    if game.player:
        game.player.jump()
    return jsonify({'status': 'success'})

@app.route('/api/game/select_slot', methods=['POST'])
def select_slot():
    """Select hotbar slot."""
    data = request.json
    slot = data.get('slot', 0)
    game = init_game()
    if game.player:
        game.player.selected_slot = slot
    return jsonify({'status': 'success'})

@app.route('/api/game/save', methods=['POST'])
def save_game():
    """Save game."""
    game = init_game()
    game.save_game()
    return jsonify({'status': 'success'})

@app.route('/api/game/pause', methods=['POST'])
def pause_game():
    """Toggle pause."""
    game = init_game()
    if game.state == "GAME":
        game.paused = True
        game.state = "PAUSE"
    elif game.state == "PAUSE":
        game.paused = False
        game.state = "GAME"
    return jsonify({'status': 'success', 'state': game.state})

@app.route('/api/game/worlds', methods=['GET'])
def get_worlds():
    """Get list of available worlds."""
    worlds = []
    if os.path.exists("saves"):
        for filename in os.listdir("saves"):
            if filename.endswith(".pkl") and not filename.startswith("player_"):
                world_name = filename.replace(".pkl", "")
                worlds.append(world_name)
    
    return jsonify({'worlds': worlds})

@app.route('/api/data/blocks', methods=['GET'])
def get_blocks():
    """Get block definitions."""
    return jsonify(BLOCK_DEFINITIONS)

@app.route('/api/data/items', methods=['GET'])
def get_items():
    """Get item definitions."""
    return jsonify(ITEM_DEFINITIONS)

@app.route('/api/data/hotbar', methods=['GET'])
def get_hotbar():
    """Get hotbar configuration."""
    return jsonify({
        'items': DEFAULT_HOTBAR_ITEMS[:9],
        'names': {k: ITEM_NAME_BY_ID[k] for k in DEFAULT_HOTBAR_ITEMS[:9]},
        'colors': {k: ITEM_ICON_COLOR_BY_ID[k] for k in DEFAULT_HOTBAR_ITEMS[:9]}
    })

@app.route('/api/data/biomes', methods=['GET'])
def get_biomes():
    """Get biome definitions."""
    return jsonify(BIOME_CHARACTERISTICS)

@app.route('/api/data/organisms', methods=['GET'])
def get_organisms():
    """Get organism definitions."""
    return jsonify(ORGANISM_DEFINITIONS)

if __name__ == '__main__':
    init_game()
    app.run(host='127.0.0.1', port=5000, debug=False, threaded=True)