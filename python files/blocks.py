# blocks.py
from typing import Dict, Tuple, Any, List
from enum import Enum, auto

class BlockType(Enum):
    AIR = auto()
    DIRT = auto()
    GRASS = auto()
    STONE = auto()
    SAND = auto()
    WATER = auto()
    LOG = auto()
    LEAVES = auto()
    COAL_ORE = auto()
    IRON_ORE = auto()
    GOLD_ORE = auto()
    DIAMOND_ORE = auto()
    BEDROCK = auto()
    GLASS = auto()
    BRICK = auto()
    PLANK = auto()
    CACTUS = auto()

# Block properties
# fmt: off
BLOCK_DEFINITIONS = {
    "air": {
        "id":          BlockType.AIR,
        "name":        "Air",
        "color":       (0, 0, 0, 0),          # fully transparent
        "hardness":    0.0,
        "transparent": True,
        "solid":       False,
        "collectible": False,
        "flammable":   False,
        "light_level": 0,
    },
    "dirt": {
        "id":          BlockType.DIRT,
        "name":        "Dirt",
        "color":       (139, 90, 43),
        "hardness":    1.0,
        "transparent": False,
        "solid":       True,
        "collectible": True,
        "flammable":   False,
        "light_level": 0,
    },
    "grass": {
        "id":          BlockType.GRASS,
        "name":        "Grass Block",
        "color":       (100, 200, 100),
        "top_color":   (126, 200, 80),  # Slightly different color for top face
        "side_color":  (95, 150, 30),   # Darker color for sides
        "hardness":    1.0,
        "transparent": False,
        "solid":       True,
        "collectible": True,
        "flammable":   False,
        "light_level": 0,
    },
    "stone": {
        "id":          BlockType.STONE,
        "name":        "Stone",
        "color":       (169, 169, 169),
        "hardness":    2.0,
        "transparent": False,
        "solid":       True,
        "collectible": True,
        "flammable":   False,
        "light_level": 0,
    },
    "sand": {
        "id":          BlockType.SAND,
        "name":        "Sand",
        "color":       (238, 214, 175),
        "hardness":    0.8,
        "transparent": False,
        "solid":       True,
        "collectible": True,
        "gravity":     True,
        "flammable":   False,
        "light_level": 0,
    },
    "water": {
        "id":          BlockType.WATER,
        "name":        "Water",
        "color":       (64, 164, 223, 140),
        "hardness":    0.0,
        "transparent": True,
        "solid":       False,
        "collectible": False,
        "viscosity":   0.8,  # How much it slows down entities
        "light_level": 1,    # Slight glow
    },
    "coal_ore": {
        "id":          BlockType.COAL_ORE,
        "name":        "Coal Ore",
        "color":       (45, 45, 45),
        "hardness":    3.0,
        "transparent": False,
        "solid":       True,
        "collectible": True,
        "drops":       "coal",
        "light_level": 0,
    },
    "iron_ore": {
        "id":          BlockType.IRON_ORE,
        "name":        "Iron Ore",
        "color":       (180, 150, 140),
        "hardness":    3.0,
        "transparent": False,
        "solid":       True,
        "collectible": True,
        "drops":       "raw_iron",
        "light_level": 0,
    },
    "gold_ore": {
        "id":          BlockType.GOLD_ORE,
        "name":        "Gold Ore",
        "color":       (250, 238, 77),
        "hardness":    3.0,
        "transparent": False,
        "solid":       True,
        "collectible": True,
        "drops":       "raw_gold",
        "light_level": 0,
    },
    "diamond_ore": {
        "id":          BlockType.DIAMOND_ORE,
        "name":        "Diamond Ore",
        "color":       (101, 240, 213),
        "hardness":    4.0,
        "transparent": False,
        "solid":       True,
        "collectible": True,
        "drops":       "diamond",
        "light_level": 2,  # Slight glow
    },
    "bedrock": {
        "id":          BlockType.BEDROCK,
        "name":        "Bedrock",
        "color":       (85, 85, 85),
        "hardness":    -1.0,  # Unbreakable
        "transparent": False,
        "solid":       True,
        "collectible": False,
        "light_level": 0,
    },
    "glass": {
        "id":          BlockType.GLASS,
        "name":        "Glass",
        "color":       (200, 220, 240, 100),
        "hardness":    0.5,
        "transparent": True,
        "solid":       True,
        "collectible": True,
        "light_level": 0,
    },
    "brick": {
        "id":          BlockType.BRICK,
        "name":        "Bricks",
        "color":       (178, 34, 34),
        "hardness":    2.0,
        "transparent": False,
        "solid":       True,
        "collectible": True,
        "light_level": 0,
    },
    "plank": {
        "id":          BlockType.PLANK,
        "name":        "Oak Planks",
        "color":       (199, 164, 96),
        "hardness":    1.0,
        "transparent": False,
        "solid":       True,
        "collectible": True,
        "flammable":   True,
        "light_level": 0,
    },
    "cactus": {
        "id":          BlockType.CACTUS,
        "name":        "Cactus",
        "color":       (88, 124, 39),
        "hardness":    0.8,
        "transparent": True,
        "solid":       True,
        "collectible": True,
        "damage":      1.0,  # Damage to entities on contact
        "light_level": 0,
    },
    # Add more blocks here later...
}

# Quick lookup tables (generated once at import)
COLOR_BY_TYPE: Dict[str, Tuple[int, ...]] = {k: v["color"] for k, v in BLOCK_DEFINITIONS.items()}
HARDNESS_BY_TYPE = {k: v["hardness"] for k, v in BLOCK_DEFINITIONS.items()}
TRANSPARENT_BY_TYPE = {k: v["transparent"] for k, v in BLOCK_DEFINITIONS.items()}
SOLID_BY_TYPE = {k: v["solid"] for k, v in BLOCK_DEFINITIONS.items()}
COLLECTIBLE_BY_TYPE = {k: v["collectible"] for k, v in BLOCK_DEFINITIONS.items()}
