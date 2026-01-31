package organism

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// OrganismType represents different types of organisms
type OrganismType string

const (
	TreeType OrganismType = "tree"
)

// Organism represents a living entity in the game
type Organism struct {
	X              float64
	Y              float64
	Type           OrganismType
	Health         float64
	MaxHealth      float64
	Alive          bool
	HitboxRadius   float64
	DropItems      []string
	DropChances    map[string]float64
	DropAmounts    map[string][2]int // [min, max]
	AxeMultiplier  float64
	OtherMultiplier float64
}

// TreeOrg represents a tree organism
type TreeOrg struct {
	*Organism
}

// NewTree creates a new tree organism
func NewTree(x, y float64) *TreeOrg {
	return &TreeOrg{
		Organism: &Organism{
			X:              x,
			Y:              y,
			Type:           TreeType,
			Health:         100,
			MaxHealth:      100,
			Alive:          true,
			HitboxRadius:   30,
			DropItems:      []string{"wood", "sapling"},
			DropChances:    map[string]float64{"wood": 1.0, "sapling": 0.3},
			DropAmounts:    map[string][2]int{"wood": {3, 6}, "sapling": {0, 1}},
			AxeMultiplier:  5.0,
			OtherMultiplier: 0.2,
		},
	}
}

// GetDamage calculates damage based on item type
func (o *Organism) GetDamage(itemType string, baseDamage float64) float64 {
	if itemType == "axe" {
		return baseDamage * o.AxeMultiplier
	}
	return baseDamage * o.OtherMultiplier
}

// TakeDamage applies damage to the organism
func (o *Organism) TakeDamage(itemType string, baseDamage float64) []Drop {
	damage := o.GetDamage(itemType, baseDamage)
	o.Health -= damage

	if o.Health <= 0 {
		o.Alive = false
		return o.GetDrops()
	}
	return nil
}

// Drop represents an item drop
type Drop struct {
	ItemType string
	Amount   int
}

// GetDrops returns items dropped when organism dies
func (o *Organism) GetDrops() []Drop {
	drops := []Drop{}
	for _, itemType := range o.DropItems {
		chance := o.DropChances[itemType]
		if rand.Float64() < chance {
			amountRange := o.DropAmounts[itemType]
			amount := rand.Intn(amountRange[1]-amountRange[0]+1) + amountRange[0]
			if amount > 0 {
				drops = append(drops, Drop{ItemType: itemType, Amount: amount})
			}
		}
	}
	return drops
}

// Draw renders the organism to the screen
func (o *Organism) Draw(screen *ebiten.Image, offsetX, offsetY float64) {
	screenX := o.X - offsetX
	screenY := o.Y - offsetY

	if o.Type == TreeType {
		o.drawTree(screen, screenX, screenY)
	}

	// Draw health bar if damaged
	if o.Health < o.MaxHealth {
		o.drawHealthBar(screen, screenX, screenY)
	}
}

// drawTree renders the tree structure
func (o *Organism) drawTree(screen *ebiten.Image, x, y float64) {
	// Trunk colors
	trunkBottomColor := color.RGBA{101, 67, 33, 255}
	trunkTopColor := color.RGBA{139, 90, 43, 255}
	leavesColor := color.RGBA{34, 139, 34, 255}
	black := color.RGBA{0, 0, 0, 255}

	// Draw trunk bottom (larger brown circle - approximated as rect)
	ebitenutil.DrawRect(screen, x-15, y+20-15, 30, 30, trunkBottomColor)
	ebitenutil.DrawRect(screen, x-15, y+20-15, 30, 30, black) // Outline

	// Draw trunk top (smaller brown circle - approximated as rect)
	ebitenutil.DrawRect(screen, x-12, y+5-12, 24, 24, trunkTopColor)
	ebitenutil.DrawRect(screen, x-12, y+5-12, 24, 24, black) // Outline

	// Draw leaves (larger green circle - approximated as rect)
	ebitenutil.DrawRect(screen, x-25, y-15-25, 50, 50, leavesColor)
	ebitenutil.DrawRect(screen, x-25, y-15-25, 50, 50, black) // Outline
}

// drawHealthBar renders the health bar
func (o *Organism) drawHealthBar(screen *ebiten.Image, x, y float64) {
	barWidth := 40.0
	barHeight := 5.0
	barX := x - barWidth/2
	barY := y - 45

	healthRatio := o.Health / o.MaxHealth
	healthWidth := barWidth * healthRatio

	// Background bar
	ebitenutil.DrawRect(screen, barX, barY, barWidth, barHeight, color.RGBA{100, 100, 100, 255})

	// Health bar (green to red gradient)
	healthColor := color.RGBA{
		R: uint8(255 * (1 - healthRatio)),
		G: uint8(255 * healthRatio),
		B: 0,
		A: 255,
	}
	ebitenutil.DrawRect(screen, barX, barY, healthWidth, barHeight, healthColor)
}