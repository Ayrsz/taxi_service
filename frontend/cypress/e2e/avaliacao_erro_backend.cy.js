// cypress/e2e/avaliacao_erro_backend.cy.js
// Objetivo: Se o backend falhar ao avaliar, a UI mantém o modo de edição e alerta o usuário.

describe('Avaliação: tratamento de erro ao enviar nota', () => {
  let rideId;

  before(function () {
    // Cria uma corrida para o cenário
    return cy.request({
      method: 'POST',
      url: 'http://localhost:3000/api/corrida',
      body: {
        PassageiroID: 654,
        Origem: '-8.050,-34.951',
        Destino: '-8.063,-34.871',
      },
    }).then((res) => {
      expect(res.status).to.eq(201);
      rideId = res.body.ID;
    });
  });

  it('exibe alert e mantém modo de avaliação quando POST falha', () => {
    // Garante que só exista um item no histórico neste teste
    cy.intercept('GET', 'http://localhost:3000/api/corridas', (req) => {
      req.reply([{
        id: rideId,
        origem: '-8.050,-34.951',
        destino: '-8.063,-34.871',
        dataInicio: new Date().toISOString(),
        preco: 24.80,
        avaliacao: null,
      }]);
    }).as('getCorridas');

    // Força erro no POST de avaliação
    cy.intercept('POST', `http://localhost:3000/api/corridas/${rideId}/avaliar`, {
      statusCode: 500,
      body: { error: 'falha ao salvar' },
    }).as('postAvaliacaoErro');

    // Captura o alert do browser
    let alerted = false;
    cy.on('window:alert', (txt) => {
      alerted = true;
      expect(txt).to.contain('Não foi possível enviar a avaliação');
    });

    cy.visit('/historico');
    cy.wait('@getCorridas');

    cy.get('.list .item .avaliar').click();
    cy.get('.list .item .stars .star-btn').eq(4 - 1).click(); // 4 estrelas por consistência

    cy.wait('@postAvaliacaoErro');

    // Continua em modo de edição (botões de estrela visíveis) e alert foi mostrado
    cy.get('.list .item .stars .star-btn').should('have.length', 5);
    cy.then(() => { expect(alerted).to.eq(true); });
    // E não deve exibir "Avaliado"
    cy.contains('.list .item .avaliar', 'Avaliado').should('not.exist');
  });
});