package gameresult

import "time"

type GameResult struct {
	GameID     int
	WinnerDisc WinnerDisc
	EndAt      time.Time
}

func NewGameResult(gameID int, winnerDisc WinnerDisc) *GameResult {
	return &GameResult{
		GameID:     gameID,
		WinnerDisc: winnerDisc,
	}
}

func (g GameResult) GetGameID() int {
	return g.GameID
}

func (g GameResult) GetWinnerDisc() WinnerDisc {
	return g.WinnerDisc
}

func (g GameResult) GetEndAt() time.Time {
	return g.EndAt
}
