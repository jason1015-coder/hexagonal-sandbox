# blocks.py
from typing import Dict, Tuple, Any

# fmt: off
BLOCK_DEFINITIONS = {
    "air": {
        "name":        "Air",
        "color":       (0, 0, 0, 0),          # fully transparent
        "hardness":    0.0,
        "transparent": True,
        "solid":       False,
        "collectible": False,
    },
    "dirt": {
        "name":        "Dirt",
        "color":       (139, 90, 43),
        "hardness":    1.0,
        "transparent": False,
        "solid":       True,
        "collectible": True,
    },
    "grass": {
        "name":        "Grass",
        "color":       (100, 200, 100),
        "hardness":    1.0,
        "transparent": False,
        "solid":       True,
        "collectible": True,
    },
    "stone": {
        "name":        "Stone",
        "color":       (169, 169, 169),
        "hardness":    2.0,
        "transparent": False,
        "solid":       True,
        "collectible": True,
    },
    "sand": {
        "name":        "Sand",
        "color":       (238, 214, 175),
        "hardness":    0.8,
        "transparent": False,
        "solid":       True,
        "collectible": True,
        "gravity":     True,           # future feature
    },
    "water": {
        "name":        "Water",
        "color":       (64, 164, 223, 140),
        "hardness":    0.0,
        "transparent": True,
        "solid":       False,
        "collectible": False,
    },
    "coal": {
        "name":        "Coal Ore",
        "color":       (40, 40, 40),
        "hardness":    2.5,
        "transparent": False,
        "solid":       True,
        "collectible": True,
    },
    "iron": {
        "name":        "Iron Ore",
        "color":       (229, 194, 159),
        "hardness":    2.8,
        "transparent": False,
        "solid":       True,
        "collectible": True,
    },
    "gold": {
        "name":        "Gold Ore",
        "color":       (255, 215, 0),
        "hardness":    3.2,
        "transparent": False,
        "solid":       True,
        "collectible": True,
    },
    "diamond": {
        "name":        "Diamond Ore",
        "color":       (0, 220, 255),
        "hardness":    4.5,
        "transparent": False,
        "solid":       True,
        "collectible": True,
    },
    "wood": {
        "name":        "Wood",
        "color":       (101, 67, 33),
        "hardness":    1.5,
        "transparent": False,
        "solid":       True,
        "collectible": True,
    },
    "glass": {
        "name":        "Glass",
        "color":       (200, 220, 255, 160),
        "hardness":    0.3,
        "transparent": True,
        "solid":       True,
        "collectible": True,
    },
    # Add more blocks here later...
}

# Quick lookup tables (generated once at import)
COLOR_BY_TYPE: Dict[str, Tuple[int, ...]] = {k: v["color"] for k, v in BLOCK_DEFINITIONS.items()}
HARDNESS_BY_TYPE = {k: v["hardness"] for k, v in BLOCK_DEFINITIONS.items()}
TRANSPARENT_BY_TYPE = {k: v["transparent"] for k, v in BLOCK_DEFINITIONS.items()}
SOLID_BY_TYPE = {k: v["solid"] for k, v in BLOCK_DEFINITIONS.items()}
COLLECTIBLE_BY_TYPE = {k: v["collectible"] for k, v in BLOCK_DEFINITIONS.items()}
