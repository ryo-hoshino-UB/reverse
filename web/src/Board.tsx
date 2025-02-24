"use client";

import { useEffect, useState } from "react";
import { Disc } from "./disc";
import { fetchApi } from "./fetch";
import { TurnRequest, TurnResponse } from "./interface";
import { Stone } from "./Stone";

const getTurn = async (turnCount: number) => {
  const res = await fetchApi(`/api/games/latest/turns/${turnCount}`);
  const turn: TurnResponse = await res.json();
  return turn;
};

const registerTurn = async (turnReq: TurnRequest) => {
  console.log(turnReq);
  const res = await fetchApi("/api/games/latest/turns", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(turnReq),
  });
  const newTurn = await res.json();
  return newTurn;
};

export const Board: React.FC = () => {
  const [board, setBoard] = useState<Disc[][]>(
    Array(8)
      .fill(null)
      .map(() => Array(8).fill(null))
  );
  const [nextDisc, setNextDisc] = useState<Disc>(0);
  const [turnCount, setTurnCount] = useState(0);

  useEffect(() => {
    const fetchTurn = async () => {
      const turn = await getTurn(turnCount);
      setBoard(turn.board);
      setNextDisc(turn.nextDisc);
    };

    fetchTurn();
  }, [turnCount]);

  const renderSquare = (y: number, x: number) => {
    const disc = board[y][x];

    const handleSquareClick = async () => {
      const nextTurnCount = turnCount + 1;
      setTurnCount(nextTurnCount);
      await registerTurn({
        turnCount: nextTurnCount,
        move: {
          disc: nextDisc,
          x,
          y,
        },
      });
    };

    return (
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
    <div className="inline-block bg-green-700 p-1">
      {[...Array(8)].map((_, y) => renderRow(y))}
    </div>
  );
};
