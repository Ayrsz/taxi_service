Feature: Notificação de novas corridas por perto
  Como um motorista
  Eu quero receber notificações de novas corridas próximas
  Para que eu possa aceitar corridas 

  Background:
    Given que o motorista de nome “João” está logado
    And sua localização está em "Rua A, 123"
    And está com o status "disponível"

  Scenario: Receber notificação de corrida próxima quando estiver disponível
    Given uma corrida é solicitada por "Maria Silva" 
    When estou num raio menor que 2 km da corrida
    Then deve receber uma notificação contendo:
      | tempo       | 23 minutos           |
      | distância   | 1.5 km               |
      | valor       | R$ 12,50             |
      | passageiro  | Maria Silva          |
    And deve ter botões de "Aceitar" e "Recusar" 

  Scenario: Notificação expira após tempo limite
    Given recebe uma notificação de corrida
    When passa 20 segundos sem interagir com a notificação
    Then a notificação deve desaparecer automaticamente
    And o status deve permanecer como "disponível"
    And outro motorista pode aceitar a corrida

  Scenario: Corrida aceita é registrada no histórico
    Given recebe uma notificação de corrida
    When clica em "Aceitar"
    Then o status deve ser alterado para "ocupado"
    And a corrida deve ser registrada no histórico como "em_andamento"

  Scenario: Motorista ocupado não recebe notificações
    Given o status está "ocupado"
    When uma corrida é solicitada por "Carlos Pereira"
    And estou num raio menor que 2 km da corrida
    Then não deve receber nenhuma notificação

  Scenario: Recusar corrida mantém disponibilidade
    Given recebe uma notificação de corrida
    When clica em "Recusar"
    Then o status deve permanecer "disponível"
    And outro motorista pode aceitar a corrida

  Scenario: Desativar notificações altera status
    Given está recebendo notificações
    When desativa a opção "Receber notificações de novas corridas"
    Then o status deve ser alterado para "ocupado"

  Scenario: Ativar notificações altera status
    Given as notificações estão desativadas
    And o status está "ocupado"
    When ativa a opção "Receber notificações de novas corridas"
    Then o status deve ser alterado para "disponível"
    And deve receber notificações de novas corridas

  Scenario: Motorista não recebe notificações fora do raio
    Given uma corrida é solicitada por "Ana Costa"
    When estou num raio maior que 2 km da corrida
    Then não deve receber nenhuma notificação
