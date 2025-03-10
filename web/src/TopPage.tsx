import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { type GameHistory, GameHistoryTable } from "./GameHistoryTable";
import { PageTitle } from "./PageTitle";
import { StartButton } from "./StartButton";
import { fetchApi } from "./fetch";

type GameHistoriesResponse = GameHistory[];

export const TopPage = () => {
  const [gameHistories, setGameHistories] = useState<GameHistory[]>([]);
  const navigate = useNavigate();

  useEffect(() => {
    const getGameHistory = async () => {
      try {
        const res = await fetchApi("/api/games");
        const gameHistories: GameHistoriesResponse = await res.json();
        setGameHistories(gameHistories);
      } catch {
        setGameHistories([]);
      }
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
