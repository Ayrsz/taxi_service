<script>
  import { onMount } from 'svelte';
  import { navigate } from 'svelte-routing';
  import axios from 'axios';

  let mapElement;
  let map;
  let originMarker, destinationMarker;
  let originCoords, destinationCoords;

  onMount(() => {
    map = L.map(mapElement).setView([-23.55052, -46.633308], 12);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    map.on('click', (e) => {
      if (!originMarker) {
        originMarker = L.marker(e.latlng).addTo(map).bindPopup('<b>Origem</b>').openPopup();
        originCoords = e.latlng;
      } else if (!destinationMarker) {
        destinationMarker = L.marker(e.latlng).addTo(map).bindPopup('<b>Destino</b>').openPopup();
        destinationCoords = e.latlng;
      }
    });
  });

  async function startRide() {
    if (!originCoords || !destinationCoords) {
      alert('Por favor, selecione os pontos de partida e destino no mapa.');
      return;
    }

    try {
      // Define um ID de passageiro fixo para a demonstração
      const passageiroId = 1;

      const response = await axios.post('http://localhost:3000/api/corrida', {
        passageiroId: passageiroId, 
        origem: `${originCoords.lat}, ${originCoords.lng}`,
        destino: `${destinationCoords.lat}, ${destinationCoords.lng}`,
        // O backend irá calcular o tempo estimado, então não precisamos enviar
      });
      
      const corrida = response.data;
      console.log('Corrida criada com sucesso! ID da Corrida:', corrida.id);
      alert(`Corrida criada com sucesso! O ID da sua corrida é: ${corrida.id}`);

      // Navega para a página da corrida específica
      navigate(`/corrida/${corrida.id}`);

    } catch (error) {
      console.error('Erro ao iniciar a corrida:', error);
      alert('Não foi possível iniciar a corrida. Tente novamente.');
    }
  }

  function clearSelection() {
      if(originMarker) map.removeLayer(originMarker);
      if(destinationMarker) map.removeLayer(destinationMarker);
      originMarker = null;
      destinationMarker = null;
      originCoords = null;
      destinationCoords = null;
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
    height: 500px;
    width: 100%;
    background-color: #eee;
    cursor: crosshair;
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
  }
  .clear-btn {
      background-color: #f44336;
      color: white;
  }
</style>

<div class="container">
  <h1>Solicitar uma Corrida</h1>
  
  <p>Clique no mapa para definir o local de <b>partida</b> e depois o de <b>destino</b>.</p>

  <div id="map" bind:this={mapElement}></div>

  <div class="actions">
    <button class="clear-btn" on:click={clearSelection}>Limpar Seleção</button>
    <button on:click={startRide}>Chamar Taxi</button>
  </div>
</div>
