import setCookie from 'set-cookie-parser'

Cypress.Commands.add('login', () => {
  cy.request({
    method: 'POST',
    url: 'http://localhost:4200/api/signin',
    body: { email: 'test@iotea.com', password: 'iotea!', method: 'password' },
  }).then((res) => {
    if (!res.headers['set-cookie']) throw new Error('did not recognize cookies in signin response')

    const cookies = setCookie.parse(res.headers['set-cookie'], {
      decodeValues: true, // default: true
    })

    cookies.forEach((cookie) => {
      cy.setCookie(cookie.name, cookie.value)
    })
  })
})
