describe('MainForm Requests Test', function () {
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

    it('Check test connection request', function () {
        cy.request({
            method: "POST",
            url: "/testSSHConnection",
            form: true,
            body: postBodySSH
        })
            .then((response) => {
                expect(response.status).to.eq(200)
                expect(response.body.result).to.be.true
                expect(response.body.message).to.be.empty
            })
    })

    it('Check connect request', function () {
        cy.request({
            method: "POST",
            url: "/connectViaSSH",
            form: true,
            body: postBodySSH
        })
            .then((response) => {
                expect(response.status).to.eq(200)
                expect(response.body.result).to.be.true
                expect(response.body.message).to.be.empty
                expect(response.body.errors).to.be.null
                expect(response.body.local_path).not.to.be.empty
                expect(response.body.local_path_separator).not.to.be.empty
                expect(response.body.remote_path).not.to.be.empty
            })
    })
})