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
        };

    it('Access main page', function () {
        cy.visit("/?for_testing=true", {
            onBeforeLoad(win) {
                cy.stub(win, 'prompt')
            },
            failOnStatusCode: false
        });

        cy.window().its('prompt').should('be.called')
    })

    it('Try to test connection', function () {
        cy.request({
            method: "POST",
            url: "/testSSHConnection/?for_testing=true",
            form: true,
            body: postBodySSH,
            failOnStatusCode: false
        })
            .then((response) => {
                expect(response.status).to.eq(403)
            })
    })
})