describe('MainForm', function () {
    it('Visit main page', function () {
        cy.visit("http://127.0.0.1:9000");

        cy.title().should('include', 'New Connection');

        cy.get("#connectionNameDisplay").click();

        cy.get(".editableform").find(".editable-input > input").type("Test");
        cy.get(".editableform").submit();

        cy.title().should('include', 'Test');

        cy.get("#connectionNameDisplay").click();

        cy.get(".editableform").find(".editable-input > input").clear().type("Another test");
        cy.get('#sshIp').click();

        cy.title().should('include', 'Another test')
    });

    it('Fill in main form', function () {
        cy.get('#sshIp')
            .type('sftp-mock-test')
            .should('have.value', 'sftp-mock-test');

        cy.get("#sshUser")
            .type('test')
            .should('have.value', 'test');
        cy.get("#sshPassword")
            .type('test')
            .should('have.value', 'test');
        cy.get("#sshPort")
            .type('22')
            .should('have.value', '2222')
    })
});

/*
describe('My First Test', function() {
    it("Gets, types and asserts", function() {
        cy.visit('https://example.cypress.io')

        cy.contains('type').click()

        // Should be on a new URL which includes '/commands/actions'
        cy.url().should('include', '/commands/actions')

        // Get an input, type into it and verify that the value has been updated
        cy.get('.action-email')
            .type('fake@email.com')
            .should('have.value', 'fake@email.com')
    })
})*/
