export const DISC = {
  empty: 0,
  black: 1,
  white: 2,
} as const;

export type Disc = (typeof DISC)[keyof typeof DISC];

export type DiscColor = keyof typeof DISC;
