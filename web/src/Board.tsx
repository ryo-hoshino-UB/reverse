export const Board: React.FC = () => {
  const renderSquare = (row: number, col: number) => {
    return (
      <div
        key={`${row}-${col}`}
        className="w-12 h-12 bg-green-600 border border-black flex items-center justify-center"
      >
        {/* ここに石を配置するロジックを追加 */}
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
