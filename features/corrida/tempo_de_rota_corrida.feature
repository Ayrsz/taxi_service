Feature: Tempo de Corrida

Scenario 1: O motorista extrapolou o tempo de espera
  Given que o tempo de espera superou o tempo estimado
  When o sistema detecta que o tempo limite foi ultrapassado
  Then uma notificação deve ser enviada para o usuário
  And o status do motorista deve ser atualizado para "atrasado"
  And é plotado na tela do motorista uma tag de "atrasado"

Scenario 2: O motorista chegou antes do tempo estimado
  Given que o motorista iniciou a corrida
  And chegou ao destino antes do tempo estimado
  When o sistema calcula o tempo total percorrido
  Then o status da corrida deve ser "concluída com antecedência"
  And uma mensagem de parabéns pode ser exibida ao motorista
  And um bônus pode ser aplicado

Scenario 3: O motorista chegou dentro do tempo estimado
  Given que o tempo decorrido está dentro da margem do tempo estimado
  When o motorista chega ao destino
  Then o status da corrida deve ser "concluída no tempo previsto"
  And nenhuma penalização ou bônus deve ser aplicado

Scenario 4: Corrida cancelada por excesso de tempo
  Given que o tempo estimado foi ultrapassado em mais de 15 minutos
  When o sistema detecta a demora excessiva
  Then a corrida deve ser automaticamente cancelada
  And uma notificação de cancelamento deve ser enviada ao motorista por um pop-up
