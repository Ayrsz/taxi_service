<!-- src/lib/HistoricoEAvaliacao/Historico.svelte -->
<script>
  import { onMount } from 'svelte';
  import axios from 'axios';
  import { navigate } from 'svelte-routing';

  let carregando = true;
  let erro = '';
  let corridas = [];
  let editandoIdx = null;
  let notaSelecionada = 0;

  // Lugares conhecidos (para trocar coordenadas por nomes)
  const HOTSPOTS = [
    { name: 'Centro de Informática', lat: -8.050, lng: -34.951 },
    { name: 'Recife Antigo',         lat: -8.063, lng: -34.871 },
    { name: 'RioMar Shopping',       lat: -8.085, lng: -34.893 },
    { name: 'UFPE',                  lat: -8.050, lng: -34.951 },
    { name: 'Avenida Boa Viagem',    lat: -8.129, lng: -34.900 },
    { name: 'Derby',                 lat: -8.054, lng: -34.898 },
    { name: 'Casa Forte',            lat: -8.028, lng: -34.918 },
    { name: 'Aeroporto',             lat: -8.129, lng: -34.918 }
  ];
  const COORD_RE = /^\s*-?\d+(\.\d+)?,\s*-?\d+(\.\d+)?\s*$/;

  function parseCoord(str) {
    if (!str || !COORD_RE.test(str)) return null;
    const [a, b] = str.split(',').map(s => Number(s.trim()));
    return { lat: a, lng: b };
  }
  function nearly(a, b, eps = 0.003) { // ~300 m
    return Math.abs(a - b) <= eps;
  }
  function coordToName(str) {
    const p = parseCoord(str);
    if (!p) return str; // já é nome/endereço
    const hit = HOTSPOTS.find(h => nearly(h.lat, p.lat) && nearly(h.lng, p.lng));
    return hit ? hit.name : str;
  }

  // Haversine (km) – para estimar preço quando não vier do backend
  function distKm(a, b) {
    if (!a || !b) return 0;
    const R = 6371;
    const dLat = (b.lat - a.lat) * Math.PI / 180;
    const dLng = (b.lng - a.lng) * Math.PI / 180;
    const la1 = a.lat * Math.PI / 180, la2 = b.lat * Math.PI / 180;
    const h = Math.sin(dLat/2)**2 + Math.cos(la1)*Math.cos(la2)*Math.sin(dLng/2)**2;
    return 2 * R * Math.asin(Math.sqrt(h));
  }

  onMount(async () => {
    try {
      const { data } = await axios.get('http://localhost:3000/api/corridas', { timeout: 10000 });
      const itens = Array.isArray(data) ? data : (data?.corridas ?? []);

      corridas = itens.map((c) => {
        const origemRaw  = c.Origem ?? c.origem ?? '';
        const destinoRaw = c.Destino ?? c.destino ?? '';

        // Data: usa DataInicio/CreatedAt; ignora "0001-01-01..."
        const d0 = c.DataInicio ?? c.dataInicio ?? c.CreatedAt ?? c.createdAt ?? null;
        const isZeroDate = typeof d0 === 'string' && d0.startsWith('0001-01-01');
        const dataInicio = isZeroDate ? null : d0 ?? null;

        // Preço: usa do backend; se faltar, estima pela distância
        let preco = c.Preco ?? c.preco ?? c.valor ?? null;
        if (!preco || Number(preco) === 0) {
          const o = parseCoord(origemRaw);
          const d = parseCoord(destinoRaw);
          const km = distKm(o, d);
          if (km > 0) {
            // bandeirada 6 + 2.8/km (ajuste se quiser)
            preco = Math.max(0, 6 + 2.8 * km);
          } else {
            preco = null;
          }
        }

        return {
          id: c.ID ?? c.id,
          origem: coordToName(origemRaw),
          destino: coordToName(destinoRaw),
          dataInicio,
          preco,
          avaliacao: c.Avaliacao ?? c.avaliacao ?? null
        };
      });
    } catch (e) {
      erro = 'Não foi possível carregar o histórico.';
      console.error(e);
    } finally {
      carregando = false;
    }
  });

  function formatarData(iso) {
    if (!iso) return '—';
    const d = new Date(iso);
    if (isNaN(d.getTime())) return '—';
    const dia = d.toLocaleDateString('pt-BR', { month: 'short', day: '2-digit' });
    const hora = d.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' });
    return `${dia} ${hora}`;
  }
  function formatarPreco(v) {
    if (v == null) return '—';
    return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(Number(v));
  }

  function abrirAvaliacao(idx, jaAvaliada) {
   if (jaAvaliada) return;
   editandoIdx = idx;
   notaSelecionada = 0;
 }

  async function enviarAvaliacao(c, idx, nota) {
    try {
        if (c.id && c.id !== 0) {
        await axios.post(`http://localhost:3000/api/corridas/${c.id}/avaliar`, { nota });
        } else {
       console.warn('Corrida sem ID válido; atualizando somente no cliente.');
     }
     // atualiza apenas a linha clicada
      corridas = corridas.map((row, i) => i === idx ? { ...row, avaliacao: nota } : row);
      editandoIdx = null;
      notaSelecionada = 0;
    } catch (e) {
      alert('Não foi possível enviar a avaliação. Tente novamente.');
      console.error(e);
    }
  }
</script>

<style>
  .page {
    background: #d9d9d9;
    min-height: 100vh;
    padding: 16px 16px 32px;
    font-family: system-ui, -apple-system, Segoe UI, Roboto, Arial, sans-serif;
  }
  .topbar {
    display: flex;
    align-items: center;
    height: 32px;
    margin-bottom: 8px;
    font-size: 14px;
    color: #222;
  }
  .title {
    text-align: center;
    font-weight: 700;
    font-size: 28px;
    margin: 8px 0 16px;
    color: #111;
  }
  .list { display: flex; flex-direction: column; gap: 12px; }

  .item {
    display: grid;
    grid-template-columns: 64px 1fr auto;
    align-items: center;
    gap: 12px;
    background: #eaeaea;
    border-radius: 16px;
    padding: 12px;
  }
  .avatar {
    width: 64px; height: 64px;
    border-radius: 16px;
    background: #fff;
  }
  .meta { display: flex; flex-direction: column; gap: 4px; }
  .route {
    font-weight: 700; color: #111;
    overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
    max-width: 60vw;
  }
  .sub { font-size: 13px; color: #444; }

  .right {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    min-width: 120px;
  }

  .stars { font-size: 22px; letter-spacing: 2px; }
  .star-btn {
    background: none; border: none; cursor: pointer; padding: 0;
    font-size: 22px; line-height: 1;
  }
  .avaliar {
    font-size: 13px; color: #2a2a2a; text-decoration: underline; cursor: pointer;
  }
  .muted { opacity: .5; pointer-events: none; }
  .empty { text-align: center; margin-top: 32px; color: #333; }

  @media (max-width: 420px) {
    .right { min-width: 96px; }
  }
</style>

<div class="page">
  <div class="topbar">
    <span style="cursor:pointer" on:click={() => navigate(-1)}>voltar</span>
  </div>

  <div class="title">Histórico de Corridas</div>

  {#if carregando}
    <div class="empty">Carregando…</div>
  {:else if erro}
    <div class="empty">{erro}</div>
  {:else if corridas.length === 0}
    <div class="empty">Você ainda não tem corridas.</div>
  {:else}
    <div class="list">
      {#each corridas as c,i}
        <div class="item">
          <div class="avatar" />

          <div class="meta">
            <!-- Título: ORIGEM → DESTINO -->
            <div class="route" title={`${c.origem} → ${c.destino}`}>
              {c.origem} → {c.destino}
            </div>
            <div class="sub">{formatarData(c.dataInicio)}</div>
            <div class="sub">{formatarPreco(c.preco)}</div>
          </div>

          <div class="right">
            {#if c.avaliacao != null}
              <!-- já avaliada -->
              <div class="stars" aria-label={`Nota ${c.avaliacao} de 5`}>
                {#each Array(5) as _, i}
                  {@const on = i < c.avaliacao}
                  <span>{on ? '★' : '☆'}</span>
                {/each}
              </div>
              <span class="avaliar muted">Avaliado</span>
            {:else if editandoIdx === i}
              <!-- modo de avaliação -->
              <div class="stars">
                {#each [1,2,3,4,5] as n}
                  <button class="star-btn"
                          aria-label={`Dar nota ${n}`}
                          on:click={() => { notaSelecionada = n; enviarAvaliacao(c, i, n); }}>
                    {n <= (notaSelecionada || 0) ? '★' : '☆'}
                  </button>
                {/each}
              </div>
              <span class="sub">Toque para enviar</span>
            {:else}
              <!-- não avaliada -->
              <div class="stars" aria-hidden="true">
                {#each Array(5) as _}<span>☆</span>{/each}
              </div>
              <span class="avaliar" on:click={() => abrirAvaliacao(i, false)}>Avaliar</span>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
