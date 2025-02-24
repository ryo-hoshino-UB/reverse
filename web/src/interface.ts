import { Disc } from "./disc";

export type TurnResponse = {
  turnCount: number;
  board: Disc[][];
  nextDisc: Disc;
  winnerDisc: Disc;
};

export type TurnRequest = {
  turnCount: number;
  move: {
    disc: Disc;
    x: number;
    y: number;
  };
};
