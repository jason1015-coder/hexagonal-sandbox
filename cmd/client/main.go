package main

import (
	"fmt"
	"log"
	"math"
	
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	
	"github.com/ninjatech/tesselbox-go/pkg/blocks"
	"github.com/ninjatech/tesselbox-go/pkg/hexagon"
	"github.com/ninjatech/tesselbox-go/pkg/items"
	"github.com/ninjatech/tesselbox-go/pkg/player"
	"github.com/ninjatech/tesselbox-go/pkg/world"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 720
	FPS          = 60
)

// Game represents the game state
type Game struct {
	world     *world.World
	player    *player.Player
	inventory *items.Inventory
	
	cameraX, cameraY float64
	
	mouseX, mouseY int
	lastClickX, lastClickY int
	
	font map[rune]*ebiten.Image
}

// NewGame creates a new game
func NewGame() *Game {
	g := &Game{
		world:     world.NewWorld(42.0), // Random seed
		player:    player.NewPlayer(0, 0),
		inventory: items.NewInventory(32),
		cameraX:   0,
		cameraY:   0,
		font:      make(map[rune]*ebiten.Image),
	}
	
	// Set default hotbar items
	defaultItems := items.DefaultHotbarItems()
	for i, item := range defaultItems {
		if i < len(g.inventory.Slots) {
			g.inventory.Slots[i] = item
		}
	}
	
	return g
}

// Update updates the game state
func (g *Game) Update() error {
	// Handle keyboard input
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.MovingLeft = true
		g.player.MovingRight = false
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.MovingRight = true
		g.player.MovingLeft = false
	} else {
		g.player.MovingLeft = false
		g.player.MovingRight = false
	}
	
	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Jump()
	}
	
	// Handle mouse input
	g.mouseX, g.mouseY = ebiten.CursorPosition()
	
	// Mining with mouse click
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.handleMining()
	}
	
	// Block placement with right click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.handleBlockPlacement()
	}
	
	// Hotbar selection
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.inventory.SelectSlot(0)
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.inventory.SelectSlot(1)
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.inventory.SelectSlot(2)
	} else if inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.inventory.SelectSlot(3)
	} else if inpututil.IsKeyJustPressed(ebiten.Key5) {
		g.inventory.SelectSlot(4)
	} else if inpututil.IsKeyJustPressed(ebiten.Key6) {
		g.inventory.SelectSlot(5)
	} else if inpututil.IsKeyJustPressed(ebiten.Key7) {
		g.inventory.SelectSlot(6)
	} else if inpututil.IsKeyJustPressed(ebiten.Key8) {
		g.inventory.SelectSlot(7)
	} else if inpututil.IsKeyJustPressed(ebiten.Key9) {
		g.inventory.SelectSlot(8)
	}
	
	// Scroll wheel for hotbar
	_, scrollY := ebiten.Wheel()
	if scrollY > 0 {
		g.inventory.PrevSlot()
	} else if scrollY < 0 {
		g.inventory.NextSlot()
	}
	
	// Update player
	g.player.Update(1.0 / FPS)
	
	// Update camera to follow player
	centerX, centerY := g.player.GetCenter()
	g.cameraX = centerX - ScreenWidth/2
	g.cameraY = centerY - ScreenHeight/2
	
	// Generate chunks around player
	px, py := g.player.GetCenter()
	g.world.GetChunksInRange(px, py)
	
	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background
	screen.Fill(g.colorToRGB(135, 206, 235)) // Sky blue
	
	// Get visible blocks
	px, py := g.player.GetCenter()
	visibleBlocks := g.world.GetVisibleBlocks(px, py)
	
	// Draw blocks
	for _, block := range visibleBlocks {
		g.drawBlock(screen, block)
	}
	
	// Draw player
	g.drawPlayer(screen)
	
	// Draw UI
	g.drawUI(screen)
	
	// Draw debug info
	g.drawDebugInfo(screen)
}

// Layout defines the game's layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// drawBlock draws a hexagonal block
func (g *Game) drawBlock(screen *ebiten.Image, block *world.Block) {
	if block.Type == "air" {
		return
	}
	
	props := blocks.BlockDefinitions[block.Type]
	if props == nil {
		return
	}
	
	// Calculate screen position
	screenX := block.X - g.cameraX
	screenY := block.Y - g.cameraY
	
	// Check if block is on screen
	if screenX < -100 || screenX > ScreenWidth+100 ||
		screenY < -100 || screenY > ScreenHeight+100 {
		return
	}
	
	// Get hexagon corners
	corners := hexagon.GetHexCorners(screenX, screenY, world.HexSize)
	
	// Create polygon vertices
	vertices := make([]ebiten.Vertex, len(corners))
	for i, corner := range corners {
		vertices[i] = ebiten.Vertex{
			DstX:   float32(corner[0]),
			DstY:   float32(corner[1]),
			ColorR: float32(props.Color.R) / 255.0,
			ColorG: float32(props.Color.G) / 255.0,
			ColorB: float32(props.Color.B) / 255.0,
			ColorA: float32(props.Color.A) / 255.0,
		}
	}
	
	// Draw filled hexagon
	indices := []uint16{0, 1, 2, 0, 2, 3, 0, 3, 4, 0, 4, 5}
	
	// Check if mouse is hovering over this block
	mouseWorldX := float64(g.mouseX) + g.cameraX
	mouseWorldY := float64(g.mouseY) + g.cameraY
	
	q, r := hexagon.PixelToHex(mouseWorldX, mouseWorldY, world.HexSize)
	hoverHex := hexagon.HexRound(q, r)
	
	if hoverHex.Q == block.Hex.Q && hoverHex.R == block.Hex.R {
		// Highlight hovered block
		for i := range vertices {
			vertices[i].ColorR = min(1.0, vertices[i].ColorR+0.2)
			vertices[i].ColorG = min(1.0, vertices[i].ColorG+0.2)
			vertices[i].ColorB = min(1.0, vertices[i].ColorB+0.2)
		}
	}
	
	// Apply damage darkening
	if block.Health < block.MaxHealth {
		damageRatio := block.Health / block.MaxHealth
		for i := range vertices {
			vertices[i].ColorR *= float32(damageRatio)
			vertices[i].ColorG *= float32(damageRatio)
			vertices[i].ColorB *= float32(damageRatio)
		}
	}
	
	screen.DrawTriangles(vertices, indices, ebiten.NewImageFromImage(nil), nil)
}

// drawPlayer draws the player
func (g *Game) drawPlayer(screen *ebiten.Image) {
	screenX := g.player.X - g.cameraX
	screenY := g.player.Y - g.cameraY
	
	// Draw player body
	ebitenutil.DrawRect(screen, screenX, screenY, g.player.Width, g.player.Height, g.colorToRGB(255, 100, 100))
	
	// Draw player head (simple representation)
	ebitenutil.DrawRect(screen, screenX+10, screenY+5, 20, 20, g.colorToRGB(255, 200, 150))
}

// drawUI draws the user interface
func (g *Game) drawUI(screen *ebiten.Image) {
	// Draw hotbar
	hotbarWidth := 400
	hotbarHeight := 50
	hotbarX := (ScreenWidth - hotbarWidth) / 2
	hotbarY := ScreenHeight - hotbarHeight - 20
	
	slotWidth := hotbarWidth / 8
	
	for i := 0; i < 8; i++ {
		slotX := hotbarX + i*slotWidth
		slotY := hotbarY
		
		// Draw slot background
		bgColor := g.colorToRGB(100, 100, 100)
		if i == g.inventory.Selected {
			bgColor = g.colorToRGB(150, 150, 150)
		}
		
		ebitenutil.DrawRect(screen, slotX, slotY, slotWidth-2, hotbarHeight, bgColor)
		
		// Draw slot border
		ebitenutil.DrawRect(screen, slotX, slotY, slotWidth-2, 2, g.colorToRGB(0, 0, 0))
		ebitenutil.DrawRect(screen, slotX, slotY+hotbarHeight-2, slotWidth-2, 2, g.colorToRGB(0, 0, 0))
		ebitenutil.DrawRect(screen, slotX, slotY, 2, hotbarHeight, g.colorToRGB(0, 0, 0))
		ebitenutil.DrawRect(screen, slotX+slotWidth-4, slotY, 2, hotbarHeight, g.colorToRGB(0, 0, 0))
		
		// Draw item if present
		if i < len(g.inventory.Slots) {
			item := g.inventory.Slots[i]
			if item.Type != items.NONE {
				itemColor := items.ItemColorByID(item.Type)
				ebitenutil.DrawRect(screen, slotX+10, slotY+10, slotWidth-20, hotbarHeight-20, 
					g.colorToRGB(int(itemColor.R), int(itemColor.G), int(itemColor.B)))
				
				// Draw quantity
				if item.Quantity > 1 {
					// Simple text representation would go here
					// For now, just draw a small dot to indicate multiple items
					if item.Quantity > 1 {
						ebitenutil.DrawRect(screen, slotX+slotWidth-20, slotY+hotbarHeight-20, 8, 8, g.colorToRGB(255, 255, 255))
					}
				}
			}
		}
	}
	
	// Draw health bar
	healthBarWidth := 200
	healthBarHeight := 20
	healthBarX := 20
	healthBarY := 20
	
	healthRatio := g.player.Health / g.player.MaxHealth
	
	// Background
	ebitenutil.DrawRect(screen, healthBarX, healthBarY, healthBarWidth, healthBarHeight, g.colorToRGB(50, 50, 50))
	
	// Health
	ebitenutil.DrawRect(screen, healthBarX, healthBarY, int(float64(healthBarWidth)*healthRatio), healthBarHeight, g.colorToRGB(200, 50, 50))
	
	// Border
	ebitenutil.DrawRect(screen, healthBarX, healthBarY, healthBarWidth, 2, g.colorToRGB(0, 0, 0))
	ebitenutil.DrawRect(screen, healthBarX, healthBarY+healthBarHeight-2, healthBarWidth, 2, g.colorToRGB(0, 0, 0))
	ebitenutil.DrawRect(screen, healthBarX, healthBarY, 2, healthBarHeight, g.colorToRGB(0, 0, 0))
	ebitenutil.DrawRect(screen, healthBarX+healthBarWidth-2, healthBarY, 2, healthBarHeight, g.colorToRGB(0, 0, 0))
}

// drawDebugInfo draws debug information
func (g *Game) drawDebugInfo(screen *ebiten.Image) {
	px, py := g.player.GetCenter()
	vx, vy := g.player.GetVelocity()
	
	info := fmt.Sprintf("Pos: (%.1f, %.1f)\nVel: (%.1f, %.1f)\nFPS: %.1f\nOnGround: %v",
		px, py, vx, vy, ebiten.ActualFPS(), g.player.IsOnGround())
	
	ebitenutil.DebugPrint(screen, info)
}

// handleMining handles block mining
func (g *Game) handleMining() {
	// Convert mouse position to world coordinates
	mouseWorldX := float64(g.mouseX) + g.cameraX
	mouseWorldY := float64(g.mouseY) + g.cameraY
	
	// Convert to hex coordinates
	q, r := hexagon.PixelToHex(mouseWorldX, mouseWorldY, world.HexSize)
	targetHex := hexagon.HexRound(q, r)
	
	// Find depth (simplified - assume top-most block)
	// In a full implementation, you'd need to raycast to find the correct depth
	depth := 0
	
	// Check if player can reach
	if !g.player.CanReach(mouseWorldX, mouseWorldY) {
		return
	}
	
	// Damage block
	destroyed := g.world.DamageBlock(targetHex, depth, 5.0) // 5 damage per frame
	
	if destroyed {
		// Use item durability
		g.inventory.UseItem()
	}
}

// handleBlockPlacement handles block placement
func (g *Game) handleBlockPlacement() {
	// Get selected item
	selectedItem := g.inventory.GetSelectedItem()
	if selectedItem == nil || selectedItem.Type == items.NONE {
		return
	}
	
	// Convert mouse position to world coordinates
	mouseWorldX := float64(g.mouseX) + g.cameraX
	mouseWorldY := float64(g.mouseY) + g.cameraY
	
	// Convert to hex coordinates
	q, r := hexagon.PixelToHex(mouseWorldX, mouseWorldY, world.HexSize)
	targetHex := hexagon.HexRound(q, r)
	
	// Determine block type from item
	var blockType string
	switch selectedItem.Type {
	case items.DIRT_BLOCK:
		blockType = "dirt"
	case items.GRASS_BLOCK:
		blockType = "grass"
	case items.STONE_BLOCK:
		blockType = "stone"
	case items.SAND_BLOCK:
		blockType = "sand"
	case items.LOG_BLOCK:
		blockType = "log"
	default:
		return // Not a placeable block
	}
	
	// Place block (simplified - at depth 0)
	depth := 0
	g.world.SetBlock(targetHex, depth, blockType)
	
	// Remove item from inventory
	g.inventory.RemoveItem(1)
}

// colorToRGB converts RGB values to a color
func (g *Game) colorToRGB(r, g, b int) ebiten.Color {
	return ebiten.Color{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Main function
func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Tesselbox Go")
	ebiten.SetTPS(FPS)
	
	game := NewGame()
	
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}