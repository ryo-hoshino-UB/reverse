export const DISC = {
  empty: 0,
  black: 1,
  white: 2,
} as const;

const empty = DISC.empty;
const black = DISC.black;
const white = DISC.white;

export const INITIAL_BOARD = [
  [empty, empty, empty, empty, empty, empty, empty, empty],
  [empty, empty, empty, empty, empty, empty, empty, empty],
  [empty, empty, empty, empty, empty, empty, empty, empty],
  [empty, empty, empty, black, white, empty, empty, empty],
  [empty, empty, empty, white, black, empty, empty, empty],
  [empty, empty, empty, empty, empty, empty, empty, empty],
  [empty, empty, empty, empty, empty, empty, empty, empty],
  [empty, empty, empty, empty, empty, empty, empty, empty],
];

export type Disc = (typeof DISC)[keyof typeof DISC];

export type DiscColor = keyof typeof DISC;
