import { useEffect, useState } from "react";
import type { Disc } from "./disc";

type Props = {
  disc: Disc;
  isFlipping?: boolean;
};

export const Stone: React.FC<Props> = ({ disc, isFlipping = false }) => {
  const [flipping, setFlipping] = useState(false);
  const [previousDisc, setPreviousDisc] = useState<Disc>(disc);
  const [displayDisc, setDisplayDisc] = useState<Disc>(disc);

  useEffect(() => {
    // 石の状態が変わった場合（かつ石が存在する場合）、アニメーションを開始
    if (disc !== previousDisc && disc !== 0 && previousDisc !== 0) {
      setFlipping(true);

      // アニメーションの前半で前の石を表示し、後半で新しい石に切り替える
      const halfwayTimer = setTimeout(() => {
        setDisplayDisc(disc);
      }, 150);

      // アニメーション終了時にフリップ状態をリセット
      const endTimer = setTimeout(() => {
        setFlipping(false);
      }, 300);

      return () => {
        clearTimeout(halfwayTimer);
        clearTimeout(endTimer);
      };
    }
    setDisplayDisc(disc);

    setPreviousDisc(disc);
  }, [disc, previousDisc]);

  // 石がない場合は何も表示しない
  if (disc === 0) return null;

  let discStyle = "";
  switch (displayDisc) {
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

  // フリップアニメーションまたは設置アニメーションのクラス
  const animationClass = flipping
    ? "animate-flip"
    : isFlipping
    ? "animate-place"
    : "";

  return (
    <div
      className={`w-12 h-12 rounded-full flex items-center justify-center
        ${discStyle} transition-all duration-300 transform ${animationClass}`}
    />
  );
};
