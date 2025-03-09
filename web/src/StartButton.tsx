import type { ComponentProps, PropsWithChildren } from "react";

type Props = ComponentProps<"button">;
export const StartButton: React.FC<PropsWithChildren<Props>> = ({
  ...props
}) => {
  return (
    <button type="button" onClick={props.onClick} className="text-white ">
      {props.children}
    </button>
  );
};
