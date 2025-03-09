import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { type GameHistory, GameHistoryTable } from "./GameHistoryTable";
import { PageTitle } from "./PageTitle";
import { StartButton } from "./StartButton";
import { fetchApi } from "./fetch";

type GameHistoryResponse = {
  games: GameHistory[];
};

export const TopPage = () => {
  const [gameHistories, setGameHistories] = useState<GameHistory[]>([
    {
      blackMoveCount: 0,
      whiteMoveCount: 0,
      winnerDisc: 1,
      startedAt: "2021-01-01T00:00:00Z",
      endedAt: "2021-01-01T00:00:00Z",
    },
    {
      blackMoveCount: 0,
      whiteMoveCount: 0,
      winnerDisc: 1,
      startedAt: "2021-01-01T00:00:00Z",
      endedAt: "2021-01-01T00:00:00Z",
    },
  ]);
  const navigate = useNavigate();

  useEffect(() => {
    const getGameHistory = async () => {
      const res = await fetchApi("/api/games");
      const data: GameHistoryResponse = await res.json();
      setGameHistories(data.games);
    };
    getGameHistory();
  }, []);

  return (
    <div className="flex flex-col items-center gap-10">
      <PageTitle title="Othello TopPage" />
      <StartButton onClick={() => navigate("/game")}>Game Start!</StartButton>
      <GameHistoryTable gameHistories={gameHistories} />
    </div>
  );
};
