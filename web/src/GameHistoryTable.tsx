import { useState } from "react";
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
  const [hoveredRow, setHoveredRow] = useState<number | null>(null);

  // 日付のフォーマットを整える関数
  const formatDate = (dateString: string) => {
    try {
      const date = new Date(dateString);
      // Check if date is valid
      if (!date.getTime()) {
        return "-";
      }

      return new Intl.DateTimeFormat("ja-JP", {
        year: "numeric",
        month: "2-digit",
        day: "2-digit",
        hour: "2-digit",
        minute: "2-digit",
      }).format(date);
    } catch {
      return "Invalid date";
    }
  };

  // 勝者の表示をレンダリングする関数
  const renderWinner = (
    winnerDisc: WinnerDisc,
    blackMoveCount: number,
    whiteMoveCount: number
  ) => {
    if (blackMoveCount === 0 && whiteMoveCount === 0) {
      return (
        <div className="flex items-center gap-2 text-gray-500">
          <span className="font-medium">-</span>
        </div>
      );
    }
    switch (winnerDisc) {
      case 1:
        return (
          <div className="flex items-center gap-2">
            <div className="w-5 h-5 rounded-full bg-gray-900 shadow-[0_0_4px_rgba(0,0,0,0.4)]" />
            <span className="font-medium">黒</span>
          </div>
        );
      case 2:
        return (
          <div className="flex items-center gap-2">
            <div className="w-5 h-5 rounded-full bg-white shadow-[0_0_4px_rgba(255,255,255,0.4)] border border-gray-300" />
            <span className="font-medium">白</span>
          </div>
        );
      default:
        return (
          <div className="flex items-center gap-2 text-gray-500">
            <span className="font-medium">引き分け</span>
          </div>
        );
    }
  };

  const thStyle =
    "bg-gradient-to-r from-gray-800 to-gray-900 text-left py-3 px-4 text-sm uppercase tracking-wide font-semibold text-emerald-400 border-b border-gray-700";

  return (
    <div className="w-full overflow-hidden rounded-xl border border-gray-700/50 shadow-[0_0_15px_rgba(16,185,129,0.15)]">
      <div className="overflow-x-auto">
        <table className="min-w-full">
          <thead>
            <tr>
              <th className={thStyle}>黒の手数</th>
              <th className={thStyle}>白の手数</th>
              <th className={thStyle}>勝者</th>
              <th className={thStyle}>開始日時</th>
              <th className={thStyle}>終了日時</th>
            </tr>
          </thead>
          <tbody>
            {gameHistories.map((game, index) => (
              <tr
                key={game.startedAt}
                className={`
                  border-b border-gray-700/20 backdrop-blur-sm
                  ${index % 2 === 0 ? "bg-gray-800/30" : "bg-gray-800/50"}
                  ${hoveredRow === index ? "bg-gray-700/50" : ""}
                  transition-colors duration-150
                `}
                onMouseEnter={() => setHoveredRow(index)}
                onMouseLeave={() => setHoveredRow(null)}
              >
                <td className="py-3 px-4 text-white/90">
                  <div className="flex items-center gap-2">
                    <span className="font-medium">{game.blackMoveCount}</span>
                    <span className="text-xs text-gray-400">手</span>
                  </div>
                </td>
                <td className="py-3 px-4 text-white/90">
                  <div className="flex items-center gap-2">
                    <span className="font-medium">{game.whiteMoveCount}</span>
                    <span className="text-xs text-gray-400">手</span>
                  </div>
                </td>
                <td className="py-3 px-4 text-white/90">
                  {renderWinner(
                    game.winnerDisc,
                    gameHistories[0].blackMoveCount,
                    gameHistories[0].whiteMoveCount
                  )}
                </td>
                <td className="py-3 px-4 text-white/70 text-sm">
                  {formatDate(game.startedAt)}
                </td>
                <td className="py-3 px-4 text-white/70 text-sm">
                  {formatDate(game.endedAt)}
                </td>
              </tr>
            ))}
            {gameHistories.length === 0 && (
              <tr>
                <td
                  colSpan={5}
                  className="py-8 text-center text-gray-400 bg-gray-800/30 italic"
                >
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
