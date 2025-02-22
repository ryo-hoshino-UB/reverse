import "./App.css";
import { Board } from "./Board";

function App() {
  return (
    <>
      <meta charSet="utf-8" />
      <title>Othello App</title>
      <body className="text-gray-700 flex flex-col gap-10 items-center justify-start h-screen">
        <header className="flex flex-row justify-items-center">
          <h1 className=" text-2xl">Othello App</h1>
        </header>
        <main>
          <div className="flex flex-col gap-6">
            <p className=" bg-red-500 text-white font-bold p-2 w-full rounded-lg text-center">
              白の番はスキップです
            </p>
            <Board />
            <p className="bg-blue-500 text-white font-bold p-2 w-full rounded-lg text-center">
              次は黒の番です
            </p>
          </div>
        </main>
      </body>
    </>
  );
}

export default App;
