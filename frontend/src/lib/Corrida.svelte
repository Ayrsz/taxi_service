<script>
  import { onMount } from 'svelte';
  import { navigate } from 'svelte-routing';
  import axios from 'axios';
  import Modal from './ConfirmacaoFimCorrida.svelte';

  export let id;

  let mapElement;
  let map;
  let ride = null;
  let rideStatus = 'Carregando informações...';
  let motoristaMarker;

  // Variáveis para controlar CADA modal individualmente
  let showCancelModal = false;
  let showFinishModal = false;

  const api = axios.create({
    baseURL: 'http://localhost:3000/api',
  });
  
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
      rideStatus = formatStatus(ride.status);
      console.log("Status da corrida: ", rideStatus);
      if (ride.MotoristaLat && ride.MotoristaLng) {
        const latLng = [ride.MotoristaLat, ride.MotoristaLng];
        if (!motoristaMarker) {
          motoristaMarker = L.marker(latLng).addTo(map).bindPopup('Motorista');
          map.setView(latLng, 15);
        } else {
          motoristaMarker.setLatLng(latLng);
        }
      }

      if (ride.status.startsWith('concluída') || ride.status.startsWith('cancelada')) {
        setTimeout(() => navigate('/'), 3000);
      }
    } catch (error) {
      console.error('Erro ao buscar dados da corrida:', error);
      rideStatus = 'Erro ao carregar dados.';
    }
  }

  async function executeCancel() {
    showCancelModal = false;
    try {
      const motoristaIdParaCancelar = "1";
      await api.post(`/corrida/${id}/cancelar/motorista`, {
        motorista_id: motoristaIdParaCancelar 
      });
      alert('Sua corrida foi cancelada.');
      navigate('/');
    } catch (error) {
      console.error('Erro ao cancelar a corrida:', error);
      const errorMessage = error.response?.data?.error || 'Não foi possível cancelar a corrida.';
      alert(errorMessage);
    }
  }
  
  async function executeFinishRide() {
    showFinishModal = false; // Esconde o modal de finalização
    try {
      await api.post(`/corrida/${id}/finalizar`);
      alert('Corrida finalizada com sucesso!');
      navigate('/');
    } catch (error) {
      console.error('Erro ao finalizar a corrida:', error);
      alert('Não foi possível finalizar a corrida.');
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
    message="Esta ação não pode ser desfeita. Você tem certeza que deseja cancelar sua corrida?"
    confirmLabel="Confirmar Cancelamento"
    on:confirm={executeCancel}
    on:close={() => showCancelModal = false}
  />
{/if}

{#if showFinishModal}
  <Modal 
    title="Finalizar Corrida"
    message="Você confirma que chegou ao seu destino e deseja finalizar a corrida?"
    confirmLabel="Confirmar Finalização"
    on:confirm={executeFinishRide}
    on:close={() => showFinishModal = false}
  />
{/if}

<div class="container">
  <h1>Sua Corrida (ID: {id})</h1>

  <div id="map" bind:this={mapElement}></div>

  <div class="status">Status: {rideStatus}</div>

  <div class="actions">
    <button class="cancel-btn" on:click={() => showCancelModal = true}>Cancelar Corrida</button>
    <button class="finish-btn" on:click={() => showFinishModal = true}>Finalizar Corrida</button>
  </div>
</div>