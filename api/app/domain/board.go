package domain

type Board struct {
	Discs [][]Disc
}

func NewBoard(discs [][]Disc) Board {
	return Board{
		Discs: discs,
	}
}

func (b Board) Place(move Move) Board {
	// 盤面におけるかチェック

	// 盤面をコピー
	newDiscs := make([][]Disc, len(b.Discs))
	for i := range b.Discs {
		newDiscs[i] = make([]Disc, len(b.Discs[i]))
		copy(newDiscs[i], b.Discs[i])
	}

	// 石を置く
	newDiscs[move.Point.Y][move.Point.X] = move.Disc

	// 石をひっくり返す

	return NewBoard(newDiscs)
}

func (b Board) GetDiscs() [][]Disc {
	return b.Discs
}

func (b Board) NewInitialBoard() Board {
	discs := make([][]Disc, 8)
	for i := range discs {
		discs[i] = make([]Disc, 8)
		for j := range discs[i] {
			discs[i][j] = Empty
		}
	}

	discs[3][3] = White
	discs[4][4] = White
	discs[3][4] = Black
	discs[4][3] = Black

	return NewBoard(discs)
}
