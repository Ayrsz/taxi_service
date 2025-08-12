<script>
  import { Router, Route } from "svelte-routing";
  import Home from "./lib/Home.svelte";
  import Corrida from "./lib/Corrida.svelte";
  import Historico from "./lib/HistoricoEAvaliacao/Historico.svelte";
</script>

<Router>
  <Route path="/" component={Home} />
  <Route path="/corrida/:id" component={Corrida} />
  <Route path="/" component={Home} exact />
  <Route path="/historico" component={Historico} />
</Router>
