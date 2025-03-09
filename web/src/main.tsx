import { createRoot } from "react-dom/client";
import { BrowserRouter, Route, Routes } from "react-router";
import { GamePage } from "./GamePage.tsx";
import "./index.css";
import { Layout } from "./Layout.tsx";
import { TopPage } from "./TopPage.tsx";

//biome-ignore lint/style/noNonNullAssertion:
createRoot(document.getElementById("root")!).render(
  <>
    <head>
      <meta charSet="utf-8" />
      <title>Othello App</title>
    </head>
    <main>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<TopPage />} />
            <Route path="/game" element={<GamePage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </main>
  </>
);
