package turn

import (
	"api/xerrors"
)

type Board struct {
	Discs       [][]Disc
	WalledDiscs [][]Disc
}

func NewBoard(discs [][]Disc) Board {
	board := Board{
		Discs: discs,
	}
	board.WalledDiscs = board.WallDiscs()
	return board
}

func (b *Board) ListFlipPoints(move Move) []Point {
	flipPoints := make([]Point, 0)

	walledX := move.Point.X + 1
	walledY := move.Point.Y + 1

	checkFlipPoints := func(xMove, yMove int) {
		flipCandidates := make([]Point, 0)

		cursorX := walledX + xMove
		cursorY := walledY + yMove

		// 打った石と逆の色の石が続く限り走査する
		for IsOppositeDisc(move.Disc, b.WalledDiscs[cursorY][cursorX]) {
			point, _ := NewPoint(cursorX-1, cursorY-1)
			flipCandidates = append(flipCandidates, point)
			cursorX += xMove
			cursorY += yMove
			if move.Disc == b.WalledDiscs[cursorY][cursorX] {
				flipPoints = append(flipPoints, flipCandidates...)
				break
			}
		}
	}

	checkFlipPoints(0, -1)
	checkFlipPoints(1, -1)
	checkFlipPoints(1, 0)
	checkFlipPoints(1, 1)
	checkFlipPoints(0, 1)
	checkFlipPoints(-1, 1)
	checkFlipPoints(-1, 0)
	checkFlipPoints(-1, -1)

	return flipPoints
}

func (b *Board) WallDiscs() [][]Disc {
	walled := make([][]Disc, 0)

	topAndBottomWall := make([]Disc, len(b.Discs[0])+2)
	for i := range topAndBottomWall {
		topAndBottomWall[i] = Wall
	}
	walled = append(walled, topAndBottomWall)

	for i := range b.Discs {
		line := make([]Disc, 0)
		line = append(line, Wall)
		line = append(line, b.Discs[i]...)
		line = append(line, Wall)
		walled = append(walled, line)
	}

	walled = append(walled, topAndBottomWall)
	return walled
}

func (b *Board) Place(move Move) (Board, error) {
	// 空のマス目ではない場合はおけない
	if b.Discs[move.Point.Y][move.Point.X] != Empty {
		return *b, xerrors.ErrBadRequest
	}

	// ひっくり返せる場所をリストアップ
	flipPoints := b.ListFlipPoints(move)

	// ひっくり返せる場所がない場合はおけない
	if len(flipPoints) == 0 {
		return *b, xerrors.ErrBadRequest
	}

	// 盤面をコピー
	newDiscs := make([][]Disc, len(b.Discs))
	for i := range b.Discs {
		newDiscs[i] = make([]Disc, len(b.Discs[i]))
		copy(newDiscs[i], b.Discs[i])
	}

	// 石を置く
	newDiscs[move.Point.Y][move.Point.X] = move.Disc

	// 石をひっくり返す
	for _, p := range flipPoints {
		newDiscs[p.Y][p.X] = move.Disc
	}

	return NewBoard(newDiscs), nil
}

func (b *Board) ExistValidMove(disc Disc) bool {
	for y := MIN_POINT; y < MAX_POINT; y++ {
		line := b.Discs[y]
		for x := MIN_POINT; x < MAX_POINT; x++ {
			discOnBoard := line[x]
			if discOnBoard != Empty {
				continue
			}
			point, _ := NewPoint(x, y)
			move := NewMove(disc, point)
			if len(b.ListFlipPoints(move)) > 0 {
				return true
			}
			continue
		}
	}
	return false
}

func (b *Board) CountDiscs(disc Disc) int {
	count := 0
	for _, line := range b.Discs {
		for _, d := range line {
			if d == disc {
				count++
			}
		}
	}
	return count
}

func (b *Board) GetDiscs() [][]Disc {
	return b.Discs
}

func NewInitialBoard() Board {
	discs := make([][]Disc, 8)
	for i := range discs {
		discs[i] = make([]Disc, 8)
		for j := range discs[i] {
			discs[i][j] = Empty
		}
	}

	discs[3][4] = White
	discs[3][3] = Black
	discs[4][4] = Black
	discs[4][3] = White

	return NewBoard(discs)
}
