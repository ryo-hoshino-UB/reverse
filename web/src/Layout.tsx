import { Outlet } from "react-router";

export const Layout: React.FC = () => {
  return (
    <div className="text-gray-800 flex flex-col gap-10 items-center justify-start w-dvw h-screen">
      <Outlet />
    </div>
  );
};
