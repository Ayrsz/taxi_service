<script>
  import { onMount } from 'svelte';
  import { navigate } from 'svelte-routing';
  import axios from 'axios';
  import Modal from './CancelCorrida.svelte'; // 1. Importe o novo componente

  export let id;

  let mapElement;
  let map;
  let ride = null;
  let rideStatus = 'Carregando informações...';
  let motoristaMarker;
  let showCancelModal = false; // 2. Crie uma variável para controlar a visibilidade do modal

  const api = axios.create({
    baseURL: 'http://localhost:3000/api',
  });


  console.log("O valor do ID recebido é:", id); // Adicione esta linha para depurar
  

  onMount(() => {
    map = L.map(mapElement).setView([-23.55052, -46.633308], 14);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    fetchRideData();
    const interval = setInterval(fetchRideData, 3000);
    return () => clearInterval(interval);
  });

  async function fetchRideData() {
    try {
      const response = await api.get(`/corrida/${id}`);
      ride = response.data;
      rideStatus = formatStatus(ride.Status);

      if (ride.MotoristaLat && ride.MotoristaLng) {
        const latLng = [ride.MotoristaLat, ride.MotoristaLng];
        if (!motoristaMarker) {
          motoristaMarker = L.marker(latLng).addTo(map).bindPopup('Motorista');
          map.setView(latLng, 15);
        } else {
          motoristaMarker.setLatLng(latLng);
        }
      }

      if (ride.Status.startsWith('concluída') || ride.Status.startsWith('cancelada')) {
        setTimeout(() => navigate('/'), 3000);
      }
    } catch (error) {
      console.error('Erro ao buscar dados da corrida:', error);
      rideStatus = 'Erro ao carregar dados.';
    }
  }

  // 3. Esta função agora APENAS faz a chamada para a API
  async function executeCancel() {
  showCancelModal = false; // Esconde o modal
  try {
    // ID do motorista que está a cancelar (para este exemplo, usamos um valor fixo)
    const motoristaIdParaCancelar = "1"; // Corresponde ao motorista_joao dos seus testes

    // 1. Endpoint corrigido para /cancelar/motorista
    // 2. Adicionado o corpo (payload) com o motorista_id
    await api.post(`/corrida/${id}/cancelar/motorista`, {
      motorista_id: motoristaIdParaCancelar 
    });

    alert('Sua corrida foi cancelada.');
    navigate('/');
  } catch (error) {
    console.error('Erro ao cancelar a corrida:', error);
    // Extrai a mensagem de erro do backend, se disponível
    const errorMessage = error.response?.data?.error || 'Não foi possível cancelar a corrida.';
    alert(errorMessage);
  }
}
  
  // Função para finalizar a corrida (não foi alterada)
  async function finishRide() {
    if (confirm('Confirmar a finalização da corrida?')) {
      try {
        await api.post(`/corrida/${id}/finalizar`);
        alert('Corrida finalizada com sucesso!');
        navigate('/');
      } catch (error) {
        console.error('Erro ao finalizar a corrida:', error);
        alert('Não foi possível finalizar a corrida.');
      }
    }
  }

  function formatStatus(status) {
    return status.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase());
  }
</script>

<style>
  .container {
    padding: 2rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  #map {
    height: 400px;
    width: 100%;
    background-color: #eee;
  }
  .status {
    font-size: 1.2rem;
    font-weight: bold;
    text-align: center;
    padding: 1rem;
    background-color: #f0f0f0;
    border-radius: 8px;
  }
  .actions {
      display: flex;
      gap: 1rem;
  }
  button {
    padding: 0.75rem;
    font-size: 1rem;
    cursor: pointer;
    flex: 1;
    border: none;
    border-radius: 5px;
    color: white;
  }
  .cancel-btn {
      background-color: #f44336;
  }
  .finish-btn {
      background-color: #4CAF50;
  }
</style>

{#if showCancelModal}
  <Modal 
    title="Cancelar Corrida"
    message="Esta ação não pode ser desfeita. Você tem certeza que deseja cancelar sua corrida? Cancelamentos frequentes podem resultar em penalidades."
    on:confirm={executeCancel}
    on:close={() => showCancelModal = false}
  />
{/if}

<div class="container">
  <h1>Sua Corrida (ID: {id})</h1>

  <div id="map" bind:this={mapElement}></div>

  <div class="status">Status: {rideStatus}</div>

  <div class="actions">
    <button class="cancel-btn" on:click={() => showCancelModal = true}>Cancelar Corrida</button>
    <button class="finish-btn" on:click={finishRide}>Finalizar Corrida</button>
  </div>
</div>