import { Suspense, useState } from "react";
import "./App.css";
import { Board } from "./Board";
import { fetchApi } from "./fetch";

const registerGame = async () => {
  const res = await fetchApi("/api/games", {
    method: "POST",
  });
  const game = await res.json();
  return game;
};

function App() {
  const [isStarted, setIsStarted] = useState(false);
  const onClick = async () => {
    await registerGame();
    setIsStarted(true);
  };

  return (
    <div className="text-gray-800 flex flex-col gap-10 items-center justify-start w-dvw h-screen">
      <header className="flex flex-row justify-items-center">
        <h1 className=" text-2xl">Othello App</h1>
      </header>
      <main>
        <div className="flex flex-col gap-6 w-[392px]">
          <p className=" bg-red-500 text-white font-bold p-2 w-full rounded text-center">
            白の番はスキップです
          </p>
          {!isStarted ? (
            <button className="text-white" onClick={onClick}>
              Game Start!
            </button>
          ) : (
            <Suspense fallback={<div>Loading...</div>}>
              <Board />
            </Suspense>
          )}
          <p className="bg-blue-500 text-white font-bold p-2 w-full rounded text-center">
            次は黒の番です
          </p>
        </div>
      </main>
    </div>
  );
}

export default App;
