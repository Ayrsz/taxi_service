// cypress/e2e/historico_vazio.cy.js
// Objetivo: Quando não há corridas, mostrar estado vazio.

describe('Histórico: estado vazio quando API retorna array vazio', () => {
  it('mostra mensagem de vazio', () => {
    cy.intercept('GET', 'http://localhost:3000/api/corridas', {
      statusCode: 200,
      body: [],
    }).as('getCorridas');

    cy.visit('/historico');
    cy.wait('@getCorridas');

    cy.get('.empty').should('contain', 'Você ainda não tem corridas.');
    cy.get('.list .item').should('not.exist');
  });
});