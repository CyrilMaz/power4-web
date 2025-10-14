async function fetchState() {
  const res = await fetch("/state");
  return await res.json();
}

async function playMove(col) {
  const sound = new Audio("/static/plop.wav");
  await fetch("/play?column=" + col);
  sound.volume = 0.6;
  sound.play();
  updateBoard();
}

async function updateBoard() {
  const state = await fetchState();
  const boardEl = document.getElementById("board");
  const statusEl = document.getElementById("status");
  boardEl.innerHTML = "";

  for (let r = 0; r < state.board.length; r++) {
    for (let c = 0; c < state.board[r].length; c++) {
      const cell = document.createElement("div");
      cell.className = "cell";

      if (state.board[r][c] === 1 || state.board[r][c] === 2) {
        const playerClass = state.board[r][c] === 1 ? "player1" : "player2";
        cell.classList.add(playerClass);

        const fallDistance = (state.board.length - r);
        cell.style.animation = `drop ${0.1 * fallDistance}s ease-out`;
      }

      cell.addEventListener("click", () => playMove(c));
      boardEl.appendChild(cell);
    }
  }

  if (state.winner !== 0) {
    statusEl.textContent = `🎉 Joueur ${state.winner} a gagné !`;
    celebrateWin(state.winner);
    const winSound = new Audio("/static/win.wav");
    winSound.play();
  } else {
    statusEl.textContent = `Tour du joueur ${state.current}`;
  }
}

function celebrateWin(winner) {
  const cells = document.querySelectorAll(".cell.player" + winner);
  cells.forEach((cell, i) => {
    setTimeout(() => cell.classList.add("winning"), i * 100);
  });
}

document.getElementById("restart").addEventListener("click", async () => {
  await fetch("/reset");
  updateBoard();
});

const themeBtn = document.getElementById("themeToggle");
themeBtn.addEventListener("click", () => {
  document.body.classList.toggle("dark");
  themeBtn.textContent = document.body.classList.contains("dark")
    ? "☀️ Mode clair"
    : "🌙 Mode sombre";
});

updateBoard();
