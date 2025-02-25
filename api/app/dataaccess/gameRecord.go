package dataaccess

import othello "api/generated"

type GameRecord struct {
	Game othello.Game
}

func (g GameRecord) GetID() int {
	return int(g.Game.ID)
}
