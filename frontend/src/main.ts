import "./style.css";
import typescriptLogo from "./typescript.svg";
import viteLogo from "/vite.svg";
import stagewiseLogo from "/stagewise-logo.png";
import { setupCounter } from "./counter.ts";

// Initialize stagewise toolbar in development mode
if (import.meta.env.DEV) {
  import('@stagewise/toolbar').then(({ initToolbar }) => {
    initToolbar({
      plugins: [],
    });
  });
}

document.querySelector<HTMLDivElement>("#app")!.innerHTML = `
  <div>
    <a href="https://vite.dev" target="_blank">
      <img src="${viteLogo}" class="logo" alt="Vite logo" />
    </a>
    <a href="https://www.typescriptlang.org/" target="_blank">
      <img src="${typescriptLogo}" class="logo vanilla" alt="TypeScript logo" />
    </a>
    <a href="https://stagewise.io" target="_blank">
      <img src="${stagewiseLogo}" class="logo" alt="Stagewise logo" />
    </a>
    <h1>Vite + TypeScript + Stagewise</h1>
    <div class="card">
      <button id="counter" type="button"></button>
    </div>
    <p class="read-the-docs">
      Click on the logos to learn more
    </p>
    <p class="read-the-docs">
      Version: <span id="version">__VERSION__</span>
    </p>
  </div>
`;

window.addEventListener("load", async () => {
  document.getElementById("version")!.textContent = await (
    await fetch(`${import.meta.env.VITE_API_URL}/version`)
  ).text();
  setupCounter(document.querySelector<HTMLButtonElement>("#counter")!);
});
