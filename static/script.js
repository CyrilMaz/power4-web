let activePower = null;

function deactivateAllPowers() {
  activePower = null;
  const powerBtns = document.querySelectorAll('.power-btn');
  powerBtns.forEach(btn => btn.classList.remove('active'));
  const cells = document.querySelectorAll('.cell');
  cells.forEach(cell => cell.classList.remove('power-target'));
  const instruction = document.getElementById('power-instruction');
  if (instruction) instruction.classList.add('hidden');
}

function markColumns() {
  const cols = {};
  document.querySelectorAll('.cell').forEach(cell => {
    const col = cell.dataset.col;
    const row = parseInt(cell.dataset.row, 10);
    if (!cols[col]) cols[col] = [];
    cols[col].push({el: cell, row});
  });
  Object.keys(cols).forEach(col => {
    cols[col].sort((a,b)=>a.row - b.row);
    const top = cols[col][0].el;
    const isFull = top.classList.contains('player1') || top.classList.contains('player2') || top.classList.contains('blocked');
    if (isFull) {
      cols[col].forEach(c => {
        c.el.classList.add('col-full');
        c.el.setAttribute('aria-disabled','true');
        c.el.style.cursor = 'not-allowed';
      });
    }
  });
}

function highlightValidTargets(power) {
  // règles simples :
  // - Détruire : cible une pièce adverse (1 case contenant l'autre joueur)
  // - Échanger : cible une pièce qui a une pièce en dessous non vide et non bloc
  // - Bloquer : cible une colonne (on utilisera la colonne via data-col)

  const cells = document.querySelectorAll('.cell');
  cells.forEach(cell => {
    cell.classList.remove('power-target');
    // ignorer les cases bloquées
    if (cell.classList.contains('blocked')) return;
    if (power === 'Détruire') {
      // target si pièce adverse
      if (cell.classList.contains('player1') || cell.classList.contains('player2')) {
        cell.classList.add('power-target');
      }
    } else if (power === 'Échanger') {
      const row = parseInt(cell.dataset.row, 10);
      const col = cell.dataset.col;
      const below = document.querySelector(`.cell[data-row='${row+1}'][data-col='${col}']`);
      if (below && !below.classList.contains('blocked') && (cell.classList.contains('player1') || cell.classList.contains('player2')) && !below.classList.contains('player1') && !below.classList.contains('player2')) {
        // ne pas permettre échange si dessous est vide ou bloc
      }
      // simplification: permettre cible si elle et la suivante sont non vides et non blocked
      if (below && !cell.classList.contains('blocked') && !below.classList.contains('blocked') && ( (cell.classList.contains('player1')||cell.classList.contains('player2')) && (below.classList.contains('player1')||below.classList.contains('player2')))) {
        cell.classList.add('power-target');
      }
    } else if (power === 'Bloquer') {
      // on surligne le top de chaque colonne vide (cell with highest row that is empty)
      // simple approche : surligner toutes les cellules vides qui ont en dessous soit un jeton, soit sont au fond
      if (!cell.classList.contains('player1') && !cell.classList.contains('player2') && !cell.classList.contains('blocked')) {
        cell.classList.add('power-target');
      }
    }
  });
}


document.addEventListener('DOMContentLoaded', () => {
  const powerBtns = document.querySelectorAll('.power-btn');
  const cells = document.querySelectorAll('.cell');
  const instruction = document.getElementById('power-instruction');

  markColumns();

  powerBtns.forEach(btn => {
    btn.addEventListener('click', (e) => {
      e.preventDefault();
      e.stopPropagation();
      if (btn.classList.contains('disabled')) return;
      if (btn.classList.contains('active')) {
        deactivateAllPowers();
        return;
      }
      deactivateAllPowers();
      activePower = btn.dataset.power;
      btn.classList.add('active');
      if (instruction) instruction.classList.remove('hidden');
      highlightValidTargets(activePower);
    });
  });

  cells.forEach(cell => {
    cell.addEventListener('click', event => {
      event.preventDefault();
      // si colonne bloquée ou pleine -> ignore
      if (cell.getAttribute('aria-disabled') === 'true') return;

      if (activePower) {
        // si cible n'est pas marquée comme valide, ignore
        if (!cell.classList.contains('power-target')) return;
        const row = cell.dataset.row;
        const col = cell.dataset.col;
        const url = `/power?power=${encodeURIComponent(activePower)}&row=${encodeURIComponent(row)}&col=${encodeURIComponent(col)}`;
        const plop = new Audio('/static/plop.wav');
        plop.volume = 0.6;
        plop.play().catch(()=>{});
        // empêcher d'autres clics
        document.querySelectorAll('.cell').forEach(c=>c.style.pointerEvents='none');
        setTimeout(()=> window.location.href = url, 300);
      } else {
        const url = cell.getAttribute('href');
        const plop = new Audio('/static/plop.wav');
        plop.volume = 0.6;
        plop.play().catch(()=>{});
        document.querySelectorAll('.cell').forEach(c=>c.style.pointerEvents='none');
        setTimeout(()=> window.location.href = url, 300);
      }
    });
  });

  const status = document.getElementById('status');
  if (status && status.textContent.includes('a gagné')) {
    const win = new Audio('/static/win.wav');
    win.volume = 0.8;
    setTimeout(()=> win.play().catch(()=>{}), 400);
  }
});