"use client";

import { useEffect, useState } from "react";
import { NextDiscBanner } from "./NextDiscBanner";
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
        className="w-12 h-12 bg-green-600 border border-black flex items-center justify-center"
      >
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
    <div className="flex flex-col gap-6 w-[392px]">
      <div className="inline-block bg-green-700 p-1">
        {[...Array(8)].map((_, y) => renderRow(y))}
      </div>
      <NextDiscBanner message={bannerMessage} />
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
