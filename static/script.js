document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll(".cell").forEach(cell => {
    cell.addEventListener("click", event => {
      event.preventDefault();

      const url = cell.getAttribute("href");
      const plop = new Audio("/static/plop.wav");
      plop.volume = 0.6;
      plop.play();

      setTimeout(() => {
        window.location.href = url;
      }, 600);
    });
  });

  const status = document.getElementById("status");
  if (status && status.textContent.includes("a gagné")) {
    const win = new Audio("/static/win.wav");
    win.volume = 0.8;
    setTimeout(() => win.play(), 400);
  }
});
