async function fetchState() {
  const res = await fetch("/state");
  return await res.json();
}

async function playMove(col) {
  await fetch("/play?column=" + col);
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
      if (state.board[r][c] === 1) cell.classList.add("player1");
      if (state.board[r][c] === 2) cell.classList.add("player2");
      cell.addEventListener("click", () => playMove(c));
      boardEl.appendChild(cell);
    }
  }

  if (state.winner !== 0) {
    statusEl.textContent = `🎉 Joueur ${state.winner} a gagné !`;
  } else {
    statusEl.textContent = `Tour du joueur ${state.current}`;
  }
}

document.getElementById("restart").addEventListener("click", async () => {
  // simple reload trick
  await fetch("/play?column=-1"); // inutile mais évite blocage
  location.reload();
});

updateBoard();

document.getElementById("restart").addEventListener("click", async () => {
  await fetch("/reset");
  updateBoard();
});

if (state.board[r][c] === 1) {
  cell.classList.add("player1");
  cell.style.animation = "none";
  setTimeout(() => (cell.style.animation = ""), 10);
}
if (state.board[r][c] === 2) {
  cell.classList.add("player2");
  cell.style.animation = "none";
  setTimeout(() => (cell.style.animation = ""), 10);
}
