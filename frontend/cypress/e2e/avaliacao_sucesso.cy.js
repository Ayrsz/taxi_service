// cypress/e2e/avaliacao_sucesso.cy.js
describe('Avaliação: usuário avalia uma corrida com sucesso', () => {
  let rideId;

  before(() => {
    return cy.request('POST', 'http://localhost:3000/api/corrida', {
      PassageiroID: 987,
      Origem: '-8.050,-34.951',   // Centro de Informática
      Destino: '-8.063,-34.871',  // Recife Antigo
    }).then((res) => {
      expect(res.status).to.eq(201);
      rideId = res.body.ID;
    });
  });

  it('envia POST /api/corridas/:id/avaliar ao clicar na 4ª estrela', () => {
    // Espionar o POST (não stubar), com regex permissivo
    cy.intercept({ method: 'POST', url: /\/api\/corridas\/\d+\/avaliar(?:\?.*)?$/ }).as('postAvaliacao');

    // Apenas observar o GET para aguardar a lista renderizar
    cy.intercept('GET', '**/api/corridas*').as('getCorridas');

    cy.visit('/historico');
    cy.wait('@getCorridas');

    // Pega o item correto pela rota renderizada
    cy.contains('.list .item .route', 'Centro de Informática → Recife Antigo')
      .parents('.item')
      .as('item');

    // Entra no modo avaliação e clica na 4ª estrela
    cy.get('@item').within(() => {
      cy.get('.avaliar').click();
      cy.get('.stars .star-btn').should('have.length', 5).eq(3).click(); // 4 estrelas
    });

    // Confirma que o POST aconteceu e valida o payload
    cy.wait('@postAvaliacao').then(({ request, response }) => {
      expect(request.body).to.deep.equal({ nota: 4 });
      expect(response?.statusCode ?? 200).to.eq(200);
    });

    // UI marcada como avaliada
    cy.get('@item').within(() => {
      cy.get('.right .stars')
        .should('have.attr', 'aria-label')
        .and('contain', 'Nota 4 de 5');
      cy.contains('.avaliar', 'Avaliado').should('exist');
    });
  });
});
