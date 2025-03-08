export const WINNER_DISC = {
  Draw: 0,
  Black: 1,
  White: 2,
} as const;

export type WinnerDisc = (typeof WINNER_DISC)[keyof typeof WINNER_DISC];
