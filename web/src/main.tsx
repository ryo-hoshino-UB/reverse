import { createRoot } from "react-dom/client";
import App from "./App.tsx";
import "./index.css";

createRoot(document.getElementById("root")!).render(
  <>
    <head>
      <meta charSet="utf-8" />
      <title>Othello App</title>
    </head>
    <App />
  </>
);
