"use client";

import { AlertCircle, Home } from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { Stone } from "./Stone";
import type { Disc } from "./disc";
import { fetchApi } from "./fetch";
import type { TurnRequest, TurnResponse } from "./interface";
import { WINNER_DISC } from "./winnerDisc";

const registerGame = async () => {
  const res = await fetchApi("/api/games", {
    method: "POST",
  });
  const game = await res.json();
  return game;
};

const getTurn = async (turnCount: number) => {
  const res = await fetchApi(`/api/games/latest/turns/${turnCount}`);
  const turn: TurnResponse = await res.json();
  return turn;
};

const registerTurn = async (turnReq: TurnRequest) => {
  const res = await fetchApi("/api/games/latest/turns", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(turnReq),
  });
  if (!res.ok) {
    throw new Error("Failed to register turn");
  }
  const newTurn = await res.json();
  return newTurn;
};

export const Board: React.FC = () => {
  const [nextDisc, setNextDisc] = useState<Disc>(0);
  const [turnCount, setTurnCount] = useState(0);
  const [board, setBoard] = useState<Disc[][]>(
    Array(8)
      .fill(null)
      .map(() => Array(8).fill(null))
  );
  const [prevBoard, setPrevBoard] = useState<Disc[][]>(
    Array(8)
      .fill(null)
      .map(() => Array(8).fill(null))
  );
  const [flippingCells, setFlippingCells] = useState<boolean[][]>(
    Array(8)
      .fill(null)
      .map(() => Array(8).fill(false))
  );
  const [newPlacedCell, setNewPlacedCell] = useState<{
    x: number;
    y: number;
  } | null>(null);
  const [bannerMessage, setBannerMessage] = useState(
    changeBannerMessage({ nextDisc: 1 })
  );
  const [gameOver, setGameOver] = useState(false);
  const [winnerDisc, setWinnerDisc] = useState<number | null>(null);
  const [showConfirmModal, setShowConfirmModal] = useState(false);
  const [animationsInProgress, setAnimationsInProgress] = useState(false);
  const [showWinnerOverlay, setShowWinnerOverlay] = useState(false);

  // アニメーション用のスタイルを追加
  useEffect(() => {
    const style = document.createElement("style");
    style.innerHTML = `
      @keyframes flip {
        0% {
          transform: perspective(400px) rotateY(0deg);
        }
        50% {
          transform: perspective(400px) rotateY(90deg);
        }
        100% {
          transform: perspective(400px) rotateY(0deg);
        }
      }
      
      @keyframes place {
        0% {
          transform: scale(0);
          opacity: 0;
        }
        100% {
          transform: scale(1);
          opacity: 1;
        }
      }
      
      .animate-flip {
        animation: flip 300ms ease-in-out forwards;
        transform-style: preserve-3d;
      }
      
      .animate-place {
        animation: place 300ms cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
      }
    `;
    document.head.appendChild(style);

    return () => {
      document.head.removeChild(style);
    };
  }, []);

  const navigate = useNavigate();

  const handleHomeClick = () => {
    if (gameOver) {
      navigate("/");
      return;
    }
    setShowConfirmModal(true);
  };

  const handleCancelClick = () => {
    setShowConfirmModal(false);
  };

  const handleConfirmClick = () => {
    setShowConfirmModal(false);
    navigate("/");
  };

  // 前回のボードと現在のボードを比較して、ひっくり返った石を特定
  useEffect(() => {
    if (prevBoard.length === 0 || board.length === 0) return;

    // 既に値が入っていて値が変わったセルをフリップ対象とする
    const newFlippingCells = Array(8)
      .fill(null)
      .map(() => Array(8).fill(false));
    let hasFlippingCells = false;

    for (let y = 0; y < 8; y++) {
      for (let x = 0; x < 8; x++) {
        if (
          prevBoard[y][x] !== 0 &&
          board[y][x] !== 0 &&
          prevBoard[y][x] !== board[y][x]
        ) {
          newFlippingCells[y][x] = true;
          hasFlippingCells = true;
        }
      }
    }

    if (hasFlippingCells) {
      setAnimationsInProgress(true);
      setFlippingCells(newFlippingCells);

      // アニメーション終了後にリセット
      const timer = setTimeout(() => {
        setFlippingCells(
          Array(8)
            .fill(null)
            .map(() => Array(8).fill(false))
        );
      }, 300);
      setAnimationsInProgress(false);

      if (gameOver) {
        setShowWinnerOverlay(true);
      }

      return () => clearTimeout(timer);
    }
  }, [board, prevBoard, gameOver]);

  useEffect(() => {
    const startGame = async () => {
      await registerGame();
    };

    const fetchTurn = async () => {
      const turn = await getTurn(0);
      setBoard(turn.board);
      setPrevBoard(turn.board);
      setNextDisc(turn.nextDisc);
    };

    startGame();
    fetchTurn();
  }, []);

  const renderConfirmModal = () => {
    if (!showConfirmModal) return null;

    return (
      <div className="fixed inset-0 flex items-center justify-center z-50 bg-black/40 backdrop-blur-sm">
        <div className="bg-gray-800 border border-gray-700 rounded-xl shadow-xl p-6 max-w-md w-full mx-4 animate-fadeIn">
          <div className="flex items-start gap-4">
            <div className="text-amber-400 flex-shrink-0 mt-1">
              <AlertCircle size={24} />
            </div>
            <div className="flex-grow">
              <h3 className="text-lg font-semibold text-white mb-2">
                対戦を終了しますか？
              </h3>
              <p className="text-gray-300 mb-6">
                ホームに戻ると現在の対戦が終了します。よろしいですか？
              </p>
              <div className="flex justify-end gap-3">
                <button
                  type="button"
                  onClick={handleCancelClick}
                  className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors duration-200"
                >
                  キャンセル
                </button>
                <button
                  type="button"
                  onClick={handleConfirmClick}
                  className="px-4 py-2 bg-emerald-600 hover:bg-emerald-500 text-white rounded-lg transition-colors duration-200"
                >
                  ホームに戻る
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  };

  const renderSquare = (y: number, x: number) => {
    const disc = board[y][x];
    const isFlipping = flippingCells[y][x];
    const isNewPlaced = newPlacedCell?.x === x && newPlacedCell?.y === y;

    const handleSquareClick = async () => {
      // ゲーム終了時はクリックを無効化
      if (gameOver || animationsInProgress) return;

      const nextTurnCount = turnCount + 1;

      try {
        // クリックされたセルを記録（新しく置かれる石の場所）
        setNewPlacedCell({ x, y });
        setAnimationsInProgress(true);

        // 現在のボードを保存（変更前の状態）
        setPrevBoard([...board.map((row) => [...row])]);

        await registerTurn({
          turnCount: nextTurnCount,
          move: {
            disc: nextDisc,
            x,
            y,
          },
        });

        // registerTurnの後に直接最新状態を取得
        const turn = await getTurn(nextTurnCount);

        if (turn.nextDisc !== 0) {
          setBannerMessage(changeBannerMessage({ nextDisc: turn.nextDisc }));
          if (nextDisc === turn.nextDisc) {
            setBannerMessage(
              `もう一度${nextDisc === 1 ? "白" : "黒"}のターンです`
            );
          }
        } else {
          // ゲーム終了時の処理
          setGameOver(true);
          setWinnerDisc(turn.winnerDisc);
          setBannerMessage("ゲーム終了");
        }

        // ボード状態を更新
        setBoard(turn.board);
        setNextDisc(turn.nextDisc);
        setTurnCount(nextTurnCount);

        // 一定時間後に新しく置いた石のマーキングをリセット
        setTimeout(() => {
          setNewPlacedCell(null);
          setAnimationsInProgress(false);

          // ゲーム終了時、アニメーション終了後にオーバーレイを表示
          if (turn.nextDisc === 0) {
            setShowWinnerOverlay(true);
          }
        }, 300);
      } catch (error) {
        console.error("Error placing disc:", error);
        setNewPlacedCell(null);
      }
    };

    return (
      // biome-ignore lint/a11y/useKeyWithClickEvents:
      <div
        onClick={handleSquareClick}
        key={`${y}-${x}`}
        className="w-14 h-14 bg-gradient-to-br from-emerald-700 to-emerald-800 border border-emerald-900/70 flex items-center justify-center relative overflow-hidden transition-all duration-300 hover:shadow-inner hover:from-emerald-600 hover:to-emerald-700"
      >
        {/* セルの光沢効果 */}
        <div className="absolute top-0 left-0 right-0 h-[1px] bg-emerald-400/30" />
        <div className="absolute top-0 left-0 bottom-0 w-[1px] bg-emerald-400/30" />
        <div className="absolute bottom-0 left-0 right-0 h-[1px] bg-emerald-900/50" />
        <div className="absolute top-0 left-0 bottom-0 w-[1px] bg-emerald-400/30" />
        <Stone disc={disc} isFlipping={isFlipping || isNewPlaced} />
      </div>
    );
  };

  const renderRow = (y: number) => {
    return (
      <div key={y} className="flex">
        {[...Array(8)].map((_, x) => renderSquare(y, x))}
      </div>
    );
  };

  // 勝者表示用のオーバーレイ
  const renderWinnerOverlay = () => {
    if (!showWinnerOverlay) return null;

    let winnerContent: React.ReactNode;
    const blackCount = board.flat().filter((d) => d === 1).length;
    const whiteCount = board.flat().filter((d) => d === 2).length;

    if (winnerDisc === WINNER_DISC.Draw) {
      winnerContent = (
        <div className="text-center">
          <h2 className="text-3xl font-bold text-white mb-4">引き分け</h2>
          <div className="flex justify-center items-center gap-8">
            <div className="flex flex-col items-center">
              <div className="w-12 h-12 rounded-full bg-gray-800 shadow-[0_0_10px_rgba(0,0,0,0.5)] mb-2" />
              <span className="text-xl font-semibold text-white">
                {blackCount}
              </span>
            </div>
            <div className="flex flex-col items-center">
              <div className="w-12 h-12 rounded-full bg-white shadow-[0_0_10px_rgba(255,255,255,0.5)] mb-2" />
              <span className="text-xl font-semibold text-white">
                {whiteCount}
              </span>
            </div>
          </div>
        </div>
      );
    } else {
      winnerContent = (
        <div className="text-</div>center">
          <h2 className="text-3xl font-bold text-white mb-4">
            {winnerDisc === 1 ? "黒の勝ち！" : "白の勝ち！"}
          </h2>
          <div className="flex justify-center items-center gap-8">
            <div className="flex flex-col items-center">
              <div className="w-12 h-12 rounded-full bg-gray-800 shadow-[0_0_10px_rgba(0,0,0,0.5)] mb-2" />
              <span className="text-xl font-semibold text-white">
                {blackCount}
              </span>
            </div>
            <div className="flex flex-col items-center">
              <div className="w-12 h-12 rounded-full bg-white shadow-[0_0_10px_rgba(255,255,255,0.5)] mb-2" />
              <span className="text-xl font-sem</div>ibold text-white">
                {whiteCount}
              </span>
            </div>
          </div>

          <button
            type="button"
            onClick={() => navigate("/")}
            className="mt-8 px-6 py-2 bg-emerald-600 hover:bg-emerald-500 text-white font-medium rounded-full transition-colors duration-200 shadow-lg flex items-center justify-center gap-2"
          >
            <Home size={18} />
            ホームに戻る
          </button>
        </div>
      );
    }

    return (
      <div className="absolute inset-0 bg-bl</div>ack/70 backdrop-blur-sm flex items-center justify-center z-10 rounded-lg animate-fadeIn">
        <div className="bg-gray-800/90 p-8 rounded-xl border border-emerald-400/30 shadow-[0_0_25px_rgba(16,185,129,0.3)]">
          {winnerContent}
        </div>
      </div>
    );
  };

  return (
    <div className="flex flex-col items-center gap-12">
      <div className="flex flex-col items-end gap-4">
        <Home
          size={36}
          className="text-emerald-700 hover:text-emerald-500 cursor-pointer"
          onClick={handleHomeClick}
        />
        <div className="inline-block bg-gradient-to-br from-emerald-800 to-emerald-900 p-1.5 rounded-lg shadow-[0_0_15px_rgba(16,185,129,0.2)] border border-emerald-700/30 relative">
          {[...Array(8)].map((_, y) => renderRow(y))}
          {renderWinnerOverlay()}
        </div>
      </div>

      <div className="bg-gray-800/90 px-6 py-3 rounded-full shadow-lg border border-emerald-500/30 backdrop-blur-sm">
        <div className="relative overflow-hidden">
          <span className="text-white font-bold text-lg tracking-wide relative z-10 flex items-center justify-center">
            {bannerMessage}
          </span>
        </div>
      </div>

      {/* ゲーム情報 */}
      <div className="flex justify-between w-full bg-gray-800/80 rounded-lg p-4 border border-gray-700/50 shadow-md text-base">
        <div className="flex items-center gap-3">
          <div className="w-5 h-5 rounded-full bg-black shadow-[0_0_5px_rgba(0,0,0,0.5)]" />
          <span className="text-white font-medium">
            黒: {board.flat().filter((d) => d === 1).length}
          </span>
        </div>
        <div className="flex items-center gap-3">
          <span className="text-white font-medium">
            白: {board.flat().filter((d) => d === 2).length}
          </span>
          <div className="w-5 h-5 rounded-full bg-white shadow-[0_0_5px_rgba(255,255,255,0.5)]" />
        </div>
      </div>
      {renderConfirmModal()}
    </div>
  );
};

type Props = {
  nextDisc: Disc;
  skip?: boolean;
};

const changeBannerMessage = ({ nextDisc }: Props): string => {
  let nextDiscColor = "";
  if (nextDisc === 1) {
    nextDiscColor = "黒";
  }
  if (nextDisc === 2) {
    nextDiscColor = "白";
  }

  return `${nextDiscColor}のターンです`;
};
