package organisms

import (
	"github.com/ninjatech/tesselbox-go/pkg/hexagon"
)

// OrganismType represents the type of organism
type OrganismType int

const (
	TREE OrganismType = iota
	BUSH
	FLOWER
)

// Organism represents a living organism in the world
type Organism struct {
	ID        string
	Type      OrganismType
	X, Y      float64
	Hex       hexagon.Hexagon
	Health    float64
	MaxHealth float64
}

// Tree represents a tree organism
type Tree struct {
	Organism
	LogCount    int
	LeafCount   int
	HasFruit    bool
}

// Bush represents a bush organism
type Bush struct {
	Organism
	LeafType    string // "oak", "birch", etc.
	BerryCount  int
}

// Flower represents a flower organism
type Flower struct {
	Organism
	FlowerType  string // "red", "yellow", "blue", etc.
}

// CreateOrganism creates a new organism based on type
func CreateOrganism(orgType OrganismType, x, y float64, hex hexagon.Hexagon) *Organism {
	switch orgType {
	case TREE:
		return createTree(x, y, hex)
	case BUSH:
		return createBush(x, y, hex)
	case FLOWER:
		return createFlower(x, y, hex)
	default:
		return nil
	}
}

// createTree creates a new tree
func createTree(x, y float64, hex hexagon.Hexagon) *Organism {
	tree := &Tree{
		Organism: Organism{
			ID:        generateID(),
			Type:      TREE,
			X:         x,
			Y:         y,
			Hex:       hex,
			Health:    100.0,
			MaxHealth: 100.0,
		},
		LogCount:  5,  // Height of the tree
		LeafCount: 20, // Number of leaf blocks
		HasFruit:  false,
	}
	return &tree.Organism
}

// createBush creates a new bush
func createBush(x, y float64, hex hexagon.Hexagon) *Organism {
	bush := &Bush{
		Organism: Organism{
			ID:        generateID(),
			Type:      BUSH,
			X:         x,
			Y:         y,
			Hex:       hex,
			Health:    30.0,
			MaxHealth: 30.0,
		},
		LeafType:   "oak",
		BerryCount: 3,
	}
	return &bush.Organism
}

// createFlower creates a new flower
func createFlower(x, y float64, hex hexagon.Hexagon) *Organism {
	flower := &Flower{
		Organism: Organism{
			ID:        generateID(),
			Type:      FLOWER,
			X:         x,
			Y:         y,
			Hex:       hex,
			Health:    10.0,
			MaxHealth: 10.0,
		},
		FlowerType: "red",
	}
	return &flower.Organism
}

// GetOrganismBlocks returns the blocks that make up an organism
func GetOrganismBlocks(org *Organism) []hexagon.Hexagon {
	blocks := []hexagon.Hexagon{}
	
	switch org.Type {
	case TREE:
		// Tree trunk (vertical)
		for i := 0; i < 5; i++ {
			hex := hexagon.AxialToHex(org.Hex.Q, org.Hex.R-i)
			blocks = append(blocks, hex)
		}
		// Leaves (around the top)
		topHex := hexagon.AxialToHex(org.Hex.Q, org.Hex.R-5)
		neighbors := hexagon.HexNeighbors(topHex)
		blocks = append(blocks, neighbors...)
		// More leaves layer
		topHex2 := hexagon.AxialToHex(org.Hex.Q, org.Hex.R-4)
		neighbors2 := hexagon.HexNeighbors(topHex2)
		blocks = append(blocks, neighbors2...)
		
	case BUSH:
		// Bush is a single block
		blocks = append(blocks, org.Hex)
		
	case FLOWER:
		// Flower is a single block
		blocks = append(blocks, org.Hex)
	}
	
	return blocks
}

// TakeDamage damages an organism
func (org *Organism) TakeDamage(amount float64) bool {
	org.Health -= amount
	if org.Health <= 0 {
		org.Health = 0
		return true // Organism destroyed
	}
	return false
}

// Heal heals an organism
func (org *Organism) Heal(amount float64) {
	org.Health += amount
	if org.Health > org.MaxHealth {
		org.Health = org.MaxHealth
	}
}

// IsAlive returns true if the organism is still alive
func (org *Organism) IsAlive() bool {
	return org.Health > 0
}

// GetPosition returns the organism's position
func (org *Organism) GetPosition() (float64, float64) {
	return org.X, org.Y
}

// GetHex returns the organism's hex position
func (org *Organism) GetHex() hexagon.Hexagon {
	return org.Hex
}

// generateID generates a unique ID for an organism
func generateID() string {
	// Simple ID generation - in production use UUID
	return "org_" + randomString(8)
}

// randomString generates a random string of the given length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[i%len(charset)] // Simplified for deterministic output
	}
	return string(result)
}

// GetDrops returns the items dropped when an organism is destroyed
func GetDrops(org *Organism) []string {
	drops := []string{}
	
	switch org.Type {
	case TREE:
		drops = append(drops, "log")
		drops = append(drops, "leaves")
	case BUSH:
		// Bush drops berries sometimes
		drops = append(drops, "leaves")
	case FLOWER:
		drops = append(drops, "flower")
	}
	
	return drops
}