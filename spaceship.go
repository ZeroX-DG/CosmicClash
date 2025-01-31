package main

import (
	"encoding/json"
	"log"
	"math"
)

// Spaceship represents a spaceship and its properties.
type Spaceship struct {
	Name               string     `json:"name"`
	Health             int        `json:"health"`
	Position           [2]float64 `json:"position"` // x, y coordinates
	Angle              float64    `json:"angle"`
	isDestroyed        bool
	isMovingForward    bool
	accelerationFactor float64
	maxAcceleration    float64
	velocity           [2]float64
	thurstVector       float64
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
	angle     float64 // angle in radians
}

// ShootCommand represents a shooting action.
type ShootCommand struct {
	Spaceship *Spaceship
}

// RegisterShipCommand can be used to register a ship
type ReigsterShipCommand struct {
	Client *Client
	Game   *Game
	Name   string
}

func newShip(name string, initialPosition [2]float64) *Spaceship {
	return &Spaceship{
		Name:               name,
		Health:             100,
		Position:           initialPosition,
		Angle:              0,
		isDestroyed:        false,
		isMovingForward:    false,
		accelerationFactor: 0,
		maxAcceleration:    0.07,
		velocity:           [2]float64{0, 0},
		thurstVector:       0,
	}
}

func (s *Spaceship) update() {
	if s.isMovingForward {
		s.thrust()
	}

	s.Position[0] += s.velocity[0]
	s.Position[1] += s.velocity[1]
}

func (s *Spaceship) thrust() {
	s.accelerationFactor = min(s.accelerationFactor+0.01, 1)
	acceleration := s.maxAcceleration * easeOutQuad(s.accelerationFactor)
	s.velocity[0] += acceleration * math.Sin(s.thurstVector)
	s.velocity[1] += -acceleration * math.Cos(s.thurstVector)
}

func easeOutQuad(t float64) float64 {
	return 1 - math.Pow((1-t), 3)
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

func (c *RotateCommand) Execute() {
	c.Spaceship.Angle = c.angle
}

func (c *ReigsterShipCommand) Execute() {
	ship := newShip(c.Name, c.Game.chooseRandomPosition())
	c.Game.ships[c.Client] = ship

	c.Game.hub.broadcast <- c.Game.toJSON()
	c.Client.send <- ship.toJSON()
	log.Println("Player Joined: " + c.Name)
}
