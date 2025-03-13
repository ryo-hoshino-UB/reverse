import { CircleDotIcon, CircleIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { type GameHistory, GameHistoryTable } from "./GameHistoryTable";
import { PageTitle } from "./PageTitle";
import { Button } from "./components/ui/button";
import { fetchApi } from "./fetch";

type GameHistoriesResponse = GameHistory[];

export const TopPage = () => {
  const [gameHistories, setGameHistories] = useState<GameHistory[]>([]);
  const navigate = useNavigate();

  useEffect(() => {
    const getGameHistory = async () => {
      try {
        const res = await fetchApi("/api/games");
        const fetchGameHistories: GameHistoriesResponse = await res.json();
        console.log(fetchGameHistories);
        setGameHistories(fetchGameHistories);
      } catch {
        setGameHistories([]);
      }
    };

    getGameHistory();
  }, []);

  return (
    <div className="relative">
      {/* Decorative elements */}
      <div className="absolute -top-6 -left-6 text-emerald-500/20">
        <CircleIcon size={48} />
      </div>
      <div className="absolute -bottom-6 -right-6 text-emerald-500/20">
        <CircleDotIcon size={48} />
      </div>

      <div className="bg-gray-800/50 backdrop-blur-sm rounded-2xl p-8 border border-gray-700 shadow-xl">
        <div className="flex flex-col items-center gap-12">
          <div className="text-center">
            <PageTitle title="Othello" />
            <p className="mt-2 text-emerald-400 text-lg">
              Challenge your strategic thinking
            </p>
          </div>

          <Button
            type="button"
            onClick={() => navigate("/game")}
            className="bg-emerald-500 hover:bg-emerald-600 text-white rounded-2xl transition-all transform hover:scale-105 shadow-lg hover:shadow-emerald-500/30 w-full h-full max-w-[16rem]"
          >
            <span className="text-xl font-bold">Game Start!</span>
          </Button>

          <div className="w-full bg-gray-900/50 rounded-xl p-6 backdrop-blur-sm border border-gray-700">
            <h2 className="text-xl font-semibold mb-6 text-emerald-400 flex items-center gap-2">
              <CircleDotIcon size={20} />
              Game History
            </h2>
            <GameHistoryTable gameHistories={gameHistories} />
          </div>
        </div>
      </div>
    </div>
  );
};
