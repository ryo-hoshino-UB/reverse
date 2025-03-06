import type { Disc } from "./disc";

type Props = {
  nextDisc: Disc;
};
export const NextDiscBanner: React.FC<Props> = ({ nextDisc }) => {
  const discColor = nextDisc === 2 ? "白" : "黒";
  return (
    <div className="bg-blue-400 text-white font-bold  p-2 w-full rounded text-center tracking-widest">
      <span className={`${discColorStyle(nextDisc)} font-bold`}>
        {discColor}の番です
      </span>
    </div>
  );
};

const discColorStyle = (disc: Disc): string => {
  return disc === 2 ? "text-white" : "text-gray-700";
};
