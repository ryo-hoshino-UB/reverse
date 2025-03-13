"use client";

import { useEffect, useState } from "react";
import { Stone } from "./Stone";
import type { Disc } from "./disc";
import { fetchApi } from "./fetch";
import type { TurnRequest, TurnResponse } from "./interface";
import { WINNER_DISC } from "./winnerDisc";

const registerGame = async () => {
  const res = await fetchApi("/api/games", {
    method: "POST",
  });
  const game = await res.json();
  return game;
};

const getTurn = async (turnCount: number) => {
  const res = await fetchApi(`/api/games/latest/turns/${turnCount}`);
  const turn: TurnResponse = await res.json();
  return turn;
};

const registerTurn = async (turnReq: TurnRequest) => {
  const res = await fetchApi("/api/games/latest/turns", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(turnReq),
  });
  if (!res.ok) {
    throw new Error("Failed to register turn");
  }
  const newTurn = await res.json();
  return newTurn;
};

export const Board: React.FC = () => {
  const [nextDisc, setNextDisc] = useState<Disc>(0);
  const [turnCount, setTurnCount] = useState(0);
  const [board, setBoard] = useState<Disc[][]>(
    Array(8)
      .fill(null)
      .map(() => Array(8).fill(null))
  );
  const [bannerMessage, setBannerMessage] = useState(
    changeBannerMessage({ nextDisc: 1 })
  );

  useEffect(() => {
    const startGame = async () => {
      await registerGame();
    };

    const fetchTurn = async () => {
      const turn = await getTurn(0);
      setBoard(turn.board);
      setNextDisc(turn.nextDisc);
    };

    startGame();
    fetchTurn();
  }, []);

  const renderSquare = (y: number, x: number) => {
    const disc = board[y][x];

    const handleSquareClick = async () => {
      const nextTurnCount = turnCount + 1;

      await registerTurn({
        turnCount: nextTurnCount,
        move: {
          disc: nextDisc,
          x,
          y,
        },
      });
      // registerTurnの後に直接最新状態を取得
      const turn = await getTurn(nextTurnCount);
      console.log(turn);
      if (turn.nextDisc !== 0) {
        setBannerMessage(changeBannerMessage({ nextDisc: turn.nextDisc }));
        if (nextDisc === turn.nextDisc) {
          setBannerMessage(
            `もう一度${nextDisc === 1 ? "白" : "黒"}のターンです`
          );
        }
      } else {
        if (turn.winnerDisc === WINNER_DISC.Draw) {
          setBannerMessage("引き分けです");
        } else {
          setBannerMessage(`${turn.winnerDisc === 1 ? "黒" : "白"}の勝ちです`);
        }
      }
      setBoard(turn.board);
      setNextDisc(turn.nextDisc);
      setTurnCount(nextTurnCount);
    };

    return (
      // biome-ignore lint/a11y/useKeyWithClickEvents:
      <div
        onClick={handleSquareClick}
        key={`${y}-${x}`}
        className="w-14 h-14 bg-gradient-to-br from-emerald-700 to-emerald-800 border border-emerald-900/70 flex items-center justify-center relative overflow-hidden transition-all duration-300 hover:shadow-inner hover:from-emerald-600 hover:to-emerald-700"
      >
        {/* セルの光沢効果 */}
        <div className="absolute top-0 left-0 right-0 h-[1px] bg-emerald-400/30" />
        <div className="absolute top-0 left-0 bottom-0 w-[1px] bg-emerald-400/30" />
        <div className="absolute bottom-0 left-0 right-0 h-[1px] bg-emerald-900/50" />
        <div className="absolute top-0 left-0 bottom-0 w-[1px] bg-emerald-400/30" />
        <Stone disc={disc} />
      </div>
    );
  };

  const renderRow = (y: number) => {
    return (
      <div key={y} className="flex">
        {[...Array(8)].map((_, x) => renderSquare(y, x))}
      </div>
    );
  };

  return (
    <div className="flex flex-col items-center gap-12">
      <div className="inline-block bg-gradient-to-br from-emerald-800 to-emerald-900 p-1.5 rounded-lg shadow-[0_0_15px_rgba(16,185,129,0.2)] border border-emerald-700/30">
        {[...Array(8)].map((_, y) => renderRow(y))}
      </div>

      <div className="bg-gray-800/90 px-6 py-3 rounded-full shadow-lg border border-emerald-500/30 backdrop-blur-sm">
        <div className="relative overflow-hidden">
          <span className="text-white font-bold text-lg tracking-wide relative z-10 flex items-center justify-center">
            {bannerMessage}
          </span>
        </div>
      </div>

      {/* ゲーム情報 */}
      <div className="flex justify-between w-full bg-gray-800/80 rounded-lg p-4 border border-gray-700/50 shadow-md text-base">
        <div className="flex items-center gap-3">
          <div className="w-5 h-5 rounded-full bg-black shadow-[0_0_5px_rgba(0,0,0,0.5)]" />
          <span className="text-white font-medium">
            黒: {board.flat().filter((d) => d === 1).length}
          </span>
        </div>
        <div className="flex items-center gap-3">
          <span className="text-white font-medium">
            白: {board.flat().filter((d) => d === 2).length}
          </span>
          <div className="w-5 h-5 rounded-full bg-white shadow-[0_0_5px_rgba(255,255,255,0.5)]" />
        </div>
      </div>
    </div>
  );
};

type Props = {
  nextDisc: Disc;
  skip?: boolean;
};

const changeBannerMessage = ({ nextDisc }: Props): string => {
  let nextDiscColor = "";
  if (nextDisc === 1) {
    nextDiscColor = "黒";
  }
  if (nextDisc === 2) {
    nextDiscColor = "白";
  }

  return `${nextDiscColor}のターンです`;
};
