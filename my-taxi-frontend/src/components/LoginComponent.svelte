<script lang="ts">
  let email = '';
  let password = '';
  let errorMessage = '';

  const handleLogin = async () => {
    errorMessage = '';
    try {
      // Endpoint do seu backend em Go
      const response = await fetch('http://localhost:3000/api/v1/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        throw new Error('Credenciais inválidas.');
      }

      const data = await response.json();
      localStorage.setItem('authToken', data.token);
      
      window.location.href = '/dashboard';
    } catch (error) {
      errorMessage = error.message;
    }
  };
</script>

<div class="login-container">
  <h2>Login do Motorista</h2>
  <form on:submit|preventDefault={handleLogin}>
    {#if errorMessage}
      <p class="error-message">{errorMessage}</p>
    {/if}
    <label for="email">Email:</label>
    <input id="email" type="email" bind:value={email} required />
    
    <label for="password">Senha:</label>
    <input id="password" type="password" bind:value={password} required />
    
    <button type="submit">Entrar</button>
  </form>
</div>