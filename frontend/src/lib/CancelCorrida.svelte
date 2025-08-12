<script>
  import { createEventDispatcher } from 'svelte';

  // Props para customizar o modal
  export let title = 'Confirmar Ação';
  export let message = 'Você tem certeza?';

  const dispatch = createEventDispatcher();

  // Funções que disparam eventos para o componente pai
  function handleConfirm() {
    dispatch('confirm');
  }

  function handleClose() {
    dispatch('close');
  }
</script>

<style>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.6);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .modal-box {
    background-color: white;
    padding: 2rem;
    border-radius: 12px;
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
    width: 90%;
    max-width: 400px;
    text-align: center;
  }

  h2 {
    margin-top: 0;
    color: #333;
  }

  p {
    color: #555;
    margin-bottom: 2rem;
  }

  .modal-actions {
    display: flex;
    gap: 1rem;
  }

  button {
    flex: 1;
    padding: 0.75rem;
    font-size: 1rem;
    border-radius: 8px;
    border: none;
    cursor: pointer;
    font-weight: bold;
    transition: transform 0.1s ease;
  }

  button:active {
      transform: scale(0.98);
  }

  .confirm-btn {
    background-color: #f44336; /* Vermelho para ação perigosa */
    color: white;
  }

  .close-btn {
    background-color: #e0e0e0;
    color: #333;
  }
</style>

<div class="modal-overlay" on:click={handleClose}>
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="modal-box" on:click|stopPropagation>
    <h2>{title}</h2>
    <p>{message}</p>
    <div class="modal-actions">
      <button class="close-btn" on:click={handleClose}>Voltar</button>
      <button class="confirm-btn" on:click={handleConfirm}>Confirmar Cancelamento</button>
    </div>
  </div>
</div>