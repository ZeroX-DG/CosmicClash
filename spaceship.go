package main

import "encoding/json"

// Spaceship represents a spaceship and its properties.
type Spaceship struct {
	Name            string     `json:"name"`
	Health          int        `json:"health"`
	Position        [2]float64 `json:"position"` // x, y coordinates
	Angle           float64    `json:"angle"`
	isDestroyed     bool
	isMovingForward bool
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

// RegisterShipCommand can be used to register a ship
type ReigsterShipCommand struct {
	Name string `json:"name"`
}

func newShip(name string, initialPosition [2]float64) *Spaceship {
	return &Spaceship{
		Name:            name,
		Health:          100,
		Position:        initialPosition,
		Angle:           0,
		isDestroyed:     false,
		isMovingForward: false,
	}
}

func (s *Spaceship) update() {
	// TODO: handle ship update here
}

func (s *Spaceship) toJSON() []byte {
	b, _ := json.Marshal(s)
	return b
}

func (c *StopCommand) Execute() {
	c.Spaceship.isMovingForward = false
}

func (c *ForwardCommand) Execute() {
	c.Spaceship.isMovingForward = true
}
