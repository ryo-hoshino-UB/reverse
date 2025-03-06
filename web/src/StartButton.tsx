import type { ComponentProps } from "react";

type Props = ComponentProps<"button">;
export const StartButton: React.FC<Props> = ({ ...props }) => {
  return (
    <button
      type="button"
      onClick={props.onClick}
      className="bg-blue-400 text-white "
    >
      Game Start!
    </button>
  );
};
