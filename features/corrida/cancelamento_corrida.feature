Feature: Tempo de Corrida
As a sistema de monitoramento de corridas
I want to acompanhar o tempo de execução de uma corrida
So that eu possa reagir de forma adequada a atrasos ou adiantamentos

Scenario: Notificar atraso do motorista após tempo limite
Given o tempo estimado de chegada ao destino é de "20" minutos
And o tempo decorrido é de "25" minutos
When o sistema detecta que o tempo foi ultrapassado
Then uma notificação é enviada para o passageiro
And o status do motorista é atualizado para "atrasado"
And é exibida a tag "atrasado" na tela do motorista

Scenario: Motorista conclui corrida antes do tempo estimado
Given o tempo estimado de chegada ao destino é de "20" minutos
And o tempo decorrido é de "15" minutos
When o motorista chega ao destino
Then o status da corrida é "concluída com antecedência"
And uma mensagem de parabéns é exibida ao motorista
And um bônus é aplicado ao motorista

Scenario: Motorista conclui corrida dentro do tempo estimado
Given o tempo estimado de chegada ao destino é de "20" minutos
And o tempo decorrido é de "19" minutos
When o motorista chega ao destino
Then o status da corrida é "concluída no tempo previsto"
And nenhuma penalização ou bônus é aplicado

Scenario: Corrida cancelada automaticamente por excesso de tempo
Given o tempo estimado de chegada é de "20" minutos
And o tempo decorrido é de "36" minutos
When o sistema detecta que o tempo ultrapassou o limite de "15" minutos extras
Then a corrida é automaticamente cancelada
And o motorista recebe uma notificação por pop-up com a mensagem "Corrida cancelada por excesso de tempo"