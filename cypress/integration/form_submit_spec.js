describe('MainForm Test', function () {
    const sshUser = Cypress.env("mock_user"),
        sshPort = Cypress.env("mock_port"),
        sshHost = Cypress.env("mock_host");

    it('Visit main page', function () {
        cy.visit("/");

        cy.title()
            .should('include', 'New Connection');

        cy.get("#connectionNameDisplay")
            .click();

        cy.get(".editableform")
            .find(".editable-input > input")
            .type("Test");

        cy.get(".editableform")
            .submit();

        cy.title()
            .should('include', 'Test');

        cy.get("#connectionNameDisplay")
            .click();

        cy.get(".editableform")
            .find(".editable-input > input")
            .clear()
            .type("Another test");

        cy.get('#sshIp')
            .click();

        cy.title().should('include', 'Another test')
    });

    it('Fill in main form', function () {
        cy.fillMainForm()
    });

    it('Check main buttons', function () {
        cy.get("#testSSHConnection")
            .click()

        cy.contains("SSH connection was established successfully to '" + sshHost + ":" + sshPort + "'")

        cy.get("#sshUser")
            .type('123')
            .should("have.value", sshUser+"123")

        cy.get("#testSSHConnection")
            .click()

        cy.contains("We could not reach '127.0.0.1:2222' OR login/password is incorrect")

        cy.get(".mainForm")
            .submit()

        cy.contains("We could not reach '127.0.0.1:2222' OR login/password is incorrect")

        cy.get("#sshUser")
            .clear()
            .type(sshUser)
            .should("have.value", sshUser)

        cy.get(".mainForm")
            .submit()

        cy.contains("Loading remote files...")
        cy.contains("Loading local files...")

        cy.get("#remoteConnectionName").trigger("mouseover")
        cy.get('.tooltip').should('be.visible')
    })
});