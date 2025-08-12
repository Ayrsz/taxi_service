Feature: Cancelamento de Corrida
  As a sistema de monitoramento de corridas
  I want to permitir e gerenciar o cancelamento de corridas
  So that apenas corridas válidas podem ser canceladas

  Contexto:
    Given que o sistema possui os seguintes motoristas
      | nome_referencia | nome_completo | cpf              | status_inicial |
      | motorista_joao  | "João Silva"  | "123.456.789-00" | "disponível"   |

    # A tabela de corridas foi atualizada para usar os status oficiais do modelo
    # e incluir os campos de embarque/desembarque.
    And que as seguintes corridas existem no sistema, alocadas para "motorista_joao":
      | id  | status_inicial | local_embarque       | local_desembarque   | tempo_estimado_minutos |
      | 101 | "pendente"     | "Rua das Flores, 10" | "Av. Principal, 20" | 5                      |
      | 102 | "andamento"    | "Rua da Praia, 30"   | "Centro Cívico, 40" | 20                     |
      | 103 | "finalizada"   | "Praça da Sé, 50"    | "Aeroporto, 60"     | 15                     |

  Scenario: Motorista cancela uma corrida pendente
    Given a corrida de id "101" tem o status "pendente"
    When o motorista de cpf "123.456.789-00" seleciona a opção "Cancelar" no aplicativo para a corrida "101"
    Then o sistema atualiza o status da corrida "101" para "cancelada"
    And o motorista recebe a mensagem de confirmação na tela "Corrida cancelada com sucesso"
    And o status do motorista de cpf "123.456.789-00" é atualizado para "disponível"

  Scenario: Sistema cancela corrida por demora na partida do motorista
    Given a corrida de id "101" está no status "pendente" há mais de "10" minutos
    When o sistema de monitoramento executa a verificação de corridas pendentes
    Then o sistema altera o status da corrida "101" para "cancelada"
    And o sistema envia uma notificação para o motorista de cpf "123.456.789-00" com a mensagem "Corrida cancelada por excesso de tempo de partida"

  Scenario: Tentativa de cancelamento de corrida em andamento pelo motorista
    Given a corrida de id "102" tem o status "andamento"
    When o motorista de cpf "123.456.789-00" tenta cancelar a corrida de id "102" através do aplicativo
    Then o sistema rejeita a solicitação de cancelamento
    And o motorista recebe uma mensagem de erro na tela "Não é possível cancelar uma corrida que já foi iniciada"
    And o status da corrida "102" permanece "andamento"

  Scenario: Tentativa de cancelamento de corrida já finalizada pelo motorista
    Given a corrida de id "103" tem o status "finalizada"
    When o motorista de cpf "123.456.789-00" tenta cancelar a corrida de id "103" através do aplicativo
    Then o sistema rejeita a solicitação de cancelamento
    And o motorista recebe uma mensagem de erro na tela "Não é possível cancelar uma corrida finalizada"
    And o status da corrida "103" permanece "finalizada"