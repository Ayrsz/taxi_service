<script lang="ts">
  export let message: string;
  export let type: 'success' | 'info' | 'error';
  export let isVisible: boolean;
</script>

{#if isVisible}
  <div class={`notification notification-${type}`}>
    <p>{message}</p>
  </div>
{/if}

<style>
  .notification {
    padding: 10px;
    border-radius: 5px;
    margin-bottom: 15px;
    color: white;
  }
  .notification-success { background-color: #4CAF50; }
  .notification-info { background-color: #2196F3; }
  .notification-error { background-color: #f44336; }
</style>