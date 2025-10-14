document.addEventListener("DOMContentLoaded", () => {
  // ---- Effet sonore "plop" quand on clique pour jouer un pion ----
  document.querySelectorAll(".cell").forEach(cell => {
    cell.addEventListener("click", () => {
      const plop = new Audio("/static/plop.wav");
      plop.volume = 0.6;
      plop.play();
    });
  });

  // ---- Si un joueur a gagné, jouer le son de victoire ----
  const status = document.getElementById("status");
  if (status && status.textContent.includes("a gagné")) {
    const win = new Audio("/static/win.wav");
    win.volume = 0.8;
    setTimeout(() => win.play(), 400); // petit délai pour le timing
  }
});
