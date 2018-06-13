describe('Remote and local navigation buttons', function () {

    before(function() {
        cy.visit("/")
        cy.fillMainForm()

        cy.get(".mainForm")
            .submit()

        cy.contains("Loading remote files...")
        cy.contains("Loading local files...")
    })

    it("Check remote back button", function() {
        cy.get(".remoteFilesBlock").find("span[data-dir='true']:first").trigger("click")
        cy.contains("Loading remote files...")
        cy.get(".remoteGoBack").trigger("click")
        cy.contains("Loading remote files...")
    })

    it("Check local back button", function() {
        cy.get(".localFilesBlock").find("span[data-dir='true']:first").trigger("click")
        cy.contains("Loading local files...")
        cy.get(".localGoBack").trigger("click")
        cy.contains("Loading local files...")
    })

    it("Check remote home button", function() {
        cy.get(".remoteGoHome").trigger("click")
        cy.contains("Loading remote files...")
    })

    it("Check local home button", function() {
        cy.get(".localGoHome").trigger("click")
        cy.contains("Loading local files...")
    })

    it("Check remote up button", function() {
        cy.get(".remoteGoUp").trigger("click")
        cy.contains("Loading remote files...")
        cy.get(".remoteGoBack").trigger("click")
        cy.contains("Loading remote files...")
    })

    it("Check local up button", function() {
        cy.get(".localGoUp").trigger("click")
        cy.contains("Loading local files...")
        cy.get(".localGoBack").trigger("click")
        cy.contains("Loading local files...")
    })

    it("Check remote refresh button", function() {
        cy.get(".remoteRefresh").trigger("click")
        cy.contains("Loading remote files...")
    })

    it("Check local refresh button", function() {
        cy.get(".localRefresh").trigger("click")
        cy.contains("Loading local files...")
    })

    it("Check create local folder", function() {
        cy.get(".localCreateNewDir").trigger("click")
        cy.get(".bootbox-input-text").type("just test local dir").should("have.value", "just test local dir")
        cy.get("button[data-bb-handler='confirm']").click()
        cy.contains("Loading local files...")
    })
})