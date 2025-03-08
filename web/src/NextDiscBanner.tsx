type Props = {
  message: string;
};
export const NextDiscBanner: React.FC<Props> = ({ message }) => {
  return (
    <div className="bg-blue-400 text-white font-bold  p-2 w-full rounded text-center tracking-widest">
      <span className="text-white font-bold">{message}</span>
    </div>
  );
};
