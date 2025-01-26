package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

var (
	gameTick      = 16 * time.Millisecond
	broadcastTick = 1 * time.Second
)

type Game struct {
	hub *Hub

	ships map[*Client]*Spaceship

	size [2]uint // width x height

	messageQueue chan struct {
		client  *Client
		message []byte
	}

	unregister chan *Client
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
		unregister: make(chan *Client),
	}
}

func (g *Game) run() {
	gameTicker := time.NewTicker(gameTick)
	broadcastTicker := time.NewTicker(broadcastTick)
	defer gameTicker.Stop()
	defer broadcastTicker.Stop()

	for {
		select {
		case message := <-g.messageQueue:
			g.processMessage(message.client, message.message)
		case client := <-g.unregister:
			if client != nil && g.ships[client] != nil {
				log.Println("Unregistered Ship: " + g.ships[client].Name)
				delete(g.ships, client)
			}
		case <-gameTicker.C:
			g.update()
		case <-broadcastTicker.C:
			g.hub.broadcast <- g.toJSON()
		}
	}
}

func (g *Game) update() {
	for _, ship := range g.ships {
		ship.update()
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
		log.Println(err)
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

		ship = newShip(registerCommand.Name, g.chooseRandomPosition())
		g.ships[client] = ship

		// if register ship success then you should receive the game state & your ship info
		g.hub.broadcast <- g.toJSON()
		client.send <- ship.toJSON()
		log.Println("New Ship: " + ship.Name)
		return
	}

	command := parseCommand(ship, message)
	command.Execute()
	// let everyone know the game state has changed with the ship info
	g.hub.broadcast <- g.toJSON()
	client.send <- ship.toJSON()
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
