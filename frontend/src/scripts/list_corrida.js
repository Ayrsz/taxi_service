const API = process.env.API_URL || 'http://localhost:3000/api';
(async () => {
  const res = await fetch(`${API}/corridas`);
  if (!res.ok) throw new Error(`HTTP ${res.status}: ${await res.text()}`);
  const data = await res.json();
  console.log(data);
})();
