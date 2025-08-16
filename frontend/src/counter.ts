export async function setupCounter() {
  const click = document.querySelector<HTMLButtonElement>("#click-count")!;
  const view = document.querySelector<HTMLButtonElement>("#view-count")!;
  const online = document.querySelector<HTMLButtonElement>("#online-count")!;
  const resp = await (
    await fetch(`${import.meta.env.VITE_API_URL}/count`)
  ).json();
  click.innerHTML = `(db) count is ${resp.click}`;
  view.innerHTML = resp.view;
  click.addEventListener("click", async () => {
    const resp = await (
      await fetch(`${import.meta.env.VITE_API_URL}/count-click`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          click: parseInt(click.innerHTML.match(/\d+/)![0]) + 1,
        }),
      })
    ).json();
    click.innerHTML = `count is ${resp.click}`;
  });
  const ws = new WebSocket(`${import.meta.env.VITE_API_URL}/ws`);
  ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    online.innerHTML = data.count;
  };
  ws.onopen = () => {
    console.log("connected to websocket");
  };
  ws.onclose = () => {
    console.log("disconnected from websocket");
  };
  ws.onerror = (event) => {
    console.log("websocket error", event);
  };
}
