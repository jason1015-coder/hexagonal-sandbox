package items

import "image/color"

// ItemType represents the type of an item
type ItemType int

const (
	NONE ItemType = iota
	DIRT_BLOCK
	GRASS_BLOCK
	STONE_BLOCK
	SAND_BLOCK
	LOG_BLOCK
	COAL
	IRON_INGOT
	GOLD_INGOT
	DIAMOND
	IRON_PICKAXE
	STONE_PICKAXE
	WOODEN_PICKAXE
)

// ItemProperties defines the properties of an item type
type ItemProperties struct {
	ID           ItemType
	Name         string
	IconColor    color.RGBA
	Description  string
	StackSize    int
	Durability   int // For tools, -1 for indestructible
	IsTool       bool
	ToolPower    float64 // Mining speed multiplier
}

// ItemDefinitions holds all item type definitions
var ItemDefinitions = map[ItemType]*ItemProperties{
	NONE: {
		ID:          NONE,
		Name:        "None",
		IconColor:   color.RGBA{0, 0, 0, 0},
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	DIRT_BLOCK: {
		ID:          DIRT_BLOCK,
		Name:        "Dirt",
		IconColor:   color.RGBA{139, 90, 43, 255},
		Description: "A block of dirt",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	GRASS_BLOCK: {
		ID:          GRASS_BLOCK,
		Name:        "Grass Block",
		IconColor:   color.RGBA{100, 200, 100, 255},
		Description: "A block of grass",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	STONE_BLOCK: {
		ID:          STONE_BLOCK,
		Name:        "Stone",
		IconColor:   color.RGBA{169, 169, 169, 255},
		Description: "A block of stone",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	SAND_BLOCK: {
		ID:          SAND_BLOCK,
		Name:        "Sand",
		IconColor:   color.RGBA{238, 214, 175, 255},
		Description: "A block of sand",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	LOG_BLOCK: {
		ID:          LOG_BLOCK,
		Name:        "Log",
		IconColor:   color.RGBA{139, 69, 19, 255},
		Description: "A wooden log",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	COAL: {
		ID:          COAL,
		Name:        "Coal",
		IconColor:   color.RGBA{45, 45, 45, 255},
		Description: "A piece of coal",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	IRON_INGOT: {
		ID:          IRON_INGOT,
		Name:        "Iron Ingot",
		IconColor:   color.RGBA{169, 166, 150, 255},
		Description: "An iron ingot",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	GOLD_INGOT: {
		ID:          GOLD_INGOT,
		Name:        "Gold Ingot",
		IconColor:   color.RGBA{255, 215, 0, 255},
		Description: "A gold ingot",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	DIAMOND: {
		ID:          DIAMOND,
		Name:        "Diamond",
		IconColor:   color.RGBA{0, 255, 255, 255},
		Description: "A diamond",
		StackSize:   64,
		Durability:  -1,
		IsTool:      false,
	},
	WOODEN_PICKAXE: {
		ID:          WOODEN_PICKAXE,
		Name:        "Wooden Pickaxe",
		IconColor:   color.RGBA{139, 69, 19, 255},
		Description: "A basic wooden pickaxe",
		StackSize:   1,
		Durability:  60,
		IsTool:      true,
		ToolPower:   2.0,
	},
	STONE_PICKAXE: {
		ID:          STONE_PICKAXE,
		Name:        "Stone Pickaxe",
		IconColor:   color.RGBA{169, 169, 169, 255},
		Description: "A sturdy stone pickaxe",
		StackSize:   1,
		Durability:  132,
		IsTool:      true,
		ToolPower:   4.0,
	},
	IRON_PICKAXE: {
		ID:          IRON_PICKAXE,
		Name:        "Iron Pickaxe",
		IconColor:   color.RGBA{169, 166, 150, 255},
		Description: "A durable iron pickaxe",
		StackSize:   1,
		Durability:  251,
		IsTool:      true,
		ToolPower:   6.0,
	},
}

// Item represents a stack of items
type Item struct {
	Type      ItemType
	Quantity  int
	Durability int // For tools
}

// Inventory represents a player's inventory
type Inventory struct {
	Slots    []Item
	Selected int
}

// NewInventory creates a new inventory with the specified number of slots
func NewInventory(slotCount int) *Inventory {
	slots := make([]Item, slotCount)
	for i := range slots {
		slots[i] = Item{Type: NONE, Quantity: 0, Durability: -1}
	}
	return &Inventory{
		Slots:    slots,
		Selected: 0,
	}
}

// DefaultHotbarItems returns the default items for the hotbar
func DefaultHotbarItems() []Item {
	return []Item{
		{Type: WOODEN_PICKAXE, Quantity: 1, Durability: 60},
		{Type: STONE_PICKAXE, Quantity: 1, Durability: 132},
		{Type: DIRT_BLOCK, Quantity: 10, Durability: -1},
		{Type: STONE_BLOCK, Quantity: 10, Durability: -1},
		{Type: LOG_BLOCK, Quantity: 5, Durability: -1},
		{Type: NONE, Quantity: 0, Durability: -1},
		{Type: NONE, Quantity: 0, Durability: -1},
		{Type: NONE, Quantity: 0, Durability: -1},
	}
}

// AddItem adds an item to the inventory
func (inv *Inventory) AddItem(itemType ItemType, quantity int) bool {
	props := ItemDefinitions[itemType]
	if props == nil {
		return false
	}
	
	remaining := quantity
	
	// First, try to stack with existing items
	if props.StackSize > 1 {
		for i := range inv.Slots {
			if inv.Slots[i].Type == itemType && inv.Slots[i].Quantity < props.StackSize {
				canAdd := props.StackSize - inv.Slots[i].Quantity
				add := min(canAdd, remaining)
				inv.Slots[i].Quantity += add
				remaining -= add
				if remaining == 0 {
					return true
				}
			}
		}
	}
	
	// Then, try to find empty slots
	for i := range inv.Slots {
		if inv.Slots[i].Type == NONE {
			inv.Slots[i].Type = itemType
			inv.Slots[i].Durability = props.Durability
			inv.Slots[i].Quantity = min(remaining, props.StackSize)
			remaining -= inv.Slots[i].Quantity
			if remaining == 0 {
				return true
			}
		}
	}
	
	return remaining == 0
}

// RemoveItem removes items from the selected slot
func (inv *Inventory) RemoveItem(quantity int) bool {
	if inv.Selected >= len(inv.Slots) {
		return false
	}
	
	slot := &inv.Slots[inv.Selected]
	if slot.Quantity < quantity {
		return false
	}
	
	slot.Quantity -= quantity
	if slot.Quantity <= 0 {
		slot.Type = NONE
		slot.Durability = -1
	}
	
	return true
}

// GetSelectedItem returns the currently selected item
func (inv *Inventory) GetSelectedItem() *Item {
	if inv.Selected >= len(inv.Slots) {
		return nil
	}
	return &inv.Slots[inv.Selected]
}

// UseItem uses the selected item (returns true if successful)
func (inv *Inventory) UseItem() bool {
	item := inv.GetSelectedItem()
	if item == nil || item.Type == NONE {
		return false
	}
	
	props := ItemDefinitions[item.Type]
	if props.IsTool && item.Durability > 0 {
		item.Durability--
		if item.Durability <= 0 {
			item.Type = NONE
			item.Quantity = 0
			item.Durability = -1
		}
		return true
	}
	
	if !props.IsTool && item.Quantity > 0 {
		item.Quantity--
		if item.Quantity <= 0 {
			item.Type = NONE
			item.Durability = -1
		}
		return true
	}
	
	return false
}

// SelectSlot selects a slot by index
func (inv *Inventory) SelectSlot(index int) bool {
	if index < 0 || index >= len(inv.Slots) {
		return false
	}
	inv.Selected = index
	return true
}

// NextSlot selects the next slot
func (inv *Inventory) NextSlot() {
	inv.Selected = (inv.Selected + 1) % len(inv.Slots)
}

// PrevSlot selects the previous slot
func (inv *Inventory) PrevSlot() {
	inv.Selected = (inv.Selected - 1 + len(inv.Slots)) % len(inv.Slots)
}

// ItemNameByID returns the name of an item type
func ItemNameByID(itemType ItemType) string {
	if props, ok := ItemDefinitions[itemType]; ok {
		return props.Name
	}
	return "Unknown"
}

// ItemColorByID returns the icon color of an item type
func ItemColorByID(itemType ItemType) color.RGBA {
	if props, ok := ItemDefinitions[itemType]; ok {
		return props.IconColor
	}
	return color.RGBA{0, 0, 0, 255}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}