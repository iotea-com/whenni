describe('Sign in', () => {
  it('Sign in with magic link', () => {
    cy.visit('http://localhost:4200/signin')
    cy.get('#email').type('test@iotea.com')
    cy.get('#submit').click()
    cy.get('.toastNotification').contains('Magic link sent')
  })

  it('Sign in with userpass', () => {
    cy.visit('http://localhost:4200/signin')
    cy.get('#email').type('test@iotea.com')
    cy.get('#userpass').click()
    cy.get('#password').type('iotea!')
    cy.get('#submit').click()
  })
})
