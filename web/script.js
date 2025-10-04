const canvas = document.getElementById("canvas");
const ctx = canvas.getContext("2d");
let width = window.innerWidth;
let height = window.innerHeight;
canvas.width = width;
canvas.height = height;

let data = null;
let ants = {}; // {L1: {from, to, progress, color}}
let stepIndex = 0;
let interval = null;
let positions = {};
let animating = false;
let startRoom = "start"; // будет переопределена из данных
let antImage = new Image(); // SVG изображение муравья
let antImageLoaded = false;
const ANT_SIZE = 24; // размер изображения муравья

const baseColors = [
  '#e6194b', '#3cb44b', '#ffe119', '#4363d8', '#f58231',
  '#911eb4', '#46f0f0', '#f032e6', '#bcf60c', '#fabebe',
  '#008080', '#e6beff', '#9a6324', '#fffac8', '#800000'
];

function getColorForAnt(name) {
  const number = parseInt(name.replace(/\D/g, ''));
  if (!isNaN(number) && number <= baseColors.length) {
    return baseColors[number - 1];
  }
  const hash = name.split('').reduce((sum, c) => sum + c.charCodeAt(0), 0);
  const hue = hash % 360;
  return `hsl(${hue}, 80%, 50%)`;
}

// Выбор и надёжная загрузка изображения муравья
function pickAntImageCandidates() {
  const params = new URLSearchParams(window.location.search);
  const userAnt = params.get('ant') || params.get('antImg');
  const candidates = [];
  if (userAnt) {
    // если пользователь передал относительный путь/имя
    candidates.push(userAnt);
    if (!userAnt.includes('/')) {
      candidates.push(`image/${userAnt}`);
      candidates.push(`images/${userAnt}`);
    }
  }
  // дефолтные варианты
  candidates.push('image/ANTS.svg');
  candidates.push('images/ANTS.svg');
  return candidates;
}

function loadAntImageWithFallback() {
  const sources = pickAntImageCandidates();
  let idx = 0;
  const tryNext = () => {
    if (idx >= sources.length) {
      console.error('Ant image failed to load from all candidates, using circle fallback');
      antImageLoaded = false;
      if (data) draw();
      return;
    }
    const src = sources[idx++];
    const img = new Image();
    img.onload = function() {
      antImage = img;
      antImageLoaded = true;
      if (data) draw();
    };
    img.onerror = function() {
      console.warn('Failed to load ant image from', src);
      tryNext();
    };
    img.src = src;
  };
  tryNext();
}

loadAntImageWithFallback();

function computePositions() {
  if (!data || !data.rooms || data.rooms.length === 0) return;
  const allX = data.rooms.map(r => r.x);
  const allY = data.rooms.map(r => r.y);
  const minX = Math.min(...allX), maxX = Math.max(...allX);
  const minY = Math.min(...allY), maxY = Math.max(...allY);
  const scale = Math.min((width - 200) / (maxX - minX + 1), (height - 200) / (maxY - minY + 1));
  positions = {};
  data.rooms.forEach((room) => {
    positions[room.name] = {
      x: (room.x - minX) * scale + 100,
      y: (room.y - minY) * scale + 100,
    };
  });
}

function draw() {
  ctx.clearRect(0, 0, width, height);
  if (!data) return;

  // edges
  data.rooms.forEach((room) => {
    const { x, y } = positions[room.name];
    room.links.forEach((link) => {
      const { x: lx, y: ly } = positions[link];
      ctx.beginPath();
      ctx.moveTo(x, y);
      ctx.lineTo(lx, ly);
      ctx.strokeStyle = "#ccc";
      ctx.lineWidth = 2;
      ctx.stroke();
    });
  });

  // rooms
  data.rooms.forEach((room) => {
    const { x, y } = positions[room.name];
    ctx.beginPath();
    ctx.arc(x, y, 20, 0, Math.PI * 2);
    ctx.fillStyle = room.isStart ? "#4caf50" : room.isEnd ? "#f44336" : "#fff";
    ctx.fill();
    ctx.strokeStyle = "#000";
    ctx.lineWidth = 2;
    ctx.stroke();
    ctx.fillStyle = "#000";
    ctx.font = "14px sans-serif";
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
    ctx.fillText(room.name, x, y);
  });

  // ants
  for (const [name, info] of Object.entries(ants)) {
    const from = positions[info.from];
    const to = positions[info.to];
    const t = info.progress;
    const x = from.x + (to.x - from.x) * t;
    const y = from.y + (to.y - from.y) * t;
    if (antImageLoaded) {
      ctx.drawImage(antImage, x - ANT_SIZE/2, y - ANT_SIZE/2, ANT_SIZE, ANT_SIZE);
    } else {
      ctx.beginPath();
      ctx.arc(x, y, 10, 0, Math.PI * 2);
      ctx.fillStyle = info.color;
      ctx.fill();
    }
    ctx.fillStyle = "#000";
    ctx.font = "10px sans-serif";
    ctx.textAlign = "center";
    ctx.textBaseline = "top";
    ctx.fillText(name, x, y + ANT_SIZE/2 + 5);
  }
}

function animateStep(moves) {
  let frame = 0;
  const frames = 20;
  const movingAnts = {};
  for (const move of moves) {
    const [ant, to] = move.split("-");
    const from = ants[ant]?.to ?? startRoom;
    ants[ant] = {
      from: from,
      to: to,
      progress: 0,
      color: ants[ant]?.color || getColorForAnt(ant),
    };
    movingAnts[ant] = ants[ant];
  }
  animating = true;
  const stepInterval = setInterval(() => {
    frame++;
    for (const ant of Object.values(movingAnts)) {
      ant.progress = frame / frames;
    }
    draw();
    if (frame >= frames) {
      clearInterval(stepInterval);
      animating = false;
    }
  }, 40);
}

function startAnim() {
  if (!data) { console.warn('startAnim: data not loaded yet'); return; }
  if (interval) { console.warn('startAnim: already running'); return; }
  const firstLine = data.moves[0];
  if (firstLine) {
    const firstMoves = firstLine.split(" ");
    for (const move of firstMoves) {
      const [ant] = move.split("-");
      ants[ant] = {
        from: startRoom,
        to: startRoom,
        progress: 1,
        color: getColorForAnt(ant),
      };
    }
    draw();
  }
  interval = setInterval(() => {
    if (animating || stepIndex >= data.moves.length) return;
    const line = data.moves[stepIndex++];
    const parts = line.split(" ");
    document.getElementById("stepCounter").textContent = `Step: ${stepIndex}`;
    animateStep(parts);
    if (stepIndex >= data.moves.length) {
      clearInterval(interval);
      interval = null;
    }
  }, 1000);
}

function pauseAnim() {
  if (interval) {
    clearInterval(interval);
    interval = null;
  }
}

function resetAnim() {
  pauseAnim();
  ants = {};
  stepIndex = 0;
  document.getElementById("stepCounter").textContent = `Step: 0`;
  draw();
}

window.addEventListener('DOMContentLoaded', () => {
  // кнопки
  const startBtn = document.getElementById('startBtn');
  const pauseBtn = document.getElementById('pauseBtn');
  const resetBtn = document.getElementById('resetBtn');
  if (startBtn) startBtn.addEventListener('click', () => { console.log('Start clicked'); startAnim(); });
  if (pauseBtn) pauseBtn.addEventListener('click', () => { console.log('Pause clicked'); pauseAnim(); });
  if (resetBtn) resetBtn.addEventListener('click', () => { console.log('Reset clicked'); resetAnim(); });

  // загрузка данных
  const params = new URLSearchParams(window.location.search);
  const file = params.get('file');
  const url = file ? `/data?file=${encodeURIComponent(file)}` : '/data';

  const stepCounter = document.getElementById('stepCounter');
  if (stepCounter) stepCounter.textContent = 'Loading...';
  const startBtn2 = document.getElementById('startBtn');
  if (startBtn2) startBtn2.disabled = true;

  fetch(url)
    .then((res) => {
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      return res.json();
    })
    .then((json) => {
      data = json;
      startRoom = (data.rooms.find(r => r.isStart) || {}).name || "0";
      computePositions();
      if (stepCounter) stepCounter.textContent = 'Step: 0';
      if (startBtn2) startBtn2.disabled = false;
      draw();
    })
    .catch((err) => {
      console.error('Failed to load /data', err);
      if (stepCounter) stepCounter.textContent = 'Load error';
      if (startBtn2) startBtn2.disabled = true;
    });
});

window.addEventListener('resize', () => {
  width = window.innerWidth;
  height = window.innerHeight;
  canvas.width = width;
  canvas.height = height;
  if (data) {
    computePositions();
    draw();
  }
});

