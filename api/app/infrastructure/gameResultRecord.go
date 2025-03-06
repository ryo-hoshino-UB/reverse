package infrastructure

import "time"

type GameResultRecord struct {
	ID         int
	GameID     int
	WinnerDisc int
	EndAt      time.Time
}

func (g GameResultRecord) GetID() int {
	return g.ID
}

func (g GameResultRecord) GetGameID() int {
	return g.GameID
}

func (g GameResultRecord) GetWinnerDisc() int {
	return g.WinnerDisc
}

func (g GameResultRecord) GetEndAt() time.Time {
	return g.EndAt
}
