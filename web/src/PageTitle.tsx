export const PageTitle: React.FC<{ title: string }> = ({ title }) => {
  return (
    <h1 className="flex flex-col items-center tracking-wide text-3xl font-bold text-gradient bg-clip-text text-transparent bg-gradient-to-r from-emerald-400 via-teal-500 to-blue-600 mb-6 py-2 border-b-2 border-gray-300 shadow-sm">
      {title}
    </h1>
  );
};
