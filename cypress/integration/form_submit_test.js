context('MainForm', () => {
    it('Visit main page, enter data and submit form', () => {
        cy.visit("http://127.0.0.1:9000");
        cy.get('#sshIp').type('sftp-mock-test').should('have.value', 'sftp-mock-test');
        cy.get("#sshUser").type('test').should('have.value', 'test');
        cy.get("#sshPassword").type('test').should('have.value', 'test');
        cy.get("#sshPort").type('22').should('have.value', '2222');
    });
});