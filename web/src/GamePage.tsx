import "./App.css";
import { Board } from "./Board";
import { PageTitle } from "./PageTitle";

export const GamePage = () => {
  return (
    <>
      <PageTitle title="Othello" />
      <Board />
    </>
  );
};
