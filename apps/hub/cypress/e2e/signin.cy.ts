describe('Sign in', () => {
  it('Sign in with magic link', () => {
    cy.visit('http://localhost:4200/signin')
    cy.get('#email').type('test@gruent.com')
    cy.get('#submit').click()
    cy.get('.toastNotification').contains('Magic link sent')
  })

  it('Sign in with userpass', () => {
    cy.visit('http://localhost:4200/signin')
    cy.get('#email').type('test@gruent.com')
    cy.get('#userpass').click()
    cy.get('#password').type('gruent!')
    cy.get('#submit').click()
  })
})
