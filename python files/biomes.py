# biomes.py
"""
Biome definitions and generation logic for terrain generation.
Provides biome-specific block distributions and characteristics.
"""

from typing import Dict, Tuple, List, Optional
import random
import math
import noise
from enum import Enum, auto

class BiomeType(Enum):
    PLAINS = auto()
    DESERT = auto()
    FOREST = auto()
    MOUNTAIN = auto()
    TAIGA = auto()
    SWAMP = auto()
    JUNGLE = auto()
    SAVANNA = auto()
    BADLANDS = auto()
    OCEAN = auto()
    DEEP_OCEAN = auto()
    BEACH = auto()
    RIVER = auto()
    FROZEN_RIVER = auto()
    STONE_SHORE = auto()

# Biome characteristics
BIOME_CHARACTERISTICS = {
    "plains": {
        "id": BiomeType.PLAINS,
        "name": "Plains",
        "surface_block": "grass",
        "subsurface_block": "dirt",
        "surface_variation": 0.3,
        "tree_density": 0.02,
        "color": (100, 200, 100),
        "temperature": 0.8,
        "humidity": 0.4,
        "height_scale": 0.5,
        "features": ["tall_grass", "flowers"],
        "structures": ["village", "pillager_outpost"],
    },
    "desert": {
        "id": BiomeType.DESERT,
        "name": "Desert",
        "surface_block": "sand",
        "subsurface_block": "sand",
        "surface_variation": 0.1,
        "tree_density": 0.0,
        "color": (238, 214, 175),
        "temperature": 2.0,
        "humidity": 0.0,
        "height_scale": 0.2,
        "features": ["dead_bush", "cactus", "sandstone"],
        "structures": ["desert_temple", "village_desert"],
    },
    "forest": {
        "id": BiomeType.FOREST,
        "name": "Forest",
        "surface_block": "grass",
        "subsurface_block": "dirt",
        "surface_variation": 0.5,
        "tree_density": 0.08,
        "color": (80, 180, 80),
        "temperature": 0.7,
        "humidity": 0.8,
        "height_scale": 0.4,
        "features": ["oak_trees", "birch_trees", "flowers", "mushrooms"],
        "structures": ["village_plains"],
    },
    "mountain": {
        "id": BiomeType.MOUNTAIN,
        "name": "Mountains",
        "surface_block": "stone",
        "subsurface_block": "stone",
        "surface_variation": 1.2,
        "tree_density": 0.01,
        "color": (150, 150, 150),
        "temperature": 0.2,
        "humidity": 0.3,
        "height_scale": 2.0,
        "features": ["emerald_ore", "silverfish"],
        "structures": ["mineshaft", "stronghold"],
    },
    "taiga": {
        "id": BiomeType.TAIGA,
        "name": "Taiga",
        "surface_block": "grass",
        "subsurface_block": "dirt",
        "surface_variation": 0.4,
        "tree_density": 0.05,
        "color": (49, 85, 44),
        "temperature": 0.25,
        "humidity": 0.8,
        "height_scale": 0.4,
        "features": ["spruce_trees", "ferns", "sweet_berry_bushes"],
        "structures": ["village_snowy", "igloo"],
    },
    "swamp": {
        "id": BiomeType.SWAMP,
        "name": "Swamp",
        "surface_block": "grass",
        "subsurface_block": "dirt",
        "surface_variation": 0.2,
        "tree_density": 0.1,
        "color": (106, 113, 57),
        "temperature": 0.8,
        "humidity": 0.9,
        "height_scale": 0.1,
        "features": ["oak_trees", "vines", "lily_pads", "blue_orchids"],
        "structures": ["swamp_hut"],
    },
    "jungle": {
        "id": BiomeType.JUNGLE,
        "name": "Jungle",
        "surface_block": "grass",
        "subsurface_block": "dirt",
        "surface_variation": 0.6,
        "tree_density": 0.15,
        "color": (83, 123, 9),
        "temperature": 0.95,
        "humidity": 0.9,
        "height_scale": 0.4,
        "features": ["jungle_trees", "cocoa_plants", "vines", "melons"],
        "structures": ["jungle_temple", "ruined_portal_jungle"],
    },
    "savanna": {
        "id": BiomeType.SAVANNA,
        "name": "Savanna",
        "surface_block": "grass",
        "subsurface_block": "dirt",
        "surface_variation": 0.3,
        "tree_density": 0.03,
        "color": (189, 178, 95),
        "temperature": 1.2,
        "humidity": 0.0,
        "height_scale": 0.2,
        "features": ["acacia_trees", "tall_grass", "village_savanna"],
        "structures": ["village_savanna", "pillager_outpost"],
    },
    "badlands": {
        "id": BiomeType.BADLANDS,
        "name": "Badlands",
        "surface_block": "red_sand",
        "subsurface_block": "red_sand",
        "surface_variation": 0.8,
        "tree_density": 0.0,
        "color": (217, 123, 51),
        "temperature": 2.0,
        "humidity": 0.0,
        "height_scale": 0.2,
        "features": ["terracotta", "gold_ore", "dead_bush"],
        "structures": ["mineshaft_mesa", "desert_pyramid"],
    },
    "ocean": {
        "id": BiomeType.OCEAN,
        "name": "Ocean",
        "surface_block": "water",
        "subsurface_block": "sand",
        "surface_variation": 0.1,
        "tree_density": 0.0,
        "color": (0, 0, 200, 100),
        "temperature": 0.5,
        "humidity": 0.5,
        "height_scale": 0.1,
        "features": ["kelp", "seagrass", "cod", "salmon"],
        "structures": ["ocean_ruin", "shipwreck"],
    },
    "beach": {
        "id": BiomeType.BEACH,
        "name": "Beach",
        "surface_block": "sand",
        "subsurface_block": "sand",
        "surface_variation": 0.05,
        "tree_density": 0.0,
        "color": (238, 214, 175),
        "temperature": 0.8,
        "humidity": 0.4,
        "height_scale": 0.05,
        "features": ["turtle_eggs", "seagrass"],
        "structures": [],
    },
    "snowy_mountains": {
        "id": BiomeType.SNOWY_MOUNTAINS,
        "name": "Snowy Mountains",
        "surface_block": "snow",
        "subsurface_block": "stone",
        "surface_variation": 1.2,
        "tree_density": 0.01,
        "color": (150, 150, 150),
        "temperature": 0.2,
        "humidity": 0.3,
        "height_scale": 2.0,
        "features": ["emerald_ore", "silverfish"],
        "structures": ["mineshaft", "stronghold"],
    },
    "snowy_tundra": {
        "id": BiomeType.SNOWY_TUNDRA,
        "name": "Snowy Tundra",
        "surface_block": "snow",
        "subsurface_block": "dirt",
        "surface_variation": 0.2,
        "tree_density": 0.05,
        "color": (49, 85, 44),
        "temperature": 0.25,
        "humidity": 0.8,
        "height_scale": 0.4,
        "features": ["spruce_trees", "ferns", "sweet_berry_bushes"],
        "structures": ["village_snowy", "igloo"],
    },
}

def get_biome_at_position(x: float, z: float, seed: int) -> str:
    """
    Get the biome at a specific (x, z) position using 2D Perlin noise.
    
    Args:
        x: X coordinate in world space
        z: Z coordinate in world space (for 3D biomes)
        seed: Random seed for consistent generation
        
    Returns:
        Biome name string
    """
    # Scale factors for different noise layers
    temp_scale = 400.0
    humidity_scale = 350.0
    height_scale = 200.0
    
    # Generate base noise values
    temp = noise.pnoise2(
        x / temp_scale, 
        z / temp_scale, 
        octaves=3, 
        persistence=0.5, 
        lacunarity=2.0, 
        repeatx=4096, 
        repeaty=4096, 
        base=seed
    )
    
    humidity = noise.pnoise2(
        (x + 1000) / humidity_scale, 
        (z + 1000) / humidity_scale, 
        octaves=3, 
        persistence=0.5, 
        lacunarity=2.0, 
        repeatx=4096, 
        repeaty=4096, 
        base=seed + 1
    )
    
    height = noise.pnoise2(
        (x + 2000) / height_scale, 
        (z + 2000) / height_scale, 
        octaves=5, 
        persistence=0.6, 
        lacunarity=1.8, 
        repeatx=4096, 
        repeaty=4096, 
        base=seed + 2
    )
    
    # Normalize values to 0-1 range
    temp = (temp + 1) / 2.0
    humidity = (humidity + 1) / 2.0
    height = (height + 1) / 2.0
    
    # Determine biome based on temperature, humidity, and height
    if height < 0.2:  # Ocean
        if height < 0.1:
            return "deep_ocean"
        return "ocean"
    elif height < 0.25:  # Beach/Shore
        return "beach"
    elif height > 0.8:  # Mountain
        if temp < 0.3:
            return "snowy_mountains"
        return "mountain"
    
    # Main biome selection based on temperature and humidity
    if temp < 0.2:  # Cold biomes
        if humidity > 0.7:
            return "taiga"
        return "snowy_tundra"
    elif temp < 0.4:  # Cool biomes
        if humidity > 0.7:
            return "forest"
        return "plains"
    elif temp < 0.6:  # Temperate biomes
        if humidity > 0.7:
            return "swamp"
        return "plains"
    elif temp < 0.8:  # Warm biomes
        if humidity > 0.7:
            return "jungle"
        return "savanna"
    else:  # Hot biomes
        if humidity > 0.5:
            return "badlands"
        return "desert"

def get_biome_block(biome: str, depth: int, world_x: float, world_y: float, seed: int) -> str:
    """
    Get the appropriate block type for a given biome, depth, and position.
    
    Args:
        biome: Biome name string
        depth: Depth below surface (in hexagons)
        world_x, world_y: World coordinates
        seed: Random seed for consistent generation
        
    Returns:
        Block type string
    """
    # Get biome info
    biome_info = BIOME_CHARACTERISTICS.get(biome, BIOME_CHARACTERISTICS["plains"])
    
    # Surface layer
    if depth == 0:
        return biome_info["surface_block"]
    
    # Subsurface layers
    elif depth < 3:
        # Add variation to subsurface layers
        if random.random() < 0.2:
            if biome in ["desert", "beach", "badlands"]:
                return "sand"
            return "dirt"
        return biome_info["subsurface_block"]
    
    # Underground layers with increasing stone frequency
    elif depth < 10:
        stone_chance = min(0.8, depth / 15.0 + 0.2)
        if random.random() < stone_chance:
            return "stone"
        return "dirt"
    
    # Deep underground layers (more stone, less dirt)
    else:
        # Generate consistent noise for ore veins
        ore_noise = noise.pnoise3(
            world_x / 50.0, 
            world_y / 50.0, 
            depth / 20.0,
            octaves=3,
            persistence=0.5,
            lacunarity=2.0,
            repeatx=1024,
            repeaty=1024,
            repeatz=1024,
            base=seed + 3
        )
        
        # Bedrock at the bottom
        if depth > 30:
            return "bedrock"
        
        # Ore generation
        if ore_noise > 0.8:
            # Rarer ores deeper down
            if depth > 20 and ore_noise > 0.95:
                if random.random() < 0.3:
                    return "diamond_ore"
                return "gold_ore"
            # More common ores
            elif depth > 10:
                if random.random() < 0.4:
                    return "iron_ore"
                return "coal_ore"
        
        # Default to stone
        return "stone"

def get_surface_height_variation(biome: str, world_x: float, world_z: float, seed: int) -> float:
    """
    Get the surface height variation for a biome at a given position.
    Uses 2D Perlin noise for more natural terrain.
    
    Args:
        biome: Biome name string
        world_x: World x coordinate
        world_z: World z coordinate (for 3D terrain)
        seed: Random seed for consistent generation
        
    Returns:
        Height offset value
    """
    # Get biome characteristics
    biome_info = BIOME_CHARACTERISTICS.get(biome, BIOME_CHARACTERISTICS["plains"])
    variation = biome_info["surface_variation"]
    
    # Generate terrain noise with multiple octaves
    scale = 100.0  # Base scale for large features
    height = 0.0
    amplitude = 1.0
    frequency = 1.0
    
    # Add multiple octaves of noise
    for _ in range(4):
        height += noise.pnoise2(
            world_x * frequency / scale, 
            world_z * frequency / scale,
            octaves=1,
            persistence=0.5,
            lacunarity=2.0,
            repeatx=4096,
            repeaty=4096,
            base=seed + 3
        ) * amplitude
        
        frequency *= 2.0
        amplitude *= 0.5
    
    # Add some smaller details
    detail = noise.pnoise2(
        world_x * 0.1, 
        world_z * 0.1,
        octaves=1,
        persistence=0.5,
        lacunarity=2.0,
        repeatx=4096,
        repeaty=4096,
        base=seed + 4
    ) * 0.1
    
    # Combine and scale by biome variation
    return (height + detail) * variation

def should_spawn_tree(biome: str, x: float, z: float, seed: int) -> bool:
    """
    Determine if a tree should spawn at a specific position in this biome.
    Uses noise for more natural tree distribution.
    
    Args:
        biome: Biome name string
        x, z: World coordinates
        seed: Random seed for consistent generation
        
    Returns:
        Boolean indicating if a tree should spawn
    """
    # Get base tree density from biome
    biome_info = BIOME_CHARACTERISTICS.get(biome, BIOME_CHARACTERISTICS["plains"])
    base_density = biome_info["tree_density"]
    
    # Skip if no trees in this biome
    if base_density <= 0:
        return False
    
    # Generate noise for tree distribution
    tree_noise = noise.pnoise2(
        x / 20.0,
        z / 20.0,
        octaves=2,
        persistence=0.5,
        lacunarity=2.0,
        repeatx=4096,
        repeaty=4096,
        base=seed + 5
    )
    
    # Map noise to 0-1 range
    noise_value = (tree_noise + 1) / 2.0
    
    # Check against threshold based on density
    threshold = 1.0 - (base_density * 10.0)
    return noise_value > threshold