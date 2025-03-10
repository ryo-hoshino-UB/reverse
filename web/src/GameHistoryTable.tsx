import type { WinnerDisc } from "./winnerDisc";

export type GameHistory = {
  id: number;
  blackMoveCount: number;
  whiteMoveCount: number;
  winnerDisc: WinnerDisc;
  startedAt: string;
  endedAt: string;
};

type Props = {
  gameHistories: GameHistory[];
};
export const GameHistoryTable: React.FC<Props> = ({ gameHistories }) => {
  // 勝者の表示をレンダリングする関数
  const renderWinner = (winnerDisc: WinnerDisc) => {
    switch (winnerDisc) {
      case 1:
        return (
          <div className="flex items-center justify-center">
            <div className="w-6 h-6 rounded-full bg-gray-800" />
            <span className="ml-2 font-medium">黒</span>
          </div>
        );
      case 2:
        return (
          <div className="flex items-center justify-center">
            <div className="w-6 h-6 rounded-full bg-white border-2 border-gray-800" />
            <span className="ml-2 font-medium">白</span>
          </div>
        );
      default:
        return <span className="text-gray-500 italic">引き分け</span>;
    }
  };

  return (
    <div className="w-full overflow-hidden shadow-lg rounded-lg">
      <div className="overflow-x-auto">
        <table className="min-w-full bg-white">
          <thead>
            <tr className="bg-gray-800 text-white">
              <th className="py-3 px-4 text-left">黒を打った回数</th>
              <th className="py-3 px-4 text-left">白を打った回数</th>
              <th className="py-3 px-4 text-left">勝者</th>
              <th className="py-3 px-4 text-left">対戦開始時刻</th>
              <th className="py-3 px-4 text-left">対戦終了時刻</th>
            </tr>
          </thead>
          <tbody>
            {gameHistories.map((game, index) => (
              <tr
                key={game.startedAt}
                className={`border-b hover:bg-gray-50 transition-colors ${
                  index % 2 === 0 ? "bg-gray-50" : "bg-white"
                }`}
              >
                <td className="py-3 px-4">{game.blackMoveCount}</td>
                <td className="py-3 px-4">{game.whiteMoveCount}</td>
                <td className="py-3 px-4">{renderWinner(game.winnerDisc)}</td>
                <td className="py-3 px-4">{game.startedAt}</td>
                <td className="py-3 px-4">{game.endedAt}</td>
              </tr>
            ))}
            {gameHistories.length === 0 && (
              <tr>
                <td colSpan={5} className="py-8 text-center text-gray-500">
                  対戦履歴がありません
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};
