import { customAlphabet } from 'nanoid'

const nanoidAlphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'

describe('Spaces', () => {
  beforeEach(() => {
    cy.login('test@iotea.com', 'iotea!')
  })

  const suiteCache = new Map([])

  it('Create a space', () => {
    const nanoid = customAlphabet(nanoidAlphabet)
    const id = nanoid(5)
    const spaceName = `Test space ${id}`

    suiteCache.set('spaceName', spaceName)

    cy.visit('http://localhost:4200/dashboard')
    cy.get('#createSpaceButton').click()
    cy.get('#name').type(spaceName)
    cy.get('#createSpaceSubmit').click()
    cy.get('.toastNotification').contains('Successfully created a new space!')
    cy.get('#spaces .spaceCard').contains(spaceName)
  })

  it('Delete a space', () => {
    const spaceName = suiteCache.get('spaceName')
    if (!spaceName) throw new Error('No spaceName in suiteCache')

    // Navigate to space
    cy.visit('http://localhost:4200/dashboard')
    cy.get('#spaces .spaceCard')
      .contains(spaceName as string)
      .parent()
      .click()

    // Delete space
    cy.get('#spacesNavbar #spaceSettings').click()
    cy.get('#deleteSpaceButton').click()
    cy.get('#deleteSpaceSubmit').click()
  })
})
