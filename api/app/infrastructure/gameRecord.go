package infrastructure

import (
	othello "api/generated"
	"time"
)

type GameRecord struct {
	Game othello.Game
}

func (g GameRecord) GetID() int {
	return int(g.Game.ID)
}

func (g GameRecord) GetStartedAt() time.Time {
	return g.Game.StartedAt
}
