describe('Remote and local navigation', function () {

    before(function() {
        cy.visit("/")
        cy.fillMainForm()

        cy.get(".mainForm")
            .submit()

        cy.contains("Loading remote files...")
        cy.contains("Loading local files...")
    })

    it('Check remote path', function() {
        cy.get("#select2-remotePath-container").then(($remotePath) => {
            cy.get(".remoteFilesNavigationBlock").find(".select2").click()
            cy.focused().should("have.value", $remotePath.text())
            cy.focused().type("{backspace}")
            cy.wait(200)
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
            cy.focused().should("have.value", $localPath.text())
            cy.focused().type("{backspace}")
            cy.wait(200)
            cy.focused().parent().parent().get(".select2-results > ul")
                .find("li").not(":contains('No results found')").should("not.be.empty")

            cy.focused().type("some_random_string")
            cy.focused().parent().parent().get(".select2-results > ul")
                .find("li").should("have.text", "No results found")
        })
    })
})