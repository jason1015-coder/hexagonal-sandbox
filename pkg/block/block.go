package block

import "image/color"

// BlockType represents different types of blocks in the game
type BlockType int

const (
	Air BlockType = iota
	Dirt
	Grass
	Stone
	Sand
	Water
	Log
	Leaves
	CoalOre
	IronOre
	GoldOre
	DiamondOre
	Bedrock
	Glass
	Brick
	Plank
	Cactus
)

// BlockDefinition holds the properties of a block type
type BlockDefinition struct {
	Name        string
	Color       color.Color
	TopColor    color.Color
	SideColor   color.Color
	Hardness    float64
	Transparent bool
	Solid       bool
	Collectible bool
	Flammable   bool
	LightLevel  int
}

// Block definitions map
var BlockDefinitions = map[BlockType]BlockDefinition{
	Air: {
		Name:        "Air",
		Color:       color.RGBA{0, 0, 0, 0},
		Hardness:    0.0,
		Transparent: true,
		Solid:       false,
		Collectible: false,
		Flammable:   false,
		LightLevel:  0,
	},
	Dirt: {
		Name:        "Dirt",
		Color:       color.RGBA{139, 90, 43, 255},
		Hardness:    1.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	Grass: {
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
	Stone: {
		Name:        "Stone",
		Color:       color.RGBA{169, 169, 169, 255},
		Hardness:    2.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	Sand: {
		Name:        "Sand",
		Color:       color.RGBA{238, 214, 175, 255},
		Hardness:    0.8,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	Water: {
		Name:        "Water",
		Color:       color.RGBA{64, 164, 223, 140},
		Hardness:    0.0,
		Transparent: true,
		Solid:       false,
		Collectible: false,
		Flammable:   false,
		LightLevel:  1,
	},
	Log: {
		Name:        "Log",
		Color:       color.RGBA{101, 67, 33, 255},
		Hardness:    2.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   true,
		LightLevel:  0,
	},
	Leaves: {
		Name:        "Leaves",
		Color:       color.RGBA{34, 139, 34, 255},
		Hardness:    0.5,
		Transparent: true,
		Solid:       true,
		Collectible: true,
		Flammable:   true,
		LightLevel:  0,
	},
	CoalOre: {
		Name:        "Coal Ore",
		Color:       color.RGBA{45, 45, 45, 255},
		Hardness:    3.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	IronOre: {
		Name:        "Iron Ore",
		Color:       color.RGBA{180, 150, 140, 255},
		Hardness:    3.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	GoldOre: {
		Name:        "Gold Ore",
		Color:       color.RGBA{250, 238, 77, 255},
		Hardness:    3.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	DiamondOre: {
		Name:        "Diamond Ore",
		Color:       color.RGBA{101, 240, 213, 255},
		Hardness:    4.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  2,
	},
	Bedrock: {
		Name:        "Bedrock",
		Color:       color.RGBA{85, 85, 85, 255},
		Hardness:    -1.0,
		Transparent: false,
		Solid:       true,
		Collectible: false,
		Flammable:   false,
		LightLevel:  0,
	},
	Glass: {
		Name:        "Glass",
		Color:       color.RGBA{200, 220, 240, 100},
		Hardness:    0.5,
		Transparent: true,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	Brick: {
		Name:        "Bricks",
		Color:       color.RGBA{178, 34, 34, 255},
		Hardness:    2.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
	Plank: {
		Name:        "Oak Planks",
		Color:       color.RGBA{199, 164, 96, 255},
		Hardness:    1.0,
		Transparent: false,
		Solid:       true,
		Collectible: true,
		Flammable:   true,
		LightLevel:  0,
	},
	Cactus: {
		Name:        "Cactus",
		Color:       color.RGBA{88, 124, 39, 255},
		Hardness:    0.8,
		Transparent: true,
		Solid:       true,
		Collectible: true,
		Flammable:   false,
		LightLevel:  0,
	},
}

// GetDefinition returns the definition for a block type
func GetDefinition(blockType BlockType) BlockDefinition {
	if def, ok := BlockDefinitions[blockType]; ok {
		return def
	}
	return BlockDefinitions[Dirt]
}