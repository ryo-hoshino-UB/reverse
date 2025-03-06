import "./App.css";
import { Board } from "./Board";
import { GameTitle } from "./GameTitle";

function App() {
  return (
    <main>
      <div className="text-gray-800 flex flex-col gap-10 items-center justify-start w-dvw h-screen">
        <GameTitle />
        <Board />
      </div>
    </main>
  );
}

export default App;
