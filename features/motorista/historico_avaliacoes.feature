Feature: Histórico e avaliação de corridas
As a usuário
I want to acessar o histórico de corridas
So that eu possa avaliar as corridas e/ou visualizar as corridas que fiz

Scenario: Exibição do histórico de corridas para o motorista
Given que estou logado como motorista,
And estou na tela de menu principal,
When clico na opção "Histórico de Corridas",
Then vejo uma lista contendo as últimas corridas realizadas,
And cada item da lista exibe destino, data, hora e valor da corrida.

Scenario: Acessar interface de avaliação de corrida 
Given que estou logado como motorista,
And tenho pelo menos uma corrida realizada que ainda não foi avaliada,
When acesso histórico de corridas,
And clico em “Avaliar corrida” ao lado da corrida pendente,
Then sou redirecionado para a tela de avaliação de corrida,
And vejo os campos de nota e comentário disponíveis para preenchimento.

Scenario: Avaliar corrida realizada pelo motorista
Given que estou logado como motorista,
And estou na tela de avaliação de corridas,
And selecionei uma corrida ainda não avaliada,
When escolho uma nota de 1 a 5 estrelas e adiciono um comentário opcional,
And clico em "Enviar avaliação",
Then vejo uma mensagem "Avaliação registrada com sucesso",
And a corrida é marcada como avaliada no histórico.

Scenario: Tentativa de avaliar uma corrida já avaliada 
Given que estou logado como motorista,
And estou na tela de histórico de corridas,
And existe uma corrida que já foi avaliada anteriormente,
When visualizo os detalhes dessa corrida,
Then não vejo o botão "Avaliar",
And os dados da avaliação anterior são exibidos

Scenario: Visualização de corrida com avaliação pendente
Given que estou logado como motorista,
And tenho uma corrida anterior ainda não avaliada,
When estou na tela de menu principal,
Then vejo um pop-up escrito “Como foi a sua corrida?”,
And o pop-up exibe os botões “Enviar avaliação” e  “pular”,
And eu clico “pular”, 
Then o pop-up é fechado.

