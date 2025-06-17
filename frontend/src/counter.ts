export async function setupCounter(element: HTMLButtonElement) {
  let counter = (
    await (await fetch(`${import.meta.env.VITE_API_URL}/count`)).json()
  ).count;
  const setCounter = (count: number) => {
    counter = count;
    element.innerHTML = `count is ${counter}`;
  };
  const setCounterRemote = async () => {
    await fetch(`${import.meta.env.VITE_API_URL}/count_update`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ count: counter }),
    });
  };
  setCounter(counter);
  element.addEventListener("click", async () => {
    setCounter(counter + 1);
    await setCounterRemote();
  });
}
