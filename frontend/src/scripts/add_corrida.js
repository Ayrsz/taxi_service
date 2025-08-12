// Uso: npm run add:corrida -- <passageiroID> "<origem lat, lng>" "<destino lat, lng>"
// Ex.: npm run add:corrida -- 1 "-23.55, -46.63" "-23.56, -46.62"

async function main() {
  const [, , passageiroIDArg = '1', origem = '-23.55, -46.63', destino = '-23.56, -46.62'] = process.argv;
  const passageiroID = Number(passageiroIDArg);

  const payload = { passageiroID, origem, destino };

  try {
    const res = await fetch('http://localhost:3000/api/corridas', {
      method: 'POST',
      headers: { 'content-type': 'application/json' },
      body: JSON.stringify(payload),
    });

    if (!res.ok) {
      const errText = await res.text();
      throw new Error(`HTTP ${res.status}: ${errText}`);
    }

    const data = await res.json();
    console.log('Corrida criada com sucesso:', data);
    console.log('Abra /historico no frontend para ver a nova corrida.');
  } catch (err) {
    console.error('Falha ao criar corrida:', err.message || err);
    process.exit(1);
  }
}

main();
