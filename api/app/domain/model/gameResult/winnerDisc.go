package gameresult

type WinnerDisc int

const (
	Draw WinnerDisc = iota
	BlackWin
	WhiteWin
)

func ToWinnerDisc(winnerDisc int) WinnerDisc {
	switch winnerDisc {
	case int(BlackWin):
		return BlackWin
	case int(WhiteWin):
		return WhiteWin
	default:
		panic("invalid winnerDisc")
	}
}
