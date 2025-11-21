let activePower = null;

document.addEventListener("DOMContentLoaded", () => {
  const powerBtns = document.querySelectorAll(".power-btn");
  const cells = document.querySelectorAll(".cell");
  const instruction = document.getElementById("power-instruction");

  powerBtns.forEach(btn => {
    btn.addEventListener("click", () => {
      if (btn.classList.contains("disabled")) return;

      if (btn.classList.contains("active")) {
        deactivateAllPowers();
        return;
      }

      deactivateAllPowers();
      activePower = btn.dataset.power;
      btn.classList.add("active");
      
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
        
        const plop = new Audio("/static/plop.wav");
        plop.volume = 0.6;
        plop.play();

        setTimeout(() => {
          window.location.href = url;
        }, 600);
      } else {
        const url = cell.getAttribute("href");
        const plop = new Audio("/static/plop.wav");
        plop.volume = 0.6;
        plop.play();

        setTimeout(() => {
          window.location.href = url;
        }, 600);
      }
    });
  });

  const status = document.getElementById("status");
  if (status && status.textContent.includes("a gagné")) {
    const win = new Audio("/static/win.wav");
    win.volume = 0.8;
    setTimeout(() => win.play(), 400);
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