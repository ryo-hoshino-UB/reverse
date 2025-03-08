import type { Disc } from "./disc";
import type { WinnerDisc } from "./winnerDisc";

export type TurnResponse = {
  turnCount: number;
  board: Disc[][];
  nextDisc: Disc;
  winnerDisc: WinnerDisc;
};

export type TurnRequest = {
  turnCount: number;
  move: {
    disc: Disc;
    x: number;
    y: number;
  };
};
