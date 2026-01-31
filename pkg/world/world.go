package world

import (
	"math"
	"math/rand"
	"time"

	"github.com/ninjatech/tesselbox-go/pkg/block"
	"github.com/ninjatech/tesselbox-go/pkg/organisms"
)

const (
	// RenderDistance is the number of chunks to render around the player
	RenderDistance = 4
	// ChunkUnloadDistance is the distance at which chunks are unloaded
	ChunkUnloadDistance = 10
)

// World represents the game world
type World struct {
	Chunks    map[[2]int]*Chunk
	Seed      int64
	Organisms []*organisms.Organism
}

// NewWorld creates a new world
func NewWorld() *World {
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	return &World{
		Chunks:    make(map[[2]int]*Chunk),
		Seed:      seed,
		Organisms: []*organisms.Organism{},
	}
}

// GetChunkCoords returns the chunk coordinates for a given world position
func (w *World) GetChunkCoords(x, y float64) (int, int) {
	chunkX := int(math.Floor(x / GetChunkWidth()))
	chunkY := int(math.Floor(y / GetChunkHeight()))
	return chunkX, chunkY
}

// GetChunk returns the chunk at the given chunk coordinates
func (w *World) GetChunk(chunkX, chunkY int) *Chunk {
	key := [2]int{chunkX, chunkY}
	chunk, exists := w.Chunks[key]
	if !exists {
		chunk = NewChunk(chunkX, chunkY)
		w.Chunks[key] = chunk
		w.generateChunk(chunk)
	}
	chunk.LastAccessed = time.Now()
	return chunk
}

// generateChunk generates terrain for a chunk
func (w *World) generateChunk(chunk *Chunk) {
	worldX, worldY := chunk.GetWorldPosition()

	for row := 0; row < ChunkSize; row++ {
		for col := 0; col < ChunkSize; col++ {
			var x, y float64

			// Calculate hexagon position with interlocking pattern
			if row%2 == 0 {
				x = worldX + float64(col)*HexWidth + HexWidth/2
			} else {
				x = worldX + float64(col)*HexWidth + HexWidth
			}
			y = worldY + float64(row)*HexVSpacing + HexSize

			// Simple terrain generation
			noise := math.Sin(x*0.01) * math.Cos(y*0.01) * 50

			// Determine block type based on height
			var blockType block.BlockType
			if y > 300+noise {
				blockType = block.Grass
			} else if y > 280+noise {
				blockType = block.Dirt
			} else {
				blockType = block.Stone
			}

			// Create hexagon
			hexagon := NewHexagon(x, y, HexSize, blockType)
			chunk.AddHexagon(x, y, hexagon)
		}
	}

	chunk.Modified = false
}

// GetNearbyHexagons returns hexagons within a radius of the given position
func (w *World) GetNearbyHexagons(x, y, radius float64) []*Hexagon {
	chunkX, chunkY := w.GetChunkCoords(x, y)
	hexagons := []*Hexagon{}

	// Check chunks around the position
	chunkRadius := int(math.Ceil(radius / GetChunkWidth()))
	for dx := -chunkRadius; dx <= chunkRadius; dx++ {
		for dy := -chunkRadius; dy <= chunkRadius; dy++ {
			chunk := w.GetChunk(chunkX+dx, chunkY+dy)
			for _, hexagon := range chunk.Hexagons {
				hx := hexagon.X - x
				hy := hexagon.Y - y
				if hx*hx+hy*hy <= radius*radius {
					hexagons = append(hexagons, hexagon)
				}
			}
		}
	}

	return hexagons
}

// GetHexagonAt returns the hexagon at the given world position
func (w *World) GetHexagonAt(x, y float64) *Hexagon {
	chunkX, chunkY := w.GetChunkCoords(x, y)
	chunk := w.GetChunk(chunkX, chunkY)
	return chunk.GetHexagon(x, y)
}

// AddHexagonAt adds a hexagon at the given world position
func (w *World) AddHexagonAt(x, y float64, blockType block.BlockType) {
	centerX, centerY, _, _ := PixelToHexCenter(x, y)
	hexagon := NewHexagon(centerX, centerY, HexSize, blockType)

	chunkX, chunkY := w.GetChunkCoords(centerX, centerY)
	chunk := w.GetChunk(chunkX, chunkY)
	chunk.AddHexagon(centerX, centerY, hexagon)
}

// RemoveHexagonAt removes the hexagon at the given world position
func (w *World) RemoveHexagonAt(x, y float64) bool {
	chunkX, chunkY := w.GetChunkCoords(x, y)
	chunk := w.GetChunk(chunkX, chunkY)
	return chunk.RemoveHexagon(x, y)
}

// UnloadDistantChunks unloads chunks that are far from the player
func (w *World) UnloadDistantChunks(playerX, playerY float64) {
	playerChunkX, playerChunkY := w.GetChunkCoords(playerX, playerY)
	toDelete := [][2]int{}

	for key, chunk := range w.Chunks {
		dx := chunk.ChunkX - playerChunkX
		dy := chunk.ChunkY - playerChunkY
		distance := math.Sqrt(float64(dx*dx + dy*dy))

		if distance > ChunkUnloadDistance {
			toDelete = append(toDelete, key)
		}
	}

	for _, key := range toDelete {
		delete(w.Chunks, key)
	}
}

// GetNearbyOrganisms returns organisms within a radius of the given position
func (w *World) GetNearbyOrganisms(x, y, radius float64) []*organisms.Organism {
	nearby := []*organisms.Organism{}
	radiusSq := radius * radius

	for _, org := range w.Organisms {
		dx := org.X - x
		dy := org.Y - y
		if dx*dx+dy*dy <= radiusSq {
			nearby = append(nearby, org)
		}
	}

	return nearby
}

// GetOrganismAt returns the organism at the given position within tolerance
func (w *World) GetOrganismAt(x, y, tolerance float64) *organisms.Organism {
	toleranceSq := tolerance * tolerance

	for _, org := range w.Organisms {
		dx := org.X - x
		dy := org.Y - y
		if dx*dx+dy*dy <= toleranceSq {
			return org
		}
	}

	return nil
}

// RemoveOrganism removes an organism from the world
func (w *World) RemoveOrganism(x, y float64) {
	for i, org := range w.Organisms {
		if org.X == x && org.Y == y {
			w.Organisms = append(w.Organisms[:i], w.Organisms[i+1:]...)
			return
		}
	}
}
