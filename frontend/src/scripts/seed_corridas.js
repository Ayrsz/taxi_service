// Cria N corridas em /api/corridas
// Uso:
//   npm run seed:corridas -- 8
//   npm run seed:corridas -- 12 --rate
//   npm run seed:corridas -- 20 --random --center="-8.05,-34.95" --radius=0.08

const API = process.env.API_URL || 'http://localhost:3000/api';

const args = process.argv.slice(2);
const count = Number(args[0] || '5');
const shouldRate = args.includes('--rate');
const useRandom = args.includes('--random');

function flag(name, def) {
  const f = args.find(a => a.startsWith(`${name}=`));
  return f ? f.split('=').slice(1).join('=') : def;
}

const [defLat, defLng] = (-8.05224 + ',' + -34.92861).split(',').map(Number);
const [centerLat, centerLng] = (flag('--center', `${defLat},${defLng}`)).split(',').map(Number);
const radius = Number(flag('--radius', '0.06'));

const hotspots = [
  { name: 'Centro de Informática', coord: '-8.050, -34.951' },
  { name: 'Recife Antigo',         coord: '-8.063, -34.871' },
  { name: 'RioMar Shopping',       coord: '-8.085, -34.893' },
  { name: 'UFPE',                  coord: '-8.050, -34.951' },
  { name: 'Avenida Boa Viagem',    coord: '-8.129, -34.900' },
  { name: 'Derby',                 coord: '-8.054, -34.898' },
  { name: 'Casa Forte',            coord: '-8.028, -34.918' },
  { name: 'Aeroporto',             coord: '-8.129, -34.918' }
];

function pickTwo() {
  let i = Math.floor(Math.random() * hotspots.length);
  let j; do { j = Math.floor(Math.random() * hotspots.length); } while (j === i);
  return { origem: hotspots[i], destino: hotspots[j] };
}
function jitter(base, d) { return base + (Math.random() * 2 - 1) * d; }
function randomPair() {
  return {
    origem:  { name: 'Origem',  coord: `${jitter(centerLat, radius).toFixed(5)}, ${jitter(centerLng, radius).toFixed(5)}` },
    destino: { name: 'Destino', coord: `${jitter(centerLat, radius).toFixed(5)}, ${jitter(centerLng, radius).toFixed(5)}` }
  };
}

async function createRide(passageiroID, origem, destino) {
  const res = await fetch(`${API}/corridas`, {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify({ passageiroID, origem, destino })
  });
  if (!res.ok) throw new Error(`HTTP ${res.status}: ${await res.text()}`);
  return res.json();
}
async function rateRide(id, nota) {
  const res = await fetch(`${API}/corridas/${id}/avaliar`, {
    method: 'POST',
    headers: { 'content-type': 'application/json' },
    body: JSON.stringify({ nota })
  });
  if (!res.ok) throw new Error(`HTTP ${res.status}: ${await res.text()}`);
}

(async () => {
  console.log(`[seed] API = ${API}`);
  for (let i = 0; i < count; i++) {
    try {
      const pair = useRandom ? randomPair() : pickTwo();
      const data = await createRide(1, pair.origem.coord, pair.destino.coord);
      const id = data.ID ?? data.id;
      process.stdout.write(`+ corrida ${id}: ${pair.origem.name} -> ${pair.destino.name}`);
      if (shouldRate) {
        const nota = 1 + Math.floor(Math.random() * 5);
        await rateRide(id, nota);
        process.stdout.write(`  | ${nota}★`);
      }
      process.stdout.write('\n');
      await new Promise(r => setTimeout(r, 120));
    } catch (e) {
      console.error('  x erro:', e.message || e);
    }
  }
  console.log('✓ pronto! Abra /historico.');
})();
