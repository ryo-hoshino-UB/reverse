import { Stone } from "./Stone";

export const Board: React.FC = () => {
  const handleClick = () => {
    console.log("クリックされました");
  };

  const renderSquare = (row: number, col: number) => {
    return (
      <div
        onClick={handleClick}
        key={`${row}-${col}`}
        className="w-12 h-12 bg-green-600 border border-black flex items-center justify-center"
      >
        <Stone disc={1} />
      </div>
    );
  };

  const renderRow = (row: number) => {
    return (
      <div key={row} className="flex">
        {[...Array(8)].map((_, col) => renderSquare(row, col))}
      </div>
    );
  };

  return (
    <div className="inline-block bg-green-700 p-1">
      {[...Array(8)].map((_, row) => renderRow(row))}
    </div>
  );
};
