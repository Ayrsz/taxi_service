Feature: Notificação de Corrida
    Como um sistema de taxi
    Eu quero gerenciar notificações relacionadas às corridas
    Para manter motoristas informados sobre o status das corridas

Scenario: Motorista recebe notificação de chegada ao destino
    Given que estou logado no sistema como motorista
    And aceitei uma corrida que está em andamento
    When a minha localização atual coincide com o destino da corrida
    Then vejo a notificação "Você chegou ao destino"

Scenario: Sistema exibe notificação ao motorista ao solicitar cancelamento de corrida
    Given que estou logado no sistema como o motorista João
    And aceitei uma corrida com 8 km de distância estimada até o local de embarque
    And estou visualizando o tempo estimado de chegada de 10 minutos
    When seleciona a opção "Cancelar corrida"
    Then o sistema exibe uma notificação na tela com a mensagem:
    """
    Tem certeza que deseja cancelar a corrida? Cancelamentos frequentes podem impactar sua avaliação.
    """
    And essa notificação aparece com as opções "Sim, quero cancelar" e "Não, continuar com a corrida"

Scenario: Corrida finalizada com sucesso é registrada no histórico
    Given que finalizei uma corrida após chegar ao destino
    When acesso meu histórico de corridas
    Then vejo a corrida com data e hora
    And o valor total e a distância percorrida são exibidos

Scenario: Sistema exibe corrida cancelada com informações no histórico do motorista
    Given que estou logado como o motorista João
    And cancelei uma corrida que havia aceitado com destino de 12 km e valor estimado de R$ 28,00
    When acesso a tela de histórico de corridas
    Then o sistema exibe uma entrada no histórico com o status "Cancelada"
    And essa entrada apresenta a data e hora do cancelamento, distância estimada e valor previsto da corrida

Scenario: Notificação de Tempo estimado de chegada é exibido ao aceitar uma corrida
    Given que aceitei uma corrida disponível
    When a corrida for carregada na tela
    Then vejo o tempo estimado de chegada até o local de embarque