

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

// Загружаем SVG изображение муравья
antImage.onload = function() {
  console.log('Ant SVG image loaded successfully');
  // Перерисовываем после загрузки изображения
  if (data) {
    draw();
  }
};
antImage.src = '/static/image/ANTS.svg';

fetch("/data")
  .then((res) => res.json())
  .then((json) => {
    data = json;
    startRoom = data.rooms.find(r => r.isStart)?.name || "0";

    const allX = data.rooms.map(r => r.x);
    const allY = data.rooms.map(r => r.y);
    const minX = Math.min(...allX), maxX = Math.max(...allX);
    const minY = Math.min(...allY), maxY = Math.max(...allY);

    const scale = Math.min((width - 200) / (maxX - minX + 1), (height - 200) / (maxY - minY + 1));
    data.rooms.forEach((room) => {
      positions[room.name] = {
        x: (room.x - minX) * scale + 100,
        y: (room.y - minY) * scale + 100,
      };
    });
    draw();
  });

function draw() {
  ctx.clearRect(0, 0, width, height);
  if (!data) return;

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

  data.rooms.forEach((room) => {
    const { x, y } = positions[room.name];
    ctx.beginPath();
    ctx.arc(x, y, 20, 0, Math.PI * 2);
    ctx.fillStyle = room.name === startRoom ? "#4caf50" : room.name === "end" ? "#f44336" : "#fff";
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

  for (const [name, info] of Object.entries(ants)) {
    const from = positions[info.from];
    const to = positions[info.to];
    const t = info.progress;
    const x = from.x + (to.x - from.x) * t;
    const y = from.y + (to.y - from.y) * t;
    
    // Рисуем SVG изображение муравья
    if (antImage.complete) {
      ctx.drawImage(antImage, x - ANT_SIZE/2, y - ANT_SIZE/2, ANT_SIZE, ANT_SIZE);
    } else {
      // Fallback: если изображение не загружено, рисуем кружок
      ctx.beginPath();
      ctx.arc(x, y, 10, 0, Math.PI * 2);
      ctx.fillStyle = info.color;
      ctx.fill();
    }
    
    // Рисуем имя муравья под изображением
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
  if (!data || interval) return;

  // 🟢 ДОБАВЛЕНО: начальное размещение всех муравьев в стартовой комнате
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
