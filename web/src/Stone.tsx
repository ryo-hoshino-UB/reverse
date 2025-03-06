import type { Disc } from "./disc";

type Props = {
  disc: Disc;
};

export const Stone: React.FC<Props> = ({ disc }) => {
  let stoneClass = "w-10 h-10 rounded-full ";

  switch (disc) {
    case 1:
      stoneClass += "bg-black";
      break;
    case 2:
      stoneClass += "bg-white";
      break;
    default:
      stoneClass += "bg-transparent";
  }

  return (
    <div className="w-full h-full flex items-center justify-center">
      <div className={stoneClass} />
    </div>
  );
};
