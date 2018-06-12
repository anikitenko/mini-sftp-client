describe('Security Test', function () {
    const sshUser = Cypress.env("mock_user"),
        sshPassword = Cypress.env("mock_pass"),
        sshPort = Cypress.env("mock_port"),
        sshHost = Cypress.env("mock_host"),
        postBodySSH = {
            ssh_ip: sshHost,
            ssh_user: sshUser,
            ssh_password: sshPassword,
            ssh_port: sshPort
        },
        postEndpoints = ["/download", "/connectViaSSH", "/testSSHConnection", "/getPath",
            "/createNewLocalDirectory", "/getLocalPathCompletion", "/localPathGoTo",
            "/getRemotePathCompletion", "/remotePathGoTo"];

    it('Access main page', function () {
        cy.visit("/")
        cy.get("#pinCode").should("not.be.empty")

        cy.get("#pinCode").then(($pinCode) => {
            cy.visit("/?for_testing=true");

            cy.focused().should('have.class', 'bootbox-input-number')
            cy.focused().type($pinCode.text()).should("have.value", $pinCode.text())
            cy.get("button[data-bb-handler='confirm']").click()

            cy.get('#sshIp').should("be.visible")

            cy.wrap(postEndpoints).each((postEndpoint) => {
                cy.request({
                    method: "POST",
                    url: postEndpoint+"/?for_testing=true",
                    failOnStatusCode: false
                })
                    .then((response) => {
                        expect(response.status).to.eq(200)
                    })
            })
        })
    })

    it('Try to test endpoints without pin code', function () {
        cy.wrap(postEndpoints).each((postEndpoint) => {
            cy.request({
                method: "POST",
                url: postEndpoint+"/?for_testing=true",
                failOnStatusCode: false
            })
                .then((response) => {
                    expect(response.status).to.eq(403)
                })
        })
    })
})