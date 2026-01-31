package hexagon

import (
	"math"
)

// Hexagon represents a hexagonal tile
type Hexagon struct {
	Q, R, S int // Cube coordinates
}

// NewHexagon creates a hexagon from cube coordinates
func NewHexagon(q, r, s int) Hexagon {
	if q+r+s != 0 {
		panic("cube coordinates must sum to zero")
	}
	return Hexagon{Q: q, R: r, S: s}
}

// AxialToHex creates a hexagon from axial coordinates
func AxialToHex(q, r int) Hexagon {
	return NewHexagon(q, r, -q-r)
}

// HexAdd adds two hexagons together
func HexAdd(a, b Hexagon) Hexagon {
	return NewHexagon(a.Q+b.Q, a.R+b.R, a.S+b.S)
}

// HexSubtract subtracts two hexagons
func HexSubtract(a, b Hexagon) Hexagon {
	return NewHexagon(a.Q-b.Q, a.R-b.R, a.S-b.S)
}

// HexScale scales a hexagon by a factor
func HexScale(a Hexagon, k int) Hexagon {
	return NewHexagon(a.Q*k, a.R*k, a.S*k)
}

// HexNeighbors returns the six neighboring hexagons
func HexNeighbors(hex Hexagon) []Hexagon {
	directions := []Hexagon{
		{1, 0, -1},
		{1, -1, 0},
		{0, -1, 1},
		{-1, 0, 1},
		{-1, 1, 0},
		{0, 1, -1},
	}
	
	neighbors := make([]Hexagon, 6)
	for i, dir := range directions {
		neighbors[i] = HexAdd(hex, dir)
	}
	return neighbors
}

// HexDistance calculates the distance between two hexagons
func HexDistance(a, b Hexagon) int {
	return (abs(a.Q-b.Q) + abs(a.Q+b.R-b.Q-b.R) + abs(a.R-b.R)) / 2
}

// PixelToHex converts pixel coordinates to hexagon coordinates
// Using axial coordinates: x = size * (3/2 * q), y = size * (sqrt(3)/2 * q + sqrt(3) * r)
func PixelToHex(x, y, size float64) (float64, float64) {
	q := (math.Sqrt(3)/3*x + 1.0/3.0*y) / size
	r := (-math.Sqrt(3)/3*x + 2.0/3.0*y) / size
	return q, r
}

// HexRound rounds floating point hex coordinates to nearest integer hexagon
func HexRound(q, r float64) Hexagon {
	x := q
	z := r
	y := -x - z
	
	rx := round(x)
	ry := round(y)
	rz := round(z)
	
	xDiff := absFloat64(rx - x)
	yDiff := absFloat64(ry - y)
	zDiff := absFloat64(rz - z)
	
	if xDiff > yDiff && xDiff > zDiff {
		rx = -ry - rz
	} else if yDiff > zDiff {
		ry = -rx - rz
	} else {
		rz = -rx - ry
	}
	
	return NewHexagon(int(rx), int(ry), int(rz))
}

// HexToPixel converts hexagon coordinates to pixel coordinates
func HexToPixel(hex Hexagon, size float64) (float64, float64) {
	x := size * (3.0/2.0 * float64(hex.Q))
	y := size * (math.Sqrt(3)/2.0*float64(hex.Q) + math.Sqrt(3)*float64(hex.R))
	return x, y
}

// GetHexCorners returns the pixel coordinates of a hexagon's six corners
func GetHexCorners(centerX, centerY, size float64) [][2]float64 {
	corners := make([][2]float64, 6)
	for i := 0; i < 6; i++ {
		angle := math.Pi / 180.0 * (30.0 + 60.0*float64(i))
		corners[i] = [2]float64{
			centerX + size*math.Cos(angle),
			centerY + size*math.Sin(angle),
		}
	}
	return corners
}

// ChunkID represents a chunk identifier
type ChunkID struct {
	X, Y int
}

// HexToChunk converts hex coordinates to chunk coordinates
func HexToChunk(hex Hexagon, chunkSize int) ChunkID {
	chunkX := hex.Q / chunkSize
	chunkY := hex.R / chunkSize
	
	// Handle negative coordinates
	if hex.Q < 0 && hex.Q%chunkSize != 0 {
		chunkX--
	}
	if hex.R < 0 && hex.R%chunkSize != 0 {
		chunkY--
	}
	
	return ChunkID{X: chunkX, Y: chunkY}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// round returns the nearest integer to a float64
func round(x float64) float64 {
	return math.Floor(x + 0.5)
}
// absFloat64 returns the absolute value of a float64
func absFloat64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
