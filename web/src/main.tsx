import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./App.tsx";
import "./index.css";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <head>
      <meta charSet="utf-8" />
      <title>Othello App</title>
    </head>
    <App />
  </StrictMode>
);
