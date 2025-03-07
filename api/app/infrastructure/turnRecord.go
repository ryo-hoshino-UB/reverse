package infrastructure

import (
	othello "api/generated"
	"time"
)

type TurnRecord struct {
	Turn othello.Turn
}

func (t TurnRecord) GetID() int {
	return int(t.Turn.ID)
}

func (t TurnRecord) GetGameID() int {
	return int(t.Turn.GameID)
}

func (t TurnRecord) GetTurnCount() int {
	return int(t.Turn.TurnCount)
}

func (t TurnRecord) GetNextDisc() int {
	return int(t.Turn.NextDisc.Int32)
}

func (t TurnRecord) GetEndAt() time.Time {
	return t.Turn.EndAt
}
