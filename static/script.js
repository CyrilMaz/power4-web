let activePower = null;

document.addEventListener("DOMContentLoaded", () => {
  const powerBtns = document.querySelectorAll(".power-btn");
  const cells = document.querySelectorAll(".cell");
  const instruction = document.getElementById("power-instruction");

  console.log("Boutons de pouvoir trouvés:", powerBtns.length);
  console.log("Cellules trouvées:", cells.length);

  powerBtns.forEach(btn => {
    btn.addEventListener("click", (e) => {
      e.preventDefault();
      e.stopPropagation();
      
      console.log("Bouton cliqué:", btn.dataset.power);

      if (btn.classList.contains("disabled")) {
        console.log("Bouton désactivé");
        return;
      }

      if (btn.classList.contains("active")) {
        console.log("Désactivation du pouvoir");
        deactivateAllPowers();
        return;
      }

      deactivateAllPowers();
      activePower = btn.dataset.power;
      btn.classList.add("active");
      
      console.log("Pouvoir activé:", activePower);
      
      if (instruction) {
        instruction.classList.remove("hidden");
      }

      cells.forEach(cell => cell.classList.add("power-target"));
    });
  });

  cells.forEach(cell => {
    cell.addEventListener("click", event => {
      event.preventDefault();

      if (activePower) {
        const row = cell.dataset.row;
        const col = cell.dataset.col;
        const url = `/power?power=${activePower}&row=${row}&col=${col}`;
        
        console.log("Utilisation du pouvoir:", activePower, "sur", row, col);
        
        const plop = new Audio("/static/plop.wav");
        plop.volume = 0.6;
        plop.play().catch(() => console.log("Pas de son"));

        setTimeout(() => {
          window.location.href = url;
        }, 300);
      } else {
        const url = cell.getAttribute("href");
        const plop = new Audio("/static/plop.wav");
        plop.volume = 0.6;
        plop.play().catch(() => console.log("Pas de son"));

        setTimeout(() => {
          window.location.href = url;
        }, 300);
      }
    });
  });

  const status = document.getElementById("status");
  if (status && status.textContent.includes("a gagné")) {
    const win = new Audio("/static/win.wav");
    win.volume = 0.8;
    setTimeout(() => win.play().catch(() => console.log("Pas de son")), 400);
  }
});

function deactivateAllPowers() {
  activePower = null;
  
  const powerBtns = document.querySelectorAll(".power-btn");
  powerBtns.forEach(btn => btn.classList.remove("active"));
  
  const cells = document.querySelectorAll(".cell");
  cells.forEach(cell => cell.classList.remove("power-target"));
  
  const instruction = document.getElementById("power-instruction");
  if (instruction) {
    instruction.classList.add("hidden");
  }
}