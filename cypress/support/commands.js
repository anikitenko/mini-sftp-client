// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add("login", (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add("drag", { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add("dismiss", { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This is will overwrite an existing command --
// Cypress.Commands.overwrite("visit", (originalFn, url, options) => { ... })
const sshUser = Cypress.env("mock_user"),
    sshPassword = Cypress.env("mock_pass"),
    sshPort = Cypress.env("mock_port"),
    sshHost = Cypress.env("mock_host");

Cypress.Commands.add('fillMainForm', () => {
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
})