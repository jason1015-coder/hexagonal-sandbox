package biomes

import (
	"math"
)

// BiomeType represents different biomes in the world
type BiomeType int

const (
	PLAINS BiomeType = iota
	FOREST
	DESERT
	MOUNTAINS
	OCEAN
	SWAMP
)

// BiomeProperties defines properties of a biome
type BiomeProperties struct {
	Name         string
	SurfaceBlock string
	UnderBlock   string
	TreeDensity  float64 // 0.0 to 1.0
	OreFrequency float64
	Temperature  float64
	Humidity     float64
}

// BiomeDefinitions holds all biome type definitions
var BiomeDefinitions = map[BiomeType]*BiomeProperties{
	PLAINS: {
		Name:         "Plains",
		SurfaceBlock: "grass",
		UnderBlock:   "dirt",
		TreeDensity:  0.1,
		OreFrequency: 1.0,
		Temperature:  0.5,
		Humidity:     0.5,
	},
	FOREST: {
		Name:         "Forest",
		SurfaceBlock: "grass",
		UnderBlock:   "dirt",
		TreeDensity:  0.4,
		OreFrequency: 1.0,
		Temperature:  0.4,
		Humidity:     0.7,
	},
	DESERT: {
		Name:         "Desert",
		SurfaceBlock: "sand",
		UnderBlock:   "sand",
		TreeDensity:  0.0,
		OreFrequency: 0.5,
		Temperature:  0.9,
		Humidity:     0.1,
	},
	MOUNTAINS: {
		Name:         "Mountains",
		SurfaceBlock: "stone",
		UnderBlock:   "stone",
		TreeDensity:  0.05,
		OreFrequency: 2.0,
		Temperature:  0.3,
		Humidity:     0.3,
	},
	OCEAN: {
		Name:         "Ocean",
		SurfaceBlock: "water",
		UnderBlock:   "sand",
		TreeDensity:  0.0,
		OreFrequency: 0.3,
		Temperature:  0.6,
		Humidity:     1.0,
	},
	SWAMP: {
		Name:         "Swamp",
		SurfaceBlock: "grass",
		UnderBlock:   "dirt",
		TreeDensity:  0.2,
		OreFrequency: 0.8,
		Temperature:  0.6,
		Humidity:     0.9,
	},
}

// SimplexNoise is a simple noise implementation for terrain generation
type SimplexNoise struct {
	seed float64
}

// NewSimplexNoise creates a new simplex noise generator
func NewSimplexNoise(seed float64) *SimplexNoise {
	return &SimplexNoise{seed: seed}
}

// Noise2D returns 2D noise value at the given coordinates
func (n *SimplexNoise) Noise2D(x, y float64) float64 {
	// Simple value-based noise for demonstration
	// In a full implementation, use a proper simplex/perlin noise library
	return n.sineNoise(x*0.01+n.seed, y*0.01+n.seed) * 0.5 +
		n.sineNoise(x*0.05+n.seed, y*0.05+n.seed) * 0.25 +
		n.sineNoise(x*0.1+n.seed, y*0.1+n.seed) * 0.25
}

// sineNoise generates a simple sine-based noise
func (n *SimplexNoise) sineNoise(x, y float64) float64 {
	return (math.Sin(x) + math.Cos(y)) / 2.0
}

// GetBiomeAtPosition returns the biome type at the given world coordinates
func GetBiomeAtPosition(x, y float64, noise *SimplexNoise) BiomeType {
	temp := noise.Noise2D(x*0.01, y*0.01)
	humid := noise.Noise2D(x*0.01+1000, y*0.01+1000)
	elev := noise.Noise2D(x*0.005, y*0.005)
	
	// Normalize values to 0-1 range
	temp = (temp + 1) / 2.0
	humid = (humid + 1) / 2.0
	elev = (elev + 1) / 2.0
	
	// Determine biome based on temperature, humidity, and elevation
	if elev < 0.3 {
		return OCEAN
	}
	
	if elev > 0.7 {
		return MOUNTAINS
	}
	
	if temp > 0.7 && humid < 0.3 {
		return DESERT
	}
	
	if temp > 0.3 && temp < 0.7 && humid > 0.7 {
		return SWAMP
	}
	
	if humid > 0.5 && temp > 0.3 {
		return FOREST
	}
	
	return PLAINS
}

// GetSurfaceHeightVariation returns the surface height variation at the given position
func GetSurfaceHeightVariation(x, y float64, noise *SimplexNoise) float64 {
	return noise.Noise2D(x*0.02, y*0.02) * 5.0
}

// ShouldSpawnTree returns whether a tree should spawn at the given position
func ShouldSpawnTree(x, y float64, noise *SimplexNoise) bool {
	biome := GetBiomeAtPosition(x, y, noise)
	props := BiomeDefinitions[biome]
	
	if props.TreeDensity <= 0 {
		return false
	}
	
	// Use noise to determine if tree should spawn
	treeNoise := noise.Noise2D(x*0.1, y*0.1)
	return treeNoise < props.TreeDensity
}

// GetBiomeBlock returns the surface block type for the given biome
func GetBiomeBlock(biome BiomeType) string {
	if props, ok := BiomeDefinitions[biome]; ok {
		return props.SurfaceBlock
	}
	return "dirt"
}

// GetUnderBlock returns the underground block type for the given biome
func GetBiomeUnderBlock(biome BiomeType) string {
	if props, ok := BiomeDefinitions[biome]; ok {
		return props.UnderBlock
	}
	return "stone"
}

// ShouldSpawnOre returns whether an ore should spawn at the given depth and position
func ShouldSpawnOre(depth int, x, y float64, noise *SimplexNoise, oreType string) bool {
	// Different ores spawn at different depths
	var minDepth, maxDepth, frequency float64
	
	switch oreType {
	case "coal_ore":
		minDepth, maxDepth, frequency = 5, 50, 0.05
	case "iron_ore":
		minDepth, maxDepth, frequency = 10, 60, 0.03
	case "gold_ore":
		minDepth, maxDepth, frequency = 20, 40, 0.02
	case "diamond_ore":
		minDepth, maxDepth, frequency = 30, 50, 0.01
	default:
		return false
	}
	
	depthFloat := float64(depth)
	if depthFloat < minDepth || depthFloat > maxDepth {
		return false
	}
	
	oreNoise := noise.Noise2D(x*0.05, y*0.05)
	return oreNoise < frequency
}