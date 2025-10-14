document.addEventListener("DOMContentLoaded", () => {
  // Toutes les cellules du plateau
  document.querySelectorAll(".cell").forEach(cell => {
    cell.addEventListener("click", event => {
      event.preventDefault(); // ❌ empêche la redirection immédiate

      const url = cell.getAttribute("href"); // garde l’URL du coup
      const plop = new Audio("/static/plop.wav");
      plop.volume = 0.6;
      plop.play();

      // ✅ redirection après 150 ms (temps de jouer le son)
      setTimeout(() => {
        window.location.href = url;
      }, 450);
    });
  });

  // Si un joueur a gagné → jouer son de victoire
  const status = document.getElementById("status");
  if (status && status.textContent.includes("a gagné")) {
    const win = new Audio("/static/win.wav");
    win.volume = 0.8;
    setTimeout(() => win.play(), 400);
  }
});
