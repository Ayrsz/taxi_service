// cypress/e2e/historico_lista.cy.js
// Objetivo: Renderizar a lista do Histórico e exibir itens quando há corridas.

describe('Histórico: renderização da lista com corridas', () => {
  let created = [];

  before(function () {
    // Cria duas corridas reais na API (seguindo o estilo dos cenários existentes)
    const makeRide = (origem, destino) => cy.request({
      method: 'POST',
      url: 'http://localhost:3000/api/corrida',
      body: {
        PassageiroID: 321,
        Origem: origem,
        Destino: destino,
      },
    }).then((res) => {
      expect(res.status).to.eq(201);
      expect(res.body).to.have.property('ID');
      created.push(res.body);
    });

    return makeRide('-8.050,-34.951', '-8.063,-34.871')
      .then(() => makeRide('-8.050,-34.951', '-8.063,-34.871'));
  });

  it('exibe título e pelo menos um item na lista', () => {
    // Garante que o GET não use cache e possamos aguardar
    cy.intercept('GET', 'http://localhost:3000/api/corridas').as('getCorridas');

    cy.visit('/historico');
    cy.wait('@getCorridas');

    cy.get('.title').should('contain', 'Histórico de Corridas');
    cy.get('.list .item').its('length').should('be.gte', 1);

    // Cada item deve ter rota, data ou preço e a área da direita
    cy.get('.list .item').first().within(() => {
      cy.get('.route').should('exist');
      cy.get('.right').should('exist');
    });
  });
});