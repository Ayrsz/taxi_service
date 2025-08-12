<script lang="ts">
  import RideHistoryComponent from '../../components/RideHistoryComponent.svelte';
  import NotificationComponent from '../../components/NotificationComponent.svelte';
  import { onMount } from 'svelte';
  
  let showNotification = false;
  let notificationMessage = '';
  
  onMount(() => {
    showNotification = true;
    notificationMessage = "Bem-vindo ao seu painel!";
    setTimeout(() => {
      showNotification = false;
    }, 5000);
  });
</script>

<main>
  <h1>Painel do Motorista</h1>
  <NotificationComponent message={notificationMessage} type="info" isVisible={showNotification} />
  <RideHistoryComponent />
</main>