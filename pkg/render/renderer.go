package render

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/ninjatech/tesselbox-go/pkg/player"
	"github.com/ninjatech/tesselbox-go/pkg/world"
)

var (
	// Colors
	White      = color.RGBA{255, 255, 255, 255}
	Black      = color.RGBA{0, 0, 0, 255}
	DarkGray   = color.RGBA{64, 64, 64, 255}
	Gray       = color.RGBA{128, 128, 128, 255}
	Green      = color.RGBA{100, 200, 100, 255}
	Red        = color.RGBA{200, 100, 100, 255}
	Blue       = color.RGBA{50, 150, 255, 255}
)

const (
	// Screen dimensions
	ScreenWidth  = 1280
	ScreenHeight = 720
	// Default render distance
	DefaultRenderDistance = 4
)

// Game represents the main game state
type Game struct {
	ScreenWidth    int
	ScreenHeight   int
	Player         *player.Player
	World          *world.World
	CameraX        float64
	CameraY        float64
	Mining         bool
	renderDistance int
	Particles      []*Particle
}

// Particle represents a visual particle effect
type Particle struct {
	X        float64
	Y        float64
	VelX     float64
	VelY     float64
	Color    color.Color
	Lifetime int
	Age      int
	Size     int
}

// NewGame creates a new game instance
func NewGame() *Game {
	return &Game{
		ScreenWidth:    ScreenWidth,
		ScreenHeight:   ScreenHeight,
		World:          world.NewWorld(),
		CameraX:        0,
		CameraY:        0,
		Mining:         false,
		renderDistance: world.RenderDistance,
		Particles:      []*Particle{},
	}
}

// Layout implements ebiten.Game interface
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.ScreenWidth, g.ScreenHeight
}

// Update implements ebiten.Game interface
func (g *Game) Update() error {
	// Spawn player if not exists
	if g.Player == nil {
		g.spawnPlayer()
	}

	// Handle input
	g.handleInput()

	// Update camera
	g.CameraX = g.Player.X - float64(g.ScreenWidth)/2
	g.CameraY = g.Player.Y - float64(g.ScreenHeight)/2

	// Update player movement
	g.Player.MovingLeft = ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft)
	g.Player.MovingRight = ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight)
	
	// TODO: Add collision detection with nearby hexagons
	// nearbyHexagons := g.World.GetNearbyHexagons(g.Player.X, g.Player.Y, 500)
	g.Player.Update(1.0 / 60.0)

	// Update hover states
	mouseX, mouseY := ebiten.CursorPosition()
	worldMX := float64(mouseX) + g.CameraX
	worldMY := float64(mouseY) + g.CameraY
	
	hoverNearby := g.World.GetNearbyHexagons(worldMX, worldMY, world.HexSize*3)
	for _, hexagon := range hoverNearby {
		hexagon.CheckHover(worldMX, worldMY, g.Player.X, g.Player.Y, player.MiningRange)
	}

	// Handle mining
	if g.Mining {
		g.mineBlock()
	}

	// Update particles
	g.updateParticles()

	// Unload distant chunks
	g.World.UnloadDistantChunks(g.Player.X, g.Player.Y)

	return nil
}

// Draw implements ebiten.Game interface
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw sky gradient
	g.drawSkyGradient(screen)

	// Get visible hexagons
	visibleHexagons := g.getVisibleHexagons()

	// Sort by Y position for proper depth rendering
	g.sortHexagonsByY(visibleHexagons)

	// Draw hexagons
	for _, hexagon := range visibleHexagons {
		g.drawHexagon(screen, hexagon)
	}

	// Draw visible organisms
	visibleOrganisms := g.World.GetNearbyOrganisms(g.Player.X, g.Player.Y, 600)
	for _, org := range visibleOrganisms {
		screenX := org.X - g.CameraX
		screenY := org.Y - g.CameraY
		if screenX > -100 && screenX < float64(g.ScreenWidth)+100 &&
			screenY > -100 && screenY < float64(g.ScreenHeight)+100 {
			// Draw organism as a simple circle for now
			ebitenutil.DrawCircle(screen, screenX, screenY, 20, Green)
		}
	}

	// Draw particles
	for _, particle := range g.Particles {
		px := particle.X - g.CameraX
		py := particle.Y - g.CameraY
		if px > -10 && px < float64(g.ScreenWidth)+10 &&
			py > -10 && py < float64(g.ScreenHeight)+10 {
			g.drawParticle(screen, particle)
		}
	}

	// Draw player
	g.drawPlayer(screen)

	// Draw UI
	g.drawUI(screen)
}

// drawSkyGradient draws a gradient background
func (g *Game) drawSkyGradient(screen *ebiten.Image) {
	for y := 0; y < g.ScreenHeight; y++ {
		ratio := float64(y) / float64(g.ScreenHeight)
		r := uint8(135 - ratio*85)
		gb := uint8(206 - ratio*106)
		b := uint8(235 - ratio*35)
		ebitenutil.DrawLine(screen, 0, float64(y), float64(g.ScreenWidth), float64(y), color.RGBA{r, gb, b, 255})
	}
}

// getVisibleHexagons returns all visible hexagons based on camera position
func (g *Game) getVisibleHexagons() []*world.Hexagon {
	playerChunkX, playerChunkY := g.World.GetChunkCoords(g.Player.X, g.Player.Y)
	visibleHexagons := []*world.Hexagon{}

	for chunkX := playerChunkX - g.renderDistance; chunkX <= playerChunkX+g.renderDistance; chunkX++ {
		for chunkY := playerChunkY - g.renderDistance; chunkY <= playerChunkY+g.renderDistance; chunkY++ {
			chunk := g.World.GetChunk(chunkX, chunkY)
			for _, hexagon := range chunk.Hexagons {
				screenX := hexagon.X - g.CameraX
				screenY := hexagon.Y - g.CameraY

				// Frustum culling
				if screenX > -world.HexSize && screenX < float64(g.ScreenWidth)+world.HexSize &&
					screenY > -world.HexSize && screenY < float64(g.ScreenHeight)+world.HexSize {
					visibleHexagons = append(visibleHexagons, hexagon)
				}
			}
		}
	}

	return visibleHexagons
}

// sortHexagonsByY sorts hexagons by Y position for depth rendering
func (g *Game) sortHexagonsByY(hexagons []*world.Hexagon) {
	// Simple bubble sort (can be optimized)
	n := len(hexagons)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if hexagons[j].Y > hexagons[j+1].Y {
				hexagons[j], hexagons[j+1] = hexagons[j+1], hexagons[j]
			}
		}
	}
}

// drawHexagon draws a single hexagon
func (g *Game) drawHexagon(screen *ebiten.Image, hexagon *world.Hexagon) {
	// Calculate screen positions
	corners := make([]ebiten.Vertex, len(hexagon.Corners))
	for i, corner := range hexagon.Corners {
		corners[i].DstX = float32(corner[0] - g.CameraX)
		corners[i].DstY = float32(corner[1] - g.CameraY)
		corners[i].ColorR = float32(hexagon.ActiveColor.(color.RGBA).R) / 255.0
		corners[i].ColorG = float32(hexagon.ActiveColor.(color.RGBA).G) / 255.0
		corners[i].ColorB = float32(hexagon.ActiveColor.(color.RGBA).B) / 255.0
		corners[i].ColorA = float32(hexagon.ActiveColor.(color.RGBA).A) / 255.0
	}

	// Draw filled hexagon (simplified to avoid complex triangle drawing)
	// In v2, we'll use simple line-based rendering for hexagons
	// For production, you'd want to implement proper triangle rendering

	// Draw outline (simplified as lines)
	outlineColor := DarkGray
	for i := 0; i < len(hexagon.Corners); i++ {
		next := (i + 1) % len(hexagon.Corners)
		x1 := float32(hexagon.Corners[i][0] - g.CameraX)
		y1 := float32(hexagon.Corners[i][1] - g.CameraY)
		x2 := float32(hexagon.Corners[next][0] - g.CameraX)
		y2 := float32(hexagon.Corners[next][1] - g.CameraY)
		ebitenutil.DrawLine(screen, float64(x1), float64(y1), float64(x2), float64(y2), outlineColor)
	}

	// Draw damage indicator
	if hexagon.Health < hexagon.MaxHealth {
		healthRatio := hexagon.Health / hexagon.MaxHealth
		crackColor := color.RGBA{
			R: 255,
			G: uint8(255 * healthRatio),
			B: uint8(255 * healthRatio),
			A: 255,
		}
		for i := 0; i < len(hexagon.Corners); i++ {
			next := (i + 1) % len(hexagon.Corners)
			x1 := float32(hexagon.Corners[i][0] - g.CameraX)
			y1 := float32(hexagon.Corners[i][1] - g.CameraY)
			x2 := float32(hexagon.Corners[next][0] - g.CameraX)
			y2 := float32(hexagon.Corners[next][1] - g.CameraY)
			ebitenutil.DrawLine(screen, float64(x1), float64(y1), float64(x2), float64(y2), crackColor)
		}
	}
}

// drawPlayer draws the player character
func (g *Game) drawPlayer(screen *ebiten.Image) {
	screenX := int(g.Player.X - g.CameraX)
	screenY := int(g.Player.Y - g.CameraY)

	// Draw player body
	playerRadius := g.Player.Width / 2
	ebitenutil.DrawRect(screen, float64(screenX-int(playerRadius)), float64(screenY-int(playerRadius)),
		playerRadius*2, playerRadius*2, Blue)

	// Draw eyes
	eyeColor := White
	eyeOffset := 5.0
	eyeSize := 4.0
	ebitenutil.DrawRect(screen, float64(screenX-int(eyeOffset)-int(eyeSize)), float64(screenY-3-int(eyeSize)),
		eyeSize*2, eyeSize*2, eyeColor)
	ebitenutil.DrawRect(screen, float64(screenX+int(eyeOffset)-int(eyeSize)), float64(screenY-3-int(eyeSize)),
		eyeSize*2, eyeSize*2, eyeColor)

	// Draw pupils
	pupilColor := Black
	pupilSize := 2.0
	ebitenutil.DrawRect(screen, float64(screenX-int(eyeOffset)-int(pupilSize)), float64(screenY-3-int(pupilSize)),
		pupilSize*2, pupilSize*2, pupilColor)
	ebitenutil.DrawRect(screen, float64(screenX+int(eyeOffset)-int(pupilSize)), float64(screenY-3-int(pupilSize)),
		pupilSize*2, pupilSize*2, pupilColor)

	// Draw mining range indicator
	if true {
		// Draw circle outline (simplified as polygon)
		rangeRadius := player.MiningRange
		steps := 32
		for i := 0; i < steps; i++ {
			angle1 := 2 * math.Pi * float64(i) / float64(steps)
			angle2 := 2 * math.Pi * float64(i+1) / float64(steps)
			x1 := screenX + int(math.Cos(angle1)*rangeRadius)
			y1 := screenY + int(math.Sin(angle1)*rangeRadius)
			x2 := screenX + int(math.Cos(angle2)*rangeRadius)
			y2 := screenY + int(math.Sin(angle2)*rangeRadius)
			ebitenutil.DrawLine(screen, float64(x1), float64(y1), float64(x2), float64(y2), color.RGBA{255, 255, 255, 50})
		}
	}
}

// drawParticle draws a particle
func (g *Game) drawParticle(screen *ebiten.Image, particle *Particle) {
	alpha := 1.0 - float64(particle.Age)/float64(particle.Lifetime)
	if alpha <= 0 {
		return
	}

	// Apply alpha to color
	c := particle.Color.(color.RGBA)
	colorWithAlpha := color.RGBA{
		R: uint8(float64(c.R) * alpha),
		G: uint8(float64(c.G) * alpha),
		B: uint8(float64(c.B) * alpha),
		A: uint8(float64(c.A) * alpha),
	}

	px := float64(int(particle.X - g.CameraX))
	py := float64(int(particle.Y - g.CameraY))
	size := float64(particle.Size)

	ebitenutil.DrawRect(screen, px-size/2, py-size/2, size, size, colorWithAlpha)
}

// updateParticles updates all particles
func (g *Game) updateParticles() {
	aliveParticles := []*Particle{}
	for _, particle := range g.Particles {
		particle.VelY += 0.2 // Gravity
		particle.X += particle.VelX
		particle.Y += particle.VelY
		particle.Age++

		if particle.Age < particle.Lifetime {
			aliveParticles = append(aliveParticles, particle)
		}
	}
	g.Particles = aliveParticles
}

// handleInput handles keyboard and mouse input
func (g *Game) handleInput() {
	// Movement is handled in Update loop with key states

	// Jump
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.Player.Jump()
	}

	// Mining
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.Mining = true
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.Mining = false
	}

	// Place block (simplified)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.placeBlock()
	}

	// Render distance adjustment
	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		g.renderDistance = min(10, g.renderDistance+1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		g.renderDistance = max(2, g.renderDistance-1)
	}
}

// mineBlock handles block mining
func (g *Game) mineBlock() {
	mouseX, mouseY := ebiten.CursorPosition()
	worldMX := float64(mouseX) + g.CameraX
	worldMY := float64(mouseY) + g.CameraY

	// Check for organisms first
	targetOrg := g.World.GetOrganismAt(worldMX, worldMY, 30)
	if targetOrg != nil {
		dx := g.Player.X - targetOrg.X
		dy := g.Player.Y - targetOrg.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		if distance <= player.MiningRange {
			// Damage organism
			destroyed := targetOrg.TakeDamage(1.0)
			
			// Create hit particles
			for i := 0; i < 5; i++ {
				particle := g.createParticle(targetOrg.X, targetOrg.Y, color.RGBA{200, 100, 50, 255})
				g.Particles = append(g.Particles, particle)
			}

			// Handle drops
			if destroyed {
				g.World.RemoveOrganism(targetOrg.X, targetOrg.Y)
				// Create drop particles
				for j := 0; j < 10; j++ {
					particle := g.createParticle(targetOrg.X, targetOrg.Y-20, color.RGBA{100, 100, 100, 255})
					g.Particles = append(g.Particles, particle)
				}
			}
		}
		return
	}

	// Mine block
	// Use raycasting or simple distance check (simplified here)
	// For full implementation, add raycasting from Python version
}

// placeBlock handles block placement
func (g *Game) placeBlock() {
	mouseX, mouseY := ebiten.CursorPosition()
	worldMX := float64(mouseX) + g.CameraX
	worldMY := float64(mouseY) + g.CameraY

	centerX, centerY, _, _ := world.PixelToHexCenter(worldMX, worldMY)

	// Distance checks
	dx := g.Player.X - centerX
	dy := g.Player.Y - centerY
	distSq := dx*dx + dy*dy

	if distSq > player.MiningRange*player.MiningRange {
		return
	}
	playerRadius := g.Player.Width / 2
	if distSq < (world.HexSize+playerRadius)*(world.HexSize+playerRadius) {
		return
	}

	// Check if position is already occupied
	if g.World.GetHexagonAt(centerX, centerY) != nil {
		return
	}

	// Place the block (simplified - always dirt)
	g.World.AddHexagonAt(centerX, centerY, 0) // Dirt block type
}

// createParticle creates a new particle
func (g *Game) createParticle(x, y float64, c color.Color) *Particle {
	return &Particle{
		X:        x,
		Y:        y,
		VelX:     rand.Float64()*4 - 2,
		VelY:     rand.Float64()*2 - 3,
		Color:    c,
		Lifetime: 30,
		Age:      0,
		Size:     rand.Intn(3) + 2,
	}
}

// spawnPlayer spawns the player in a safe location
func (g *Game) spawnPlayer() {
	spawnX := float64(g.ScreenWidth) / 2
	spawnY := 200.0

	// Find the surface
	foundGround := false
	for y := 0; y < g.ScreenHeight*2; y += 20 {
		if g.World.GetHexagonAt(spawnX, float64(y)) != nil {
			spawnY = float64(y - 50)
			foundGround = true
			break
		}
	}

	if !foundGround {
		spawnY = 300
	}

	g.Player = player.NewPlayer(spawnX, spawnY)
}

// drawUI draws the user interface
func (g *Game) drawUI(screen *ebiten.Image) {
	// Draw hotbar (simplified)
	hotbarWidth := 400.0
	hotbarHeight := 50.0
	hotbarX := float64(g.ScreenWidth)/2 - hotbarWidth/2
	hotbarY := float64(g.ScreenHeight) - hotbarHeight - 10

	slotWidth := hotbarWidth / 9

	for i := 0; i < 9; i++ {
		slotX := hotbarX + float64(i)*slotWidth

		// Highlight selected slot
		if i == g.Player.SelectedSlot {
			ebitenutil.DrawRect(screen, slotX, hotbarY, slotWidth-2, hotbarHeight, White)
		} else {
			ebitenutil.DrawRect(screen, slotX, hotbarY, slotWidth-2, hotbarHeight, Gray)
		}
	}

	// Draw render distance info
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Render Distance: %d", g.renderDistance))
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
