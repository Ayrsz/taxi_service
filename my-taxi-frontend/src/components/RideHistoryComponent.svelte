<script lang="ts">
  import { onMount } from 'svelte';
  
  // Exemplo de como você pode definir os tipos. Você pode criar um arquivo separado, como src/types/models.ts
  interface Ride {
    id: number;
    status: string;
    ActualDistanceKM: number;
    ActualValue: number;
    // ... outros campos
  }

  let rides: Ride[] = [];
  let isLoading = true;

  onMount(async () => {
    // ID do motorista (em um app real, viria do login)
    const driverId = 1; 
    const token = localStorage.getItem('authToken');
    
    if (!token) {
      // Redirecionar para login
      return;
    }
    
    try {
      // Endpoint do seu backend em Go
      const response = await fetch(`http://localhost:8080/api/v1/drivers/${driverId}/rides/history`, {
        headers: { 'Authorization': `Bearer ${token}` },
      });

      if (!response.ok) {
        throw new Error('Falha ao buscar histórico.');
      }
      
      const data = await response.json();
      rides = data.rides;
    } catch (error) {
      console.error(error);
    } finally {
      isLoading = false;
    }
  });
</script>

<div class="history-container">
  <h2>Histórico de Corridas</h2>
  {#if isLoading}
    <p>Carregando histórico...</p>
  {:else if rides.length === 0}
    <p>Nenhuma corrida no histórico.</p>
  {:else}
    <ul>
      {#each rides as ride (ride.id)}
        <li class="ride-item">
          <p>Corrida #{ride.id} - Status: <span>{ride.status}</span></p>
          <p>Distância: {ride.ActualDistanceKM} km</p>
          <p>Valor: R$ {ride.ActualValue}</p>
        </li>
      {/each}
    </ul>
  {/if}
</div>