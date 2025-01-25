package main

type Game struct {
	hub *Hub

	ships map[*Client]*Spaceship

	messageQueue chan struct {
		client  *Client
		message string
	}
}

func newGame(hub *Hub) *Game {
	return &Game{
		hub:   hub,
		ships: make(map[*Client]*Spaceship),
		messageQueue: make(chan struct {
			client  *Client
			message string
		}),
	}
}

func (g *Game) run() {
	for {
		select {
		case message := <-g.messageQueue:
			g.processMessage(message.client, message.message)
			// let everyone know the game state
			g.hub.broadcast <- []byte(g.toJSON())
		}
	}
}

// Translate the current game state to JSON to boardcast to all players
func (g *Game) toJSON() string {
	return ""
}

func (g *Game) processMessage(client *Client, message string) {
	ship := g.ships[client]
	command := parseCommand(ship, message)
	command.Execute()
}

func parseCommand(ship *Spaceship, message string) Command {
	// TODO: parse message JSON to command here
	return StopCommand{
		Spaceship: ship,
	}
}
