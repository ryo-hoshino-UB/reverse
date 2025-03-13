import { Outlet } from "react-router";

export const Layout: React.FC = () => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 to-gray-800 text-white w-full">
      <div className="min-w-screen mx-auto px-12 py-12">
        <Outlet />
      </div>
    </div>
  );
};
