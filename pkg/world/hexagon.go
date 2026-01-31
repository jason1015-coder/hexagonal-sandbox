package world

import (
	"image/color"
	"math"

	"github.com/ninjatech/tesselbox-go/pkg/block"
)

var (
	// HexSize is the radius of a hexagon
	HexSize = 30.0
	// HexWidth is the width of a hexagon
	HexWidth = math.Sqrt(3) * HexSize
	// HexHeight is the height of a hexagon
	HexHeight = 2.0 * HexSize
	// HexVSpacing is the vertical spacing between hexagon rows
	HexVSpacing = HexHeight * 0.75
)

// Hexagon represents a single hexagonal block in the world
type Hexagon struct {
	X            float64
	Y            float64
	Size         float64
	BlockType    block.BlockType
	Color        color.Color
	ActiveColor  color.Color
	Hovered      bool
	Health       float64
	MaxHealth    float64
	Transparent  bool
	Corners      [][2]float64
	ChunkX       int
	ChunkY       int
}

// NewHexagon creates a new hexagon at the specified position
func NewHexagon(x, y, size float64, blockType block.BlockType) *Hexagon {
	def := block.GetDefinition(blockType)

	h := &Hexagon{
		X:           x,
		Y:           y,
		Size:        size,
		BlockType:   blockType,
		Color:       def.Color,
		ActiveColor: def.Color,
		Hovered:     false,
		Health:      100,
		MaxHealth:   100,
		Transparent: def.Transparent,
		Corners:     make([][2]float64, 6),
	}

	h.calculateCorners()
	return h
}

// calculateCorners calculates the corner positions of the hexagon
func (h *Hexagon) calculateCorners() {
	for i := 0; i < 6; i++ {
		angle := math.Pi/6 + math.Pi/3*float64(i)
		px := h.X + h.Size*math.Cos(angle)
		py := h.Y + h.Size*math.Sin(angle)
		h.Corners[i] = [2]float64{px, py}
	}
}

// CheckHover determines if the hexagon is being hovered by the mouse
func (h *Hexagon) CheckHover(mouseX, mouseY, playerX, playerY, miningRange float64) {
	dx := mouseX - h.X
	dy := mouseY - h.Y
	distanceSq := dx*dx + dy*dy

	pdx := playerX - h.X
	pdy := playerY - h.Y
	playerDistanceSq := pdx*pdx + pdy*pdy
	inRange := playerDistanceSq < miningRange*miningRange

	hexRadiusSq := (h.Size * 0.866) * (h.Size * 0.866)
	h.Hovered = distanceSq < hexRadiusSq && inRange

	if h.Hovered {
		h.ActiveColor = h.brightenColor(h.Color, 30)
	} else {
		h.ActiveColor = h.Color
	}
}

// brightenColor increases the brightness of a color
func (h *Hexagon) brightenColor(c color.Color, amount int) color.Color {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(min(255, int(r/257)+amount)),
		G: uint8(min(255, int(g/257)+amount)),
		B: uint8(min(255, int(b/257)+amount)),
		A: uint8(a / 257),
	}
}

// TakeDamage applies damage to the hexagon
func (h *Hexagon) TakeDamage(amount float64) bool {
	def := block.GetDefinition(h.BlockType)
	hardness := def.Hardness
	if hardness <= 0 {
		return false
	}
	h.Health -= amount / hardness
	return h.Health <= 0
}

// GetTopSurfaceY returns the Y coordinate of the top surface at a given X
func (h *Hexagon) GetTopSurfaceY(x float64) float64 {
	left := h.Corners[3]
	top := h.Corners[4]
	right := h.Corners[5]

	if left[0] <= x && x <= top[0] {
		if left[0] != top[0] {
			t := (x - left[0]) / (top[0] - left[0])
			return left[1] + t*(top[1]-left[1])
		}
		return left[1]
	} else if top[0] <= x && x <= right[0] {
		if top[0] != right[0] {
			t := (x - top[0]) / (right[0] - top[0])
			return top[1] + t*(right[1]-top[1])
		}
		return top[1]
	}

	return h.Y
}

// PixelToHexCenter converts pixel coordinates to hexagon center coordinates
func PixelToHexCenter(wx, wy float64) (centerX, centerY, col, row float64) {
	hexSize := HexSize
	q := (math.Sqrt(3)/3*wx + 1/3*wy) / hexSize
	r := (-math.Sqrt(3)/3*wx + 2/3*wy) / hexSize

	// Cube rounding
	x := q
	z := r
	y := -x - z

	rx := math.Round(x)
	ry := math.Round(y)
	rz := math.Round(z)

	xDiff := math.Abs(rx - x)
	yDiff := math.Abs(ry - y)
	zDiff := math.Abs(rz - z)

	if xDiff > yDiff && xDiff > zDiff {
		rx = -ry - rz
	} else if yDiff > zDiff {
		ry = -rx - rz
	} else {
		rz = -rx - ry
	}

	// Convert cube coordinates back to axial
	q = rx
	r = rz

	// Calculate pixel center
	centerX = hexSize * (3/2*q)
	centerY = hexSize * (math.Sqrt(3)/2*q+math.Sqrt(3)*r)

	col = q
	row = r

	return centerX, centerY, col, row
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}