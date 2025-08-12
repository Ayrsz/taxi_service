<script>
  import { onMount } from 'svelte';
  import axios from 'axios';

  export let show = false;
  export let onClose = () => {};

  let corrida = null;
  let nota = 0;
  let enviando = false;

  // Busca dados da última corrida ao montar
  onMount(async () => {
    if (show) {
      try {
        const res = await axios.get('http://localhost:3000/api/ultima-corrida');
        corrida = res.data;
      } catch {
        corrida = null;
      }
    }
  });

  function selecionarNota(n) {
    nota = n;
  }

  async function enviarAvaliacao() {
    if (!nota || !corrida) return;
    enviando = true;
    try {
      await axios.post('http://localhost:3000/api/avaliacao', {
        corridaId: corrida.id,
        nota
      });
      onClose();
    } catch {
      // Trate erro se necessário
    } finally {
      enviando = false;
    }
  }
</script>

{#if show && corrida}
  <div class="popup-backdrop">
    <div class="popup">
      <div class="header">
        <span class="titulo">Como foi a sua corrida?</span>
        <span class="pular" on:click={onClose}>pular</span>
      </div>
      <div class="info">
        <div class="img-placeholder"></div>
        <div>
          <div class="destino"><b>{corrida.destino}</b></div>
          <div>{corrida.dataHora}</div>
          <div>R$ {corrida.valor}</div>
        </div>
      </div>
      <div class="estrelas">
        {#each Array(5) as _, i}
          <span
            class="estrela"
            on:click={() => selecionarNota(i + 1)}
            style="color: {i < nota ? '#FFD700' : '#000'}"
            >&#9733;</span>
        {/each}
      </div>
      <button class="enviar" on:click={enviarAvaliacao} disabled={!nota || enviando}>
        Enviar avaliação
      </button>
    </div>
  </div>
{/if}

<style>
.popup-backdrop {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.popup {
  background: #fff;
  border-radius: 16px;
  border: 2px solid #222;
  padding: 24px;
  min-width: 380px;
  box-shadow: 0 2px 16px rgba(0,0,0,0.08);
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.titulo {
  font-size: 1.3rem;
  font-weight: bold;
}
.pular {
  font-size: 1rem;
  color: #444;
  cursor: pointer;
}
.info {
  display: flex;
  align-items: flex-start;
  margin: 18px 0 10px 0;
  gap: 12px;
}
.img-placeholder {
  width: 48px;
  height: 48px;
  background: #eee;
  border-radius: 10px;
  margin-right: 8px;
}
.destino {
  font-weight: bold;
}
.estrelas {
  display: flex;
  justify-content: center;
  margin: 18px 0 10px 0;
  font-size: 2rem;
  gap: 8px;
}
.estrela {
  cursor: pointer;
  transition: transform 0.1s;
}
.estrela:hover {
  transform: scale(1.2);
}
.enviar {
  display: block;
  margin: 0 auto;
  background: #888;
  color: #fff;
  border: none;
  border-radius: 16px;
  padding: 8px 32px;
  font-size: 1rem;
  margin-top: 10px;
  cursor: pointer;
  opacity: 1;
  transition: opacity 0.2s;
}
.enviar:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>