type Props = {
  color: "black" | "white" | "empty";
};

export const Stone: React.FC<Props> = ({ color }) => {
  let stoneClass = "w-10 h-10 rounded-full ";

  switch (color) {
    case "black":
      stoneClass += "bg-black";
      break;
    case "white":
      stoneClass += "bg-white";
      break;
    default:
      stoneClass += "bg-transparent";
  }

  return (
    <div className="w-full h-full flex items-center justify-center">
      <div className={stoneClass}></div>
    </div>
  );
};
