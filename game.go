package main

type Game struct {
	ships map[*Client]*Spaceship

	messageQueue chan struct {
		client  *Client
		message string
	}
}

func newGame() *Game {
	return &Game{
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
		}
	}
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
