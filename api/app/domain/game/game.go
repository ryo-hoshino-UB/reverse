package game

import "time"

type Game struct {
	ID        int
	StartedAt time.Time
}

func NewGame(id int, started_at time.Time) Game {
	return Game{
		ID:        id,
		StartedAt: time.Now(),
	}
}

func (g Game) GetID() int {
	return g.ID
}

func (g Game) GetStartedAt() time.Time {
	return g.StartedAt
}
