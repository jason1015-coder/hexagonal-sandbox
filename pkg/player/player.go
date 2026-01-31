package player

import (
	"github.com/ninjatech/tesselbox-go/pkg/hexagon"
)

const (
	Gravity       = 0.5
	PlayerSpeed   = 5.0
	JumpForce     = -12.0
	Friction      = 0.8
	MiningRange   = 150.0
	PlayerWidth   = 40.0
	PlayerHeight  = 80.0
)

// Player represents a player in the game
type Player struct {
	X, Y          float64
	VX, VY        float64
	Width, Height float64
	
	// Movement state
	MovingLeft   bool
	MovingRight  bool
	Jumping      bool
	OnGround     bool
	
	// Mining state
	Mining       bool
	MiningTarget *hexagon.Hexagon
	MiningProgress float64
	
	// Inventory (reference to inventory)
	SelectedSlot int
	
	// Health and stats
	Health       float64
	MaxHealth    float64
}

// NewPlayer creates a new player at the specified position
func NewPlayer(x, y float64) *Player {
	return &Player{
		X:          x,
		Y:          y,
		VX:         0,
		VY:         0,
		Width:      PlayerWidth,
		Height:     PlayerHeight,
		Health:     100.0,
		MaxHealth:  100.0,
		SelectedSlot: 0,
	}
}

// Update updates the player's physics
func (p *Player) Update(deltaTime float64) {
	// Apply horizontal movement
	if p.MovingLeft {
		p.VX = -PlayerSpeed
	} else if p.MovingRight {
		p.VX = PlayerSpeed
	} else {
		p.VX *= Friction // Apply friction
	}
	
	// Apply gravity
	p.VY += Gravity
	
	// Jump
	if p.Jumping && p.OnGround {
		p.VY = JumpForce
		p.OnGround = false
	}
	
	// Update position
	p.X += p.VX * deltaTime
	p.Y += p.VY * deltaTime
	
	// Reset jumping flag
	p.Jumping = false
	
	// Update mining progress
	if p.Mining && p.MiningTarget != nil {
		p.MiningProgress += deltaTime * 10.0 // Mining speed
		if p.MiningProgress >= 100.0 {
			p.MiningProgress = 0
			p.Mining = false
		}
	} else {
		p.MiningProgress = 0
		p.Mining = false
	}
}

// GetCenter returns the center position of the player
func (p *Player) GetCenter() (float64, float64) {
	return p.X + p.Width/2.0, p.Y + p.Height/2.0
}

// GetPosition returns the top-left position of the player
func (p *Player) GetPosition() (float64, float64) {
	return p.X, p.Y
}

// SetPosition sets the player's position
func (p *Player) SetPosition(x, y float64) {
	p.X = x
	p.Y = y
}

// Move moves the player by the specified offset
func (p *Player) Move(dx, dy float64) {
	p.X += dx
	p.Y += dy
}

// GetVelocity returns the player's velocity
func (p *Player) GetVelocity() (float64, float64) {
	return p.VX, p.VY
}

// SetVelocity sets the player's velocity
func (p *Player) SetVelocity(vx, vy float64) {
	p.VX = vx
	p.VY = vy
}

// Jump makes the player jump if on ground
func (p *Player) Jump() {
	if p.OnGround {
		p.Jumping = true
	}
}

// IsOnGround returns true if the player is on the ground
func (p *Player) IsOnGround() bool {
	return p.OnGround
}

// SetOnGround sets the player's on-ground state
func (p *Player) SetOnGround(onGround bool) {
	p.OnGround = onGround
}

// GetBounds returns the player's bounding box
func (p *Player) GetBounds() (float64, float64, float64, float64) {
	return p.X, p.Y, p.X + p.Width, p.Y + p.Height
}

// GetHealth returns the player's current health
func (p *Player) GetHealth() float64 {
	return p.Health
}

// SetHealth sets the player's health
func (p *Player) SetHealth(health float64) {
	p.Health = health
	if p.Health > p.MaxHealth {
		p.Health = p.MaxHealth
	}
	if p.Health < 0 {
		p.Health = 0
	}
}

// GetMaxHealth returns the player's maximum health
func (p *Player) GetMaxHealth() float64 {
	return p.MaxHealth
}

// TakeDamage reduces the player's health
func (p *Player) TakeDamage(amount float64) {
	p.Health -= amount
	if p.Health < 0 {
		p.Health = 0
	}
}

// Heal increases the player's health
func (p *Player) Heal(amount float64) {
	p.Health += amount
	if p.Health > p.MaxHealth {
		p.Health = p.MaxHealth
	}
}

// IsAlive returns true if the player is alive
func (p *Player) IsAlive() bool {
	return p.Health > 0
}

// StartMining starts mining at the target hexagon
func (p *Player) StartMining(target *hexagon.Hexagon) {
	p.Mining = true
	p.MiningTarget = target
	p.MiningProgress = 0
}

// StopMining stops mining
func (p *Player) StopMining() {
	p.Mining = false
	p.MiningTarget = nil
	p.MiningProgress = 0
}

// GetMiningProgress returns the current mining progress (0-100)
func (p *Player) GetMiningProgress() float64 {
	return p.MiningProgress
}

// IsMining returns true if the player is currently mining
func (p *Player) IsMining() bool {
	return p.Mining
}

// GetMiningTarget returns the hexagon being mined
func (p *Player) GetMiningTarget() *hexagon.Hexagon {
	return p.MiningTarget
}

// GetMiningRange returns the player's mining range
func (p *Player) GetMiningRange() float64 {
	return MiningRange
}

// DistanceTo returns the distance from the player to a point
func (p *Player) DistanceTo(x, y float64) float64 {
	centerX, centerY := p.GetCenter()
	dx := centerX - x
	dy := centerY - y
	return dx*dx + dy*dy // Return squared distance for efficiency
}

// CanReach returns true if the player can reach a point
func (p *Player) CanReach(x, y float64) bool {
	return p.DistanceTo(x, y) <= MiningRange*MiningRange
}

// SetSelectedSlot sets the currently selected inventory slot
func (p *Player) SetSelectedSlot(slot int) {
	p.SelectedSlot = slot
}

// GetSelectedSlot returns the currently selected inventory slot
func (p *Player) GetSelectedSlot() int {
	return p.SelectedSlot
}