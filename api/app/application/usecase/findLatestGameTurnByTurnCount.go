package usecase

import (
	"api/domain/model/game"
	"api/domain/model/gameresult"
	"api/domain/model/turn"
	"context"
	"database/sql"
	"log"
)

type FindLatestGameTurnByTurnCount struct{
	turnPort turn.TurnRepository
	gamePort game.GameRepository
	gameResPort gameresult.GameResultRepository
}

func NewFindLatestGameTurnByTurnCount(
	turnPort turn.TurnRepository,
	gamePort game.GameRepository,
	gameResPort gameresult.GameResultRepository,
) *FindLatestGameTurnByTurnCount {
	return &FindLatestGameTurnByTurnCount{
		turnPort: turnPort,
		gamePort: gamePort,
		gameResPort: gameResPort,
	}
}

type FindLatestGameTurnByTurnCountOutput struct {
	TurnCount  int             `json:"turnCount"`
	Board      [][]turn.Disc `json:"board"`
	NextDisc   int             `json:"nextDisc"`
	WinnerDisc int             `json:"winnerDisc"`
}

func NewFindLatestGameTurnByTurnCountOutput(turnCount int, board [][]turn.Disc, nextDisc int, winnerDisc int) FindLatestGameTurnByTurnCountOutput {
	return FindLatestGameTurnByTurnCountOutput{
		TurnCount:  turnCount,
		Board:      board,
		NextDisc:   nextDisc,
		WinnerDisc: winnerDisc,
	}
}

func (t *FindLatestGameTurnByTurnCount) Run(ctx context.Context, db *sql.DB, turnCount int) (res FindLatestGameTurnByTurnCountOutput, err error) {
	game, err := t.gamePort.FindLatest(ctx, db)

	turn, err := t.turnPort.FindForGameIDAndTurnCount(ctx, db, int(game.GetID()), turnCount)
	if err != nil {
		log.Println(err)
		return FindLatestGameTurnByTurnCountOutput{}, err
	}

	var gameResult gameresult.GameResult
	if turn.GameEnded() {
		gameResult, err = t.gameResPort.FindForGameID(ctx, db, int(game.GetID()))
		if err != nil {
			log.Println(err)
			return FindLatestGameTurnByTurnCountOutput{}, err
		}
	}

	res = NewFindLatestGameTurnByTurnCountOutput(turn.GetTurnCount(), turn.Board.Discs, int(turn.GetNextDisc()), int(gameResult.GetWinnerDisc()))

	return res, nil
}
