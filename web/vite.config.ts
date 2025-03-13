import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react-swc";
import path from "node:path";
import { defineConfig } from "vite";

// https://vite.dev/config/
export default defineConfig({
  base: "./",
  root: "src",
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  plugins: [react(), tailwindcss()],
  build: {
    outDir: "dist", // ビルド出力ディレクトリを指定
    // 存在しないときはフォルダを作成する
    emptyOutDir: true,
    copyPublicDir: true,
    rollupOptions: {
      input: {
        main: "src/index.html",
      },
    },
  },
  server: {
    port: 5001,
  },
  preview: {
    port: 5001,
  },
});
