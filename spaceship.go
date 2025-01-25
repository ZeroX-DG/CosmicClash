package main

// Spaceship represents a spaceship and its properties.
type Spaceship struct {
	Name            string
	Health          int
	Position        [2]float64 // x, y coordinates
	IsDestroyed     bool
	IsMovingForward bool
}

// Command is the interface for all spaceship actions.
type Command interface {
	Execute()
}

// ForwardCommand represents a forward movement action.
type ForwardCommand struct {
	Spaceship *Spaceship
}

// StopCommand represents a stop movement action.
type StopCommand struct {
	Spaceship *Spaceship
}

// RotateCommand represents a rotate movement action.
type RotateCommand struct {
	Spaceship *Spaceship
	angle     float64 // angle in radians, the rotation will be relative to the ship current angle
}

// ShootCommand represents a shooting action.
type ShootCommand struct {
	Spaceship *Spaceship
}

func newShip() *Spaceship {
	return &Spaceship{}
}

func (c *StopCommand) Execute() {
	c.Spaceship.IsMovingForward = false
}

func (c *ForwardCommand) Execute() {
	c.Spaceship.IsMovingForward = true
}
