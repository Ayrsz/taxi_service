Feature: Notificação de novas corridas por perto

  Scenario: Receber notificação de corrida próxima quando estiver disponível
    Given o motorista de nome “João” está logado
    And está com o status "disponível"
    When uma nova corrida é solicitada por um passageiro próximo
    Then deve receber uma notificação com os detalhes da corrida
    And deve ter a opção de aceitar ou recusar a corrida

  Scenario: Notificação expira após tempo limite
    Given que o motorista de nome “João” está logado
    And está com o status "disponível"
    And recebe uma notificação de corrida
    When não interage com a notificação em até “20” segundos
    Then a notificação deve desaparecer

  Scenario: Corrida aceita deve ser cadastrada no histórico e impedir novas notificações
    Given o motorista de nome “João” está logado
    And está com o status "disponível"
    And recebi uma notificação de uma nova corrida próxima
    When aceita a corrida
    Then o status deve ser alterado para "ocupado"
    And a corrida deve ser registrada no histórico como "aceita"
    And enquanto o status de “João” estiver como "ocupado", não deve receber notificações de novas corridas

  Scenario: Motorista recusa a corrida
    Given o motorista de nome “João” está logado
    And está com o status "disponível"
    And recebe uma notificação de uma nova corrida próxima
    When recusa a corrida
    Then o status deve permanecer como "disponível"
    And não deve ser registrada nenhuma corrida no histórico
    And deve continuar recebendo notificações de novas corridas próximas

  Scenario: Motorista desativa notificações de novas corridas
    Given o motorista de nome “João” está logado
    When eu desativo a opção "Receber notificações de novas corridas"
    Then o status deve estar “ocupado”
    And não deve receber notificações de novas corridas

 Scenario: Motorista ativa notificações de novas corridas
    Given o motorista de nome “João” está logado
    When ativa a opção "Receber notificações de novas corridas"
    Then o status deve estar “disponível”
    And deve voltar a receber notificações de novas corridas
