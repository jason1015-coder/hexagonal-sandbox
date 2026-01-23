import pygame
import math
import sys
import os
import random
import warnings
import pickle
import numpy as np
import time
from blocks import (
    BLOCK_DEFINITIONS,
    COLOR_BY_TYPE,
    HARDNESS_BY_TYPE,
    TRANSPARENT_BY_TYPE,
    SOLID_BY_TYPE,          # now used below if you want
    COLLECTIBLE_BY_TYPE
)



warnings.filterwarnings("ignore", category=RuntimeWarning, module="importlib._bootstrap")
os.environ['SDL_AUDIODRIVER'] = 'dummy'

pygame.init()
pygame.mixer.init()

# Constants
SCREEN_WIDTH = 1280
SCREEN_HEIGHT = 720
FPS = 60

# Colors
WHITE = (255, 255, 255)
BLACK = (0, 0, 0)
BLUE= (50,150,255)
GRAY = (128, 128, 128)
DARK_GRAY = (64, 64, 64)
GREEN = (100, 200, 100)
RED = (200, 100, 100)


# Physics constants
GRAVITY = 0.5
PLAYER_SPEED = 5
JUMP_FORCE = -12
FRICTION = 0.8
MINING_RANGE = 150

# Hexagon constants
HEX_SIZE = 30
HEX_WIDTH = math.sqrt(3) * HEX_SIZE
HEX_HEIGHT = 2 * HEX_SIZE
HEX_V_SPACING = HEX_HEIGHT * 0.75

# Chunk constants
CHUNK_SIZE = 32
CHUNK_WIDTH = CHUNK_SIZE * HEX_WIDTH
CHUNK_HEIGHT = CHUNK_SIZE * HEX_V_SPACING
RENDER_DISTANCE = 4  # Number of chunks to render around player



class Particle:
    def __init__(self, x, y, color):
        self.x = x
        self.y = y
        self.vel_x = random.uniform(-2, 2)
        self.vel_y = random.uniform(-3, -1)
        self.color = color
        self.lifetime = 30
        self.age = 0
        self.size = random.randint(2, 4)
    
    def update(self):
        self.vel_y += 0.2
        self.x += self.vel_x
        self.y += self.vel_y
        self.age += 1
    
    def draw_offset(self, screen, offset_x, offset_y):
        alpha = 1 - (self.age / self.lifetime)
        if alpha > 0:
            color = tuple(int(c * alpha) for c in self.color[:3])
            pygame.draw.circle(screen, color, (int(self.x - offset_x), int(self.y - offset_y)), self.size)
    
    def is_dead(self):
        return self.age >= self.lifetime

class Hexagon:
    def __init__(self, x, y, size, block_type='dirt'):
        self.x = x
        self.y = y
        self.size = size
        self.block_type = block_type
        
        # Fallback if invalid block type
        if block_type not in COLOR_BY_TYPE:
            block_type = 'dirt'
            
        self.color = COLOR_BY_TYPE[block_type]
        self.active_color = self.color
        self.hovered = False
        self.health = 100
        self.max_health = 100
        self.transparent = TRANSPARENT_BY_TYPE[block_type]
        self.corners = self._calculate_corners()
        self.chunk_x = 0
        self.chunk_y = 0
        
    def _calculate_corners(self):
        corners = []
        for i in range(6):
            angle = math.radians(30 + 60 * i)
            px = self.x + self.size * math.cos(angle)
            py = self.y + self.size * math.sin(angle)
            corners.append((px, py))
        return corners
    
    def check_hover(self, mouse_x, mouse_y, player_x, player_y):
        dx = mouse_x - self.x
        dy = mouse_y - self.y
        distance_sq = dx * dx + dy * dy
        
        pdx = player_x - self.x
        pdy = player_y - self.y
        player_distance_sq = pdx * pdx + pdy * pdy
        in_range = player_distance_sq < MINING_RANGE * MINING_RANGE
        
        hex_radius_sq = (self.size * 0.866) ** 2
        self.hovered = distance_sq < hex_radius_sq and in_range
        
        if self.hovered:
            self.active_color = tuple(min(c + 30, 255) for c in self.color[:3])
        else:
            self.active_color = self.color
    
    def take_damage(self, amount):
        hardness = HARDNESS_BY_TYPE.get(self.block_type, 1.0)
        if hardness <= 0:
            return False
        self.health -= amount / hardness
        return self.health <= 0
    
    def get_top_surface_y(self, x):
        left = self.corners[3]
        top = self.corners[4]
        right = self.corners[5]
        
        if left[0] <= x <= top[0]:
            if left[0] != top[0]:
                t = (x - left[0]) / (top[0] - left[0])
                return left[1] + t * (top[1] - left[1])
            else:
                return left[1]
        elif top[0] <= x <= right[0]:
            if top[0] != right[0]:
                t = (x - top[0]) / (right[0] - top[0])
                return top[1] + t * (right[1] - top[1])
            else:
                return top[1]
        
        return None

class Chunk:
    def __init__(self, chunk_x, chunk_y):
        self.chunk_x = chunk_x
        self.chunk_y = chunk_y
        self.hexagons = {}
        self.modified = False
        self.last_accessed = time.time()
        
    def get_world_position(self):
        world_x = self.chunk_x * CHUNK_WIDTH
        world_y = self.chunk_y * CHUNK_HEIGHT
        return world_x, world_y
    
    def add_hexagon(self, x, y, hexagon):
        # Convert world position to local chunk coordinates
        local_x = int((x - self.chunk_x * CHUNK_WIDTH) // HEX_WIDTH)
        local_y = int((y - self.chunk_y * CHUNK_HEIGHT) // HEX_V_SPACING)
        hexagon.chunk_x = self.chunk_x
        hexagon.chunk_y = self.chunk_y
        self.hexagons[(local_x, local_y)] = hexagon
        self.modified = True
        
    def get_hexagon(self, x, y):
        local_x = int((x - self.chunk_x * CHUNK_WIDTH) // HEX_WIDTH)
        local_y = int((y - self.chunk_y * CHUNK_HEIGHT) // HEX_V_SPACING)
        return self.hexagons.get((local_x, local_y))
    
    def remove_hexagon(self, x, y):
        local_x = int((x - self.chunk_x * CHUNK_WIDTH) // HEX_WIDTH)
        local_y = int((y - self.chunk_y * CHUNK_HEIGHT) // HEX_V_SPACING)
        if (local_x, local_y) in self.hexagons:
            del self.hexagons[(local_x, local_y)]
            self.modified = True
            return True
        return False

class World:
    def __init__(self):
        self.chunks = {}
        self.seed = random.randint(0, 1000000)
        
    def get_chunk_coords(self, x, y):
        chunk_x = int(x // CHUNK_WIDTH)
        chunk_y = int(y // CHUNK_HEIGHT)
        return chunk_x, chunk_y
    
    def get_chunk(self, chunk_x, chunk_y):
        if (chunk_x, chunk_y) not in self.chunks:
            self.chunks[(chunk_x, chunk_y)] = Chunk(chunk_x, chunk_y)
            self.generate_chunk(chunk_x, chunk_y)
        return self.chunks[(chunk_x, chunk_y)]
    
    def generate_chunk(self, chunk_x, chunk_y):
        chunk = self.chunks[(chunk_x, chunk_y)]
        world_x, world_y = chunk.get_world_position()
        
        # Fixed & more reliable ground level
        base_ground_y = 500                     # start around middle of screen
        height_variation = int(60 * math.sin(chunk_x * 0.4 + self.seed * 0.007) +
                               45 * math.cos(chunk_y * 0.55 + self.seed * 0.009))
        
        ground_y = base_ground_y + height_variation
        
        # Fill from surface down
        for local_row in range(CHUNK_SIZE):
            for local_col in range(CHUNK_SIZE):
                x_offset = (HEX_WIDTH / 2) if local_row % 2 == 1 else 0
                hx = world_x + local_col * HEX_WIDTH + x_offset - HEX_WIDTH / 2
                hy = ground_y + local_row * HEX_V_SPACING
                
                depth = local_row
                
                if depth == 0:
                    block_type = "grass"
                elif depth <= 4:
                    block_type = "dirt"
                elif depth <= 12:
                    block_type = "stone"
                    if random.random() < 0.07:
                        block_type = "coal"
                elif depth <= 25:
                    block_type = "stone"
                    r = random.random()
                    if r < 0.04: block_type = "iron"
                    elif r < 0.07: block_type = "coal"
                else:
                    block_type = "stone"
                    r = random.random()
                    if r < 0.025: block_type = "gold"
                    elif r < 0.045: block_type = "diamond"
                
                hex_obj = Hexagon(hx, hy, HEX_SIZE, block_type)
                chunk.add_hexagon(hx, hy, hex_obj)
        
        for local_row in range(CHUNK_SIZE):
            for local_col in range(CHUNK_SIZE):
                x_offset = (HEX_WIDTH / 2) if local_row % 2 == 1 else 0
                x = world_x + local_col * HEX_WIDTH + x_offset - HEX_WIDTH / 2
                y = ground_y + local_row * HEX_V_SPACING
                
                # Use simple noise for biome determination
                
                # Determine biome based on chunk position
                biome = "plains"
                biome_noise = np.sin(chunk_x * 0.3 + self.seed) * np.cos(chunk_y * 0.3 + self.seed)
                if biome_noise < -0.3:
                    biome = "desert"
                elif biome_noise > 0.3:
                    biome = "mountain"
                elif -0.1 < biome_noise < 0.1:
                    biome = "forest"
                
                # Determine block type based on depth and biome
                depth = local_row
                if depth == 0:
                    if biome == "desert":
                        block_type = 'sand'
                    elif biome == "forest":
                        block_type = 'grass'
                    else:
                        block_type = 'grass'
                elif depth < 3:
                    if biome == "desert":
                        block_type = 'sand'
                    else:
                        block_type = 'dirt'
                elif depth < 6:
                    if random.random() < 0.1:
                        block_type = 'coal'
                    else:
                        block_type = 'stone'
                elif depth < 12:
                    rand_val = random.random()
                    if rand_val < 0.05:
                        block_type = 'coal'
                    elif rand_val < 0.08:
                        block_type = 'iron'
                    else:
                        block_type = 'stone'
                else:
                    rand_val = random.random()
                    if rand_val < 0.03:
                        block_type = 'coal'
                    elif rand_val < 0.05:
                        block_type = 'iron'
                    elif rand_val < 0.06:
                        block_type = 'gold'
                    elif rand_val < 0.07 and biome == "mountain":
                        block_type = 'diamond'
                    else:
                        block_type = 'stone'
                
                hexagon = Hexagon(x, y, HEX_SIZE, block_type)
                chunk.add_hexagon(x, y, hexagon)
    
    def get_nearby_hexagons(self, center_x, center_y, radius=200):
        nearby = []
        
        min_chunk_x, min_chunk_y = self.get_chunk_coords(center_x - radius, center_y - radius)
        max_chunk_x, max_chunk_y = self.get_chunk_coords(center_x + radius, center_y + radius)
        
        for chunk_x in range(min_chunk_x - 1, max_chunk_x + 2):
            for chunk_y in range(min_chunk_y - 1, max_chunk_y + 2):
                if (chunk_x, chunk_y) in self.chunks:
                    chunk = self.chunks[(chunk_x, chunk_y)]
                    chunk.last_accessed = time.time()
                    
                    for hexagon in chunk.hexagons.values():
                        if abs(hexagon.x - center_x) <= radius + HEX_SIZE and abs(hexagon.y - center_y) <= radius + HEX_SIZE:
                            nearby.append(hexagon)
        
        return nearby
    
    def get_hexagon_at(self, x, y):
        chunk_x, chunk_y = self.get_chunk_coords(x, y)
        if (chunk_x, chunk_y) in self.chunks:
            return self.chunks[(chunk_x, chunk_y)].get_hexagon(x, y)
        return None
    
    def remove_hexagon_at(self, x, y):
        chunk_x, chunk_y = self.get_chunk_coords(x, y)
        if (chunk_x, chunk_y) in self.chunks:
            return self.chunks[(chunk_x, chunk_y)].remove_hexagon(x, y)
        return False
    
    def add_hexagon_at(self, x, y, block_type):
        chunk_x, chunk_y = self.get_chunk_coords(x, y)
        chunk = self.get_chunk(chunk_x, chunk_y)
        
        # Find nearest valid hex position
        world_x, world_y = chunk.get_world_position()
        for local_row in range(CHUNK_SIZE):
            for local_col in range(CHUNK_SIZE):
                hex_x = world_x + local_col * HEX_WIDTH + (HEX_WIDTH / 2 if local_row % 2 == 1 else 0) - HEX_WIDTH / 2
                hex_y = world_y + local_row * HEX_V_SPACING
                
                if math.sqrt((hex_x - x) ** 2 + (hex_y - y) ** 2) < HEX_SIZE:
                    hexagon = Hexagon(hex_x, hex_y, HEX_SIZE, block_type)
                    chunk.add_hexagon(hex_x, hex_y, hexagon)
                    return hexagon
        
        # If no nearby position found, place at exact location
        hexagon = Hexagon(x, y, HEX_SIZE, block_type)
        chunk.add_hexagon(x, y, hexagon)
        return hexagon
    
    def unload_distant_chunks(self, player_x, player_y, max_chunks=50):
        # Sort chunks by last accessed time
        chunk_list = list(self.chunks.items())
        chunk_list.sort(key=lambda x: x[1].last_accessed, reverse=True)
        
        # Keep only the most recently accessed chunks
        if len(chunk_list) > max_chunks:
            for coords, chunk in chunk_list[max_chunks:]:
                del self.chunks[coords]
    
    def save_to_file(self, filename):
        data = {
            'chunks': {},
            'seed': self.seed
        }
        
        for (chunk_x, chunk_y), chunk in self.chunks.items():
            if chunk.modified:
                chunk_data = {
                    'hexagons': {},
                    'modified': chunk.modified
                }
                for (local_x, local_y), hexagon in chunk.hexagons.items():
                    chunk_data['hexagons'][(local_x, local_y)] = {
                        'x': hexagon.x,
                        'y': hexagon.y,
                        'block_type': hexagon.block_type,
                        'health': hexagon.health
                    }
                data['chunks'][(chunk_x, chunk_y)] = chunk_data
        
        with open(filename, 'wb') as f:
            pickle.dump(data, f)
    
    def load_from_file(self, filename):
        try:
            with open(filename, 'rb') as f:
                data = pickle.load(f)
            
            if not isinstance(data, dict) or 'chunks' not in data:
                return False
            
            self.seed = data.get('seed', random.randint(0, 1000000))
            
            for (chunk_x, chunk_y), chunk_data in data['chunks'].items():
                chunk = Chunk(chunk_x, chunk_y)
                chunk.modified = chunk_data.get('modified', True)
                
                for (local_x, local_y), hex_data in chunk_data['hexagons'].items():
                    hexagon = Hexagon(
                        hex_data['x'],
                        hex_data['y'],
                        HEX_SIZE,
                        hex_data['block_type']
                    )
                    hexagon.health = hex_data.get('health', 100)
                    hexagon.chunk_x = chunk_x
                    hexagon.chunk_y = chunk_y
                    chunk.hexagons[(local_x, local_y)] = hexagon
                
                self.chunks[(chunk_x, chunk_y)] = chunk
            
            return True
        except (FileNotFoundError, EOFError, pickle.UnpicklingError):
            return False

class Player:
    def __init__(self, x, y):
        self.x = x
        self.y = y
        self.radius = 15
        self.vel_x = 0
        self.vel_y = 0
        self.on_ground = False
        self.color = BLUE
        self.inventory = {bt: 0 for bt in COLOR_BY_TYPE}
        self.selected_slot = 0
        self.flight_mode = False  # Creative mode flying
        
    def get_selected_block(self):
        block_types_list = list(COLOR_BY_TYPE.keys())
        if self.selected_slot < len(block_types_list):
            return block_types_list[self.selected_slot]
        return 'dirt'
    
    def update(self, keys, hexagons):
        if self.flight_mode:
            # Creative mode flight controls
            if keys[pygame.K_a]:
                self.vel_x = -PLAYER_SPEED
            elif keys[pygame.K_d]:
                self.vel_x = PLAYER_SPEED
            else:
                self.vel_x *= 0.9
            
            if keys[pygame.K_w]:
                self.vel_y = -PLAYER_SPEED
            elif keys[pygame.K_s]:
                self.vel_y = PLAYER_SPEED
            else:
                self.vel_y *= 0.9
        else:
            # Survival mode physics
            if keys[pygame.K_a]:
                self.vel_x = -PLAYER_SPEED
            elif keys[pygame.K_d]:
                self.vel_x = PLAYER_SPEED
            else:
                self.vel_x *= FRICTION
            
            self.vel_y += GRAVITY
        
        self.x += self.vel_x
        self.y += self.vel_y
        
        # Collision with hexagons
        if not self.flight_mode:
            self.on_ground = False
            hexagons_sorted = sorted(hexagons, key=lambda h: (h.x - self.x)**2 + (h.y - self.y)**2)
            for hexagon in hexagons_sorted:
                self.check_collision(hexagon)
    
    def check_collision(self, hexagon):
        dx = self.x - hexagon.x
        dy = self.y - hexagon.y
        distance_sq = dx * dx + dy * dy
        
        max_distance_sq = (self.radius + hexagon.size * 1.5) ** 2
        if distance_sq > max_distance_sq:
            return
        
        # Skip collision with transparent blocks
        if hexagon.transparent:
            return
        
        player_bottom = self.y + self.radius
        surface_y = hexagon.get_top_surface_y(self.x)
        
        if surface_y is not None:
            if self.vel_y >= 0 and player_bottom >= surface_y - 10 and player_bottom <= surface_y + 15:
                self.y = surface_y - self.radius
                self.vel_y = 0
                self.on_ground = True
                return
        
        if distance_sq < (self.radius + hexagon.size * 0.9) ** 2:
            if self.y < hexagon.y:
                overlap = (self.radius + hexagon.size * 0.9) - math.sqrt(distance_sq)
                if overlap > 0:
                    self.y -= overlap
                    self.vel_y = 0
                    self.on_ground = True
                    return
            
            if abs(dy) < hexagon.size * 0.7:
                if dx > 0:
                    self.x = hexagon.x + hexagon.size * 0.95 + self.radius
                else:
                    self.x = hexagon.x - hexagon.size * 0.95 - self.radius
                self.vel_x = 0
    
    def jump(self):
        if self.on_ground or self.flight_mode:
            self.vel_y = JUMP_FORCE
            if not self.flight_mode:
                self.on_ground = False
    
    def add_to_inventory(self, block_type):
        if block_type in self.inventory:
            self.inventory[block_type] += 1
    
    def remove_from_inventory(self, block_type):
        if self.inventory.get(block_type, 0) > 0:
            self.inventory[block_type] -= 1
            return True
        return False

    def draw_offset(self, screen, offset_x, offset_y):
        screen_x = int(self.x - offset_x)
        screen_y = int(self.y - offset_y)
        pygame.draw.circle(screen, self.color, (screen_x, screen_y), self.radius)
        # Draw eyes
        pygame.draw.circle(screen, WHITE, (screen_x - 5, screen_y - 3), 4)
        pygame.draw.circle(screen, WHITE, (screen_x + 5, screen_y - 3), 4)
        pygame.draw.circle(screen, BLACK, (screen_x - 5, screen_y - 3), 2)
        pygame.draw.circle(screen, BLACK, (screen_x + 5, screen_y - 3), 2)

        # Draw mining range indicator
        if not self.flight_mode:
            pygame.draw.circle(screen, (255, 255, 255, 50), (screen_x, screen_y), MINING_RANGE, 1)

class Button:
    def __init__(self, x, y, width, height, text, color, hover_color, action):
        self.rect = pygame.Rect(x, y, width, height)
        self.text = text
        self.color = color
        self.hover_color = hover_color
        self.action = action
        self.hovered = False
        
    def check_hover(self, mouse_pos):
        self.hovered = self.rect.collidepoint(mouse_pos)
        return self.hovered
    
    def handle_event(self, event):
        if event.type == pygame.MOUSEBUTTONDOWN and event.button == 1 and self.hovered:
            self.action()
            return True
        return False
    
    def draw(self, screen, font):
        color = self.hover_color if self.hovered else self.color
        pygame.draw.rect(screen, color, self.rect, border_radius=8)
        pygame.draw.rect(screen, WHITE, self.rect, 2, border_radius=8)
        
        text_surface = font.render(self.text, True, WHITE)
        text_rect = text_surface.get_rect(center=self.rect.center)
        screen.blit(text_surface, text_rect)

class Game:
    def __init__(self):
        self.screen = pygame.display.set_mode((SCREEN_WIDTH, SCREEN_HEIGHT))
        pygame.display.set_caption("Tesselbox")
        self.clock = pygame.time.Clock()
        self.state = "MENU"
        self.running = True
        self.player = None
        self.world = World()
        self.particles = []
        self.mining = False
        self.camera_x = 0
        self.camera_y = 0
        self.sky_surface = self._create_sky_gradient()
        self.paused = False
        self.game_mode = "survival"  # "survival" or "creative"
        self.render_distance = RENDER_DISTANCE
        
        # UI Elements
        self.buttons = []
        self._create_menu_buttons()
        self._create_pause_buttons()
        
    def _create_sky_gradient(self):
        surface = pygame.Surface((SCREEN_WIDTH, SCREEN_HEIGHT))
        for y in range(SCREEN_HEIGHT):
            ratio = y / SCREEN_HEIGHT
            color = (
                int(135 - ratio * 85),
                int(206 - ratio * 106),
                int(235 - ratio * 35)
            )
            pygame.draw.line(surface, color, (0, y), (SCREEN_WIDTH, y))
        return surface

    def _create_menu_buttons(self):
        button_width = 200
        button_height = 50
        center_x = SCREEN_WIDTH // 2 - button_width // 2
        
        self.menu_buttons = [
            Button(center_x, 300, button_width, button_height, "Survival Mode", 
                   GREEN, (120, 220, 120), lambda: self.start_game("survival")),
            Button(center_x, 370, button_width, button_height, "Creative Mode", 
                   BLUE, (70, 170, 255), lambda: self.start_game("creative")),
            Button(center_x, 440, button_width, button_height, "Exit", 
                   RED, (220, 120, 120), lambda: self.quit_game())
        ]
    
    def _create_pause_buttons(self):
        button_width = 200
        button_height = 50
        center_x = SCREEN_WIDTH // 2 - button_width // 2
        
        self.pause_buttons = [
            Button(center_x, 250, button_width, button_height, "Resume", 
                   GREEN, (120, 220, 120), self.resume_game),
            Button(center_x, 320, button_width, button_height, "Save Game", 
                   BLUE, (70, 170, 255), self.save_game),
            Button(center_x, 390, button_width, button_height, "Back to Menu", 
                   RED, (220, 120, 120), self.back_to_menu)
        ]
    
    def quit_game(self):
        self.running = False
    
    def start_game(self, mode):
        self.game_mode = mode
        
        # Try to load existing world
        if self.world.load_from_file('world.pkl'):
            # Find player spawn point from saved data
            try:
                with open('player.pkl', 'rb') as f:
                    player_data = pickle.load(f)
                self.player = Player(player_data['x'], player_data['y'])
                self.player.inventory = player_data['inventory']
                self.player.selected_slot = player_data['selected_slot']
            except:
                # If player data fails, spawn at default location
                self._spawn_player()
        else:
            self._spawn_player()
        
        self.player.flight_mode = (mode == "creative")
        self.particles = []
        self.state = "GAME"
    
    def _spawn_player(self):
        # Safer spawn — start above likely ground level
        spawn_x = SCREEN_WIDTH // 2 + 100
        spawn_y = 200   # high enough to fall onto surface in most cases
        
        self.player = Player(spawn_x, spawn_y)
        
        if self.game_mode == "creative":
            for bt in COLOR_BY_TYPE:
                if COLLECTIBLE_BY_TYPE.get(bt, False):
                    self.player.inventory[bt] = 999
    
    def resume_game(self):
        self.paused = False
        self.state = "GAME"
    
    def save_game(self):
        if self.player:
            self.world.save_to_file('world.pkl')
            player_data = {
                'x': self.player.x,
                'y': self.player.y,
                'inventory': self.player.inventory,
                'selected_slot': self.player.selected_slot
            }
            with open('player.pkl', 'wb') as f:
                pickle.dump(player_data, f)
    
    def back_to_menu(self):
        self.save_game()
        self.state = "MENU"
        self.paused = False
    
    def handle_menu_input(self, event):
        mouse_pos = pygame.mouse.get_pos()
        for button in self.menu_buttons:
            button.check_hover(mouse_pos)
            button.handle_event(event)
        
        if event.type == pygame.KEYDOWN:
            if event.key == pygame.K_ESCAPE:
                self.quit_game()
    
    def handle_game_input(self, event):
        if event.type == pygame.KEYDOWN:
            if event.key == pygame.K_SPACE:
                self.player.jump()
            elif event.key == pygame.K_ESCAPE:
                self.paused = not self.paused
                if self.paused:
                    self.state = "PAUSE"
                else:
                    self.state = "GAME"
            elif event.key == pygame.K_1:
                self.player.selected_slot = 0
            elif event.key == pygame.K_2:
                self.player.selected_slot = 1
            elif event.key == pygame.K_3:
                self.player.selected_slot = 2
            elif event.key == pygame.K_4:
                self.player.selected_slot = 3
            elif event.key == pygame.K_5:
                self.player.selected_slot = 4
            elif event.key == pygame.K_6:
                self.player.selected_slot = 5
            elif event.key == pygame.K_7:
                self.player.selected_slot = 6
            elif event.key == pygame.K_8:
                self.player.selected_slot = 7
            elif event.key == pygame.K_9:
                self.player.selected_slot = 8
            elif event.key == pygame.K_EQUALS:  # Increase render distance
                self.render_distance = min(self.render_distance + 1, 10)
            elif event.key == pygame.K_MINUS:  # Decrease render distance
                self.render_distance = max(self.render_distance - 1, 2)
            elif event.key == pygame.K_f and self.game_mode == "creative":
                self.player.flight_mode = not self.player.flight_mode
        
        if event.type == pygame.MOUSEBUTTONDOWN:
            if event.button == 1:
                self.mining = True
            elif event.button == 3:
                self.place_block()
            elif event.button == 4:  # Mouse wheel up
                self.player.selected_slot = max(0, self.player.selected_slot - 1)
            elif event.button == 5:  # Mouse wheel down
                self.player.selected_slot = min(8, self.player.selected_slot + 1)
        
        if event.type == pygame.MOUSEBUTTONUP:
            if event.button == 1:
                self.mining = False
    
    def handle_pause_input(self, event):
        mouse_pos = pygame.mouse.get_pos()
        for button in self.pause_buttons:
            button.check_hover(mouse_pos)
            button.handle_event(event)
        
        if event.type == pygame.KEYDOWN:
            if event.key == pygame.K_ESCAPE:
                self.resume_game()
    
    def raycast_to_block(self, start_x, start_y, target_x, target_y, max_distance=MINING_RANGE):
        """Raycast from player position to find the first block in line of sight"""
        dx = target_x - start_x
        dy = target_y - start_y
        distance = math.sqrt(dx * dx + dy * dy)
        
        if distance > max_distance or distance == 0:
            return None, None
        
        # Normalize direction
        dx /= distance
        dy /= distance
        
        # Step along the ray
        step_size = 5
        for d in range(0, int(distance), step_size):
            x = start_x + dx * d
            y = start_y + dy * d
            
            # Check for block at this position
            block = self.world.get_hexagon_at(x, y)
            if block and not block.transparent:
                # Calculate the position for placing a block (just outside the hit block)
                place_x = start_x + dx * (d - step_size)
                place_y = start_y + dy * (d - step_size)
                return block, (place_x, place_y)
        
        return None, None
    
    def mine_block(self):
        mouse_pos = pygame.mouse.get_pos()
        world_mx = mouse_pos[0] + self.camera_x
        world_my = mouse_pos[1] + self.camera_y
        
        # Use raycasting for accurate block selection
        target_block, _ = self.raycast_to_block(self.player.x, self.player.y, world_mx, world_my)
        
        if target_block:
            if target_block.take_damage(5):
                # Create particles
                for _ in range(10):
                    particle = Particle(target_block.x, target_block.y, target_block.color)
                    self.particles.append(particle)
                
                self.player.add_to_inventory(target_block.block_type)
                self.world.remove_hexagon_at(target_block.x, target_block.y)
    
    def place_block(self):
        mouse_pos = pygame.mouse.get_pos()
        world_mx = mouse_pos[0] + self.camera_x
        world_my = mouse_pos[1] + self.camera_y
        
        selected_block = self.player.get_selected_block()
        
        # In creative mode, we have infinite blocks
        if self.game_mode == "creative":
            # Use raycasting to find where to place
            target_block, place_pos = self.raycast_to_block(self.player.x, self.player.y, world_mx, world_my)
            
            if place_pos:
                place_x, place_y = place_pos
                
                # Check if position is valid (not overlapping player or too far)
                dx = self.player.x - place_x
                dy = self.player.y - place_y
                if dx * dx + dy * dy > MINING_RANGE * MINING_RANGE:
                    return
                
                # Check if not overlapping player
                if dx * dx + dy * dy < (HEX_SIZE + self.player.radius) ** 2:
                    return
                
                # Check if there's already a block there
                existing = self.world.get_hexagon_at(place_x, place_y)
                if existing:
                    return
                
                self.world.add_hexagon_at(place_x, place_y, selected_block)
        else:
            # Survival mode requires inventory
            if not self.player.remove_from_inventory(selected_block):
                return
            
            target_block, place_pos = self.raycast_to_block(self.player.x, self.player.y, world_mx, world_my)
            
            if place_pos:
                place_x, place_y = place_pos
                
                dx = self.player.x - place_x
                dy = self.player.y - place_y
                if dx * dx + dy * dy > MINING_RANGE * MINING_RANGE:
                    self.player.add_to_inventory(selected_block)
                    return
                
                if dx * dx + dy * dy < (HEX_SIZE + self.player.radius) ** 2:
                    self.player.add_to_inventory(selected_block)
                    return
                
                existing = self.world.get_hexagon_at(place_x, place_y)
                if existing:
                    self.player.add_to_inventory(selected_block)
                    return
                
                self.world.add_hexagon_at(place_x, place_y, selected_block)

    def update_game(self):
        keys = pygame.key.get_pressed()
        
        # Update camera - no clamping = can explore left & up
        self.camera_x = self.player.x - SCREEN_WIDTH // 2
        self.camera_y = self.player.y - SCREEN_HEIGHT // 2
        # self.camera_x = max(0, self.camera_x)   # comment these out
        # self.camera_y = max(0, self.camera_y)
        
        # Only process hexagons near player
        nearby_hexagons = self.world.get_nearby_hexagons(self.player.x, self.player.y, radius=500)
        self.player.update(keys, nearby_hexagons)
        
        # Update hover states
        mouse_pos = pygame.mouse.get_pos()
        world_mx = mouse_pos[0] + self.camera_x
        world_my = mouse_pos[1] + self.camera_y
        hover_nearby = self.world.get_nearby_hexagons(world_mx, world_my, radius=HEX_SIZE * 3)
        
        for hexagon in hover_nearby:
            hexagon.check_hover(world_mx, world_my, self.player.x, self.player.y)
        
        if self.mining:
            self.mine_block()
        
        # Update particles
        self.particles = [p for p in self.particles if not (p.update() or p.is_dead())]
        
        # Unload distant chunks
        self.world.unload_distant_chunks(self.player.x, self.player.y)
    
    def draw_menu(self):
        self.screen.fill(BLACK)
        
        title_font = pygame.font.Font(None, 74)
        title = title_font.render("TESSELBOX", True, WHITE)
        title_rect = title.get_rect(center=(SCREEN_WIDTH // 2, 150))
        self.screen.blit(title, title_rect)
        
        subtitle_font = pygame.font.Font(None, 36)
        subtitle = subtitle_font.render("Enhanced Minecraft-like Experience", True, GRAY)
        subtitle_rect = subtitle.get_rect(center=(SCREEN_WIDTH // 2, 200))
        self.screen.blit(subtitle, subtitle_rect)
        
        mouse_pos = pygame.mouse.get_pos()
        for button in self.menu_buttons:
            button.check_hover(mouse_pos)
            button.draw(self.screen, pygame.font.Font(None, 32))
        
        # Draw controls
        font = pygame.font.Font(None, 24)
        controls = [
            "Controls:",
            "WASD - Move",
            "Space - Jump",
            "Left Click - Mine",
            "Right Click - Place",
            "1-9 / Scroll - Select Block",
            "+/- - Render Distance",
            "F - Toggle Flight (Creative)",
            "ESC - Pause/Menu"
        ]
        
        y_pos = 520
        for control in controls:
            text = font.render(control, True, DARK_GRAY)
            text_rect = text.get_rect(center=(SCREEN_WIDTH // 2, y_pos))
            self.screen.blit(text, text_rect)
            y_pos += 25
    
    def draw_pause_menu(self):
        # Draw semi-transparent overlay
        overlay = pygame.Surface((SCREEN_WIDTH, SCREEN_HEIGHT))
        overlay.set_alpha(128)
        overlay.fill(BLACK)
        self.screen.blit(overlay, (0, 0))
        
        title_font = pygame.font.Font(None, 74)
        title = title_font.render("PAUSED", True, WHITE)
        title_rect = title.get_rect(center=(SCREEN_WIDTH // 2, 150))
        self.screen.blit(title, title_rect)
        
        mouse_pos = pygame.mouse.get_pos()
        for button in self.pause_buttons:
            button.check_hover(mouse_pos)
            button.draw(self.screen, pygame.font.Font(None, 32))
    
    def draw_game(self):
        # Use pre-rendered sky gradient
        self.screen.blit(self.sky_surface, (0, 0))
        
        # Get visible chunks based on render distance
        player_chunk_x, player_chunk_y = self.world.get_chunk_coords(self.player.x, self.player.y)
        
        visible_hexagons = []
        for chunk_x in range(player_chunk_x - self.render_distance, player_chunk_x + self.render_distance + 1):
            for chunk_y in range(player_chunk_y - self.render_distance, player_chunk_y + self.render_distance + 1):
                if (chunk_x, chunk_y) in self.world.chunks:
                    chunk = self.world.chunks[(chunk_x, chunk_y)]
                    # Frustum culling - only add hexagons that are on screen
                    for hexagon in chunk.hexagons.values():
                        screen_x = hexagon.x - self.camera_x
                        screen_y = hexagon.y - self.camera_y
                        
                        if -HEX_SIZE < screen_x < SCREEN_WIDTH + HEX_SIZE and -HEX_SIZE < screen_y < SCREEN_HEIGHT + HEX_SIZE:
                            visible_hexagons.append(hexagon)
        
        # Sort by y position for proper depth rendering
        visible_hexagons.sort(key=lambda h: h.y)
        
        # Draw hexagons
        for hexagon in visible_hexagons:
            offset_corners = [(cx - self.camera_x, cy - self.camera_y) for cx, cy in hexagon.corners]
            pygame.draw.polygon(self.screen, hexagon.active_color, offset_corners)
            pygame.draw.polygon(self.screen, DARK_GRAY, offset_corners, 2)
            
            if hexagon.health < hexagon.max_health:
                health_ratio = hexagon.health / hexagon.max_health
                crack_color = (255, int(255 * health_ratio), int(255 * health_ratio))
                pygame.draw.polygon(self.screen, crack_color, offset_corners, 3)
        
        # Draw particles
        for particle in self.particles:
            px = particle.x - self.camera_x
            py = particle.y - self.camera_y
            if -10 < px < SCREEN_WIDTH + 10 and -10 < py < SCREEN_HEIGHT + 10:
                particle.draw_offset(self.screen, self.camera_x, self.camera_y)
        
        self.player.draw_offset(self.screen, self.camera_x, self.camera_y)
        self.draw_ui()
    
    def draw_ui(self):
        # Draw hotbar
        hotbar_width = 400
        hotbar_height = 50
        hotbar_x = SCREEN_WIDTH // 2 - hotbar_width // 2
        hotbar_y = SCREEN_HEIGHT - hotbar_height - 10
        
        block_types_list = list(COLOR_BY_TYPE.keys())
        slot_width = hotbar_width // 9
        
        for i in range(9):
            slot_x = hotbar_x + i * slot_width
            slot_rect = pygame.Rect(slot_x, hotbar_y, slot_width - 2, hotbar_height)
            
            # Highlight selected slot
            if i == self.player.selected_slot:
                pygame.draw.rect(self.screen, WHITE, slot_rect, 3)
            else:
                pygame.draw.rect(self.screen, GRAY, slot_rect, 1)
            
            if i < len(block_types_list):
                block_type = block_types_list[i]
                color = COLOR_BY_TYPE[block_type]
                
                # Draw block preview (handle possible RGBA colors)
                draw_color = color[:3] if len(color) == 4 else color
                pygame.draw.rect(self.screen, draw_color,
                               (slot_x + 5, hotbar_y + 5, slot_width - 12, hotbar_height - 10))
                
                # Draw count (only in survival mode)
                if self.game_mode == "survival":
                    count = self.player.inventory.get(block_type, 0)
                    font = pygame.font.Font(None, 20)
                    count_text = font.render(str(count), True, WHITE)
                    count_rect = count_text.get_rect(bottomright=(slot_x + slot_width - 5, hotbar_y + hotbar_height - 5))
                    self.screen.blit(count_text, count_rect)
                
                # Draw block name on hover
                mouse_pos = pygame.mouse.get_pos()
                if slot_rect.collidepoint(mouse_pos):
                    font = pygame.font.Font(None, 24)
                    name_text = font.render(BLOCK_DEFINITIONS[block_type]['name'], True, WHITE)
                    name_rect = name_text.get_rect(midtop=(slot_x + slot_width // 2, hotbar_y - 5))
                    
                    # Draw tooltip background
                    tooltip_rect = name_rect.inflate(10, 10)
                    pygame.draw.rect(self.screen, (50, 50, 50), tooltip_rect, border_radius=5)
                    pygame.draw.rect(self.screen, WHITE, tooltip_rect, 1, border_radius=5)
                    self.screen.blit(name_text, name_rect)
        
        # Draw mode indicator
        font = pygame.font.Font(None, 24)
        mode_text = f"Mode: {self.game_mode.capitalize()}"
        mode_surface = font.render(mode_text, True, WHITE)
        self.screen.blit(mode_surface, (10, 10))
        
        # Draw render distance
        render_text = f"Render Distance: {self.render_distance}"
        render_surface = font.render(render_text, True, WHITE)
        self.screen.blit(render_surface, (10, 35))
        
        # Draw flight status in creative mode
        if self.game_mode == "creative":
            flight_text = f"Flight: {'ON' if self.player.flight_mode else 'OFF'}"
            flight_surface = font.render(flight_text, True, WHITE)
            self.screen.blit(flight_surface, (10, 60))
    
    def run(self):
        while self.running:
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    self.save_game()
                    self.running = False
                
                if self.state == "MENU":
                    self.handle_menu_input(event)
                elif self.state == "GAME":
                    self.handle_game_input(event)
                elif self.state == "PAUSE":
                    self.handle_pause_input(event)
            
            if self.state == "MENU":
                self.draw_menu()
            elif self.state == "GAME":
                self.update_game()
                self.draw_game()
            elif self.state == "PAUSE":
                self.draw_pause_menu()
            
            pygame.display.flip()
            self.clock.tick(FPS)
        
        pygame.quit()
        sys.exit()

if __name__ == "__main__":
    game = Game()
    game.run()
