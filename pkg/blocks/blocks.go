package blocks

import "image/color"

// BlockType represents the type of a block
type BlockType int

const (
	AIR BlockType = iota
	DIRT
	GRASS
	STONE
	SAND
	WATER
	LOG
	LEAVES
	COAL_ORE
	IRON_ORE
	GOLD_ORE
	DIAMOND_ORE
	BEDROCK
	GLASS
	BRICK
	PLANK
	CACTUS
)

// BlockProperties defines the properties of a block type
type BlockProperties struct {
	ID          BlockType
	Name        string
	Color       color.RGBA
	TopColor    color.RGBA // For blocks with different top face colors
	SideColor   color.RGBA // For blocks with different side face colors
	Hardness    float64
	Transparent bool
	Solid       bool
	Collectible bool
	Flammable   bool
	LightLevel  int
	Gravity     bool
	Viscosity   float64 // For liquids
}

// BlockDefinitions holds all block type definitions
var BlockDefinitions = map[string]*BlockProperties{
	"air": {
		ID:          AIR,
		Name:        "Air",
		Color:       color.RGBA{0, 0, 0, 0},
		Hardness:    0,
		Transparent: true,
		Solid:       false,
		Collectible: false,
		Flammable:   false,
		LightLevel:  0,
	},
	"dirt": {
		ID:          DIRT,
		Name:        "Dirt",
		Color:       color.RGBA{139, 90, 43, 255},
		Hardness:    1.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	"grass": {
		ID:          GRASS,
		Name:        "Grass Block",
		Color:       color.RGBA{100, 200, 100, 255},
		TopColor:    color.RGBA{126, 200, 80, 255},
		SideColor:   color.RGBA{95, 150, 30, 255},
		Hardness:    1.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	"stone": {
		ID:          STONE,
		Name:        "Stone",
		Color:       color.RGBA{169, 169, 169, 255},
		Hardness:    2.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	"sand": {
		ID:          SAND,
		Name:        "Sand",
		Color:       color.RGBA{238, 214, 175, 255},
		Hardness:    0.8,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Gravity:     true,
		Flammable:   false,
		LightLevel:  0,
	},
	"water": {
		ID:          WATER,
		Name:        "Water",
		Color:       color.RGBA{64, 164, 223, 140},
		Hardness:    0,
		Transparent: true,
		Solid:       false,
		Collectible: false,
		Viscosity:   0.8,
		LightLevel:  1,
	},
	"log": {
		ID:          LOG,
		Name:        "Log",
		Color:       color.RGBA{139, 69, 19, 255},
		Hardness:    2.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   true,
		LightLevel:  0,
	},
	"leaves": {
		ID:          LEAVES,
		Name:        "Leaves",
		Color:       color.RGBA{34, 139, 34, 200},
		Hardness:    0.5,
		Transparent: true,
		Solid:       true,
		Collectible: true,
		Flammable:   true,
		LightLevel:  0,
	},
	"coal_ore": {
		ID:          COAL_ORE,
		Name:        "Coal Ore",
		Color:       color.RGBA{45, 45, 45, 255},
		Hardness:    3.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	"iron_ore": {
		ID:          IRON_ORE,
		Name:        "Iron Ore",
		Color:       color.RGBA{169, 166, 150, 255},
		Hardness:    3.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	"gold_ore": {
		ID:          GOLD_ORE,
		Name:        "Gold Ore",
		Color:       color.RGBA{255, 215, 0, 255},
		Hardness:    3.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	"diamond_ore": {
		ID:          DIAMOND_ORE,
		Name:        "Diamond Ore",
		Color:       color.RGBA{0, 255, 255, 255},
		Hardness:    3.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	"bedrock": {
		ID:          BEDROCK,
		Name:        "Bedrock",
		Color:       color.RGBA{30, 30, 30, 255},
		Hardness:    -1, // Unbreakable
		Transparent: false,
		Solid:       true,
		Collectible: false,
		Flammable:   false,
		LightLevel:  0,
	},
	"glass": {
		ID:          GLASS,
		Name:        "Glass",
		Color:       color.RGBA{200, 220, 255, 100},
		Hardness:    0.5,
		Transparent: true,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	"brick": {
		ID:          BRICK,
		Name:        "Brick",
		Color:       color.RGBA{178, 34, 34, 255},
		Hardness:    2.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	"plank": {
		ID:          PLANK,
		Name:        "Wooden Plank",
		Color:       color.RGBA{222, 184, 135, 255},
		Hardness:    1.5,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   true,
		LightLevel:  0,
	},
	"cactus": {
		ID:          CACTUS,
		Name:        "Cactus",
		Color:       color.RGBA{34, 139, 34, 255},
		Hardness:    0.5,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
}

// ColorByType returns the color for a block type string
func ColorByType(blockType string) color.RGBA {
	if props, ok := BlockDefinitions[blockType]; ok {
		return props.Color
	}
	return BlockDefinitions["dirt"].Color
}

// HardnessByType returns the hardness for a block type string
func HardnessByType(blockType string) float64 {
	if props, ok := BlockDefinitions[blockType]; ok {
		return props.Hardness
	}
	return 1.0
}

// TransparentByType returns if a block type is transparent
func TransparentByType(blockType string) bool {
	if props, ok := BlockDefinitions[blockType]; ok {
		return props.Transparent
	}
	return false
}

// SolidByType returns if a block type is solid
func SolidByType(blockType string) bool {
	if props, ok := BlockDefinitions[blockType]; ok {
		return props.Solid
	}
	return true
}

// CollectibleByType returns if a block type is collectible
func CollectibleByType(blockType string) bool {
	if props, ok := BlockDefinitions[blockType]; ok {
		return props.Collectible
	}
	return true
}

// ValidBlockType checks if a block type string is valid
func ValidBlockType(blockType string) bool {
	_, ok := BlockDefinitions[blockType]
	return ok
}