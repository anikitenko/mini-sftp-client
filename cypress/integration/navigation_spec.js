describe('Remote and local navigation', function () {
    const sshUser = Cypress.env("mock_user"),
        sshPassword = Cypress.env("mock_pass"),
        sshPort = Cypress.env("mock_port"),
        sshHost = Cypress.env("mock_host");

    before(function() {
        cy.visit("/")
        cy.get('#sshIp')
            .type(sshHost)
            .should('have.value', sshHost);

        cy.get("#sshUser")
            .type(sshUser)
            .should('have.value', sshUser);

        cy.get("#sshPassword")
            .type(sshPassword)
            .should('have.value', sshPassword);

        cy.get("#sshPort")
            .clear()
            .type(sshPort)
            .should('have.value', sshPort)

        cy.get(".mainForm")
            .submit()

        cy.contains("Loading remote files...")
        cy.contains("Loading local files...")
    })

    it('Check remote path', function() {
        cy.get("#select2-remotePath-container").then(($remotePath) => {
            cy.get(".remoteFilesNavigationBlock").find(".select2").click()
            cy.focused().should("have.value", $remotePath.attr("title"))
            cy.focused().type("{backspace}")
            cy.wait(100)
            cy.focused().parent().parent().get(".select2-results > ul")
                .find("li").not(":contains('No results found')").should("not.be.empty")

            cy.focused().type("some_random_string")
            cy.focused().parent().parent().get(".select2-results > ul")
                .find("li").should("have.text", "No results found")
        })
    })

    it('Check local path', function() {
        cy.get("#select2-localPath-container").then(($localPath) => {
            cy.get(".localFilesNavigationBlock").find(".select2").click()
            cy.focused().should("have.value", $localPath.attr("title"))
            cy.focused().type("{backspace}")
            cy.wait(100)
            cy.focused().parent().parent().get(".select2-results > ul")
                .find("li").not(":contains('No results found')").should("not.be.empty")

            cy.focused().type("some_random_string")
            cy.focused().parent().parent().get(".select2-results > ul")
                .find("li").should("have.text", "No results found")
        })
    })
})