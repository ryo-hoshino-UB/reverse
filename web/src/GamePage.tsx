import "./App.css";
import { Board } from "./Board";
import { PageTitle } from "./PageTitle";

export const GamePage = () => {
  return (
    <div className="flex flex-col items-center">
      <div className="mb-8">
        <PageTitle title="Othello" />
      </div>
      <div className="relative w-full max-w-2xl mx-auto">
        <div className="absolute -z-10 -top-20 -left-20 w-64 h-64 bg-emerald-500/10 rounded-full blur-3xl" />
        <div className="absolute -z-10 -bottom-20 -right-20 w-64 h-64 bg-emerald-500/10 rounded-full blur-3xl" />

        <div
          className="backdrop-blur-sm rounded-2xl p-6 border border-emerald-500/70 
            shadow-[0_10px_40px_-15px_rgba(16,185,129,0.3)] 
            bg-gradient-to-br from-gray-900/90 to-gray-800/90 
            hover:translate-y-[-2px] transition-transform duration-300
            dark:shadow-emerald-500/20"
        >
          <Board />
        </div>
      </div>
    </div>
  );
};
