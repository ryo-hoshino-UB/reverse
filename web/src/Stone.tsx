import type { Disc } from "./disc";

type Props = {
  disc: Disc;
};

export const Stone: React.FC<Props> = ({ disc }) => {
  let discStyle = "";
  switch (disc) {
    case 1:
      discStyle =
        "bg-gradient-to-br from-gray-900 to-gray-700 shadow-[inset_0_-2px_4px_rgba(255,255,255,0.1)]";
      break;
    case 2:
      discStyle =
        "bg-gradient-to-br from-gray-100 to-gray-300 shadow-[inset_0_-2px_4px_rgba(0,0,0,0.1)]";
      break;
    default:
      break;
  }
  return (
    <div
      className={`w-12 h-12 rounded-full flex items-center justify-center
        ${discStyle} transition-all duration-300 transform`}
    />
  );
};
