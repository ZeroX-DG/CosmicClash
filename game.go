package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

type Game struct {
	hub *Hub

	ships map[*Client]*Spaceship

	size [2]uint // width x height

	messageQueue chan struct {
		client  *Client
		message []byte
	}
}

func newGame(hub *Hub) *Game {
	return &Game{
		hub:   hub,
		ships: make(map[*Client]*Spaceship),
		size:  [2]uint{800, 800},
		messageQueue: make(chan struct {
			client  *Client
			message []byte
		}),
	}
}

func (g *Game) run() {
	for {
		select {
		case message := <-g.messageQueue:
			g.processMessage(message.client, message.message)
		case client := <-g.hub.unregister:
			delete(g.ships, client)
		}
	}
}

// Translate the current game state to JSON to boardcast to all players
func (g *Game) toJSON() []byte {
	ships := make([]*Spaceship, 0, len(g.ships))

	for _, ship := range g.ships {
		ships = append(ships, ship)
	}

	b, err := json.Marshal(&struct {
		Ships []*Spaceship `json:"ships"`
		Size  [2]uint      `json:"size"`
	}{
		Ships: ships,
		Size:  g.size,
	})

	if err != nil {
		fmt.Println(err)
	}

	return b
}

func (g *Game) processMessage(client *Client, message []byte) {
	ship := g.ships[client]

	// If client have no ship then the first command is always register ship
	if ship == nil {
		var registerCommand struct {
			Command string
			Name    string
		}
		err := json.Unmarshal(message, &registerCommand)

		if err != nil {
			g.hub.broadcast <- makeJSONError("Invalid command format. The only command available now is registerShip")
			return
		}

		g.ships[client] = newShip(registerCommand.Name, g.chooseRandomPosition())

		// if register ship success then you should receive the game state
		g.hub.broadcast <- g.toJSON()
		return
	}

	command := parseCommand(ship, message)
	command.Execute()
	// let everyone know the game state has changed
	g.hub.broadcast <- g.toJSON()
}

func (g *Game) chooseRandomPosition() [2]float64 {
	return [2]float64{
		float64(rand.Intn(int(g.size[0]))),
		float64(rand.Intn(int(g.size[1]))),
	}
}

func makeJSONError(err string) []byte {
	b, _ := json.Marshal(&struct {
		Err string `json:"error"`
	}{Err: err})
	return b
}

func parseCommand(ship *Spaceship, message []byte) Command {
	// TODO: parse message JSON to command here
	return &StopCommand{
		Spaceship: ship,
	}
}
