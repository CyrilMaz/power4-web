document.querySelectorAll(".cell").forEach(cell => {
  cell.addEventListener("click", () => {
    const sound = new Audio("/static/plop.mp3");
    sound.volume = 0.5;
    sound.play();
  });
});
