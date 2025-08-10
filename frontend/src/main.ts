import "./style.css";
import typescriptLogo from "./typescript.svg";
import viteLogo from "/vite.svg";
import { setupCounter } from "./counter.ts";

document.querySelector<HTMLDivElement>("#app")!.innerHTML = `
  <div>
    <a href="https://vite.dev" target="_blank">
      <img src="${viteLogo}" class="logo" alt="Vite logo" />
    </a>
    <a href="https://www.typescriptlang.org/" target="_blank">
      <img src="${typescriptLogo}" class="logo vanilla" alt="TypeScript logo" />
    </a>
    <h1>Vite + TypeScript</h1>
    <div class="card">
      <button id="click-count" type="button"></button>
    </div>
    <p class="read-the-docs">
      Version: <span id="version">__VERSION__</span>
    </p>
    <div class="stats">
      <p>
        Online users: <span id="online-count">0</span>
      </p>
      <p>
        Views since last restart: <span id="view-count">0</span>
      </p>
    </div>
  </div>
`;

window.addEventListener("load", async () => {
  document.getElementById("version")!.textContent = await (
    await fetch(`${import.meta.env.VITE_API_URL}/version`)
  ).text();
  setupCounter();
});
