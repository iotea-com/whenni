import search from './search'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { User } from '@prisma/client'

type ResponseData = User[]

describe('users/search', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return users when search is successful', async () => {
    // Given
    const mockResponseData: ResponseData = [
      {
        id: 'id-fake',
        email: 'email-fake',
        name: null,
        emailVerified: null,
        password: null,
        image: null,
      },
    ]

    const mockResponseCode = 200
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await search('tea_fake', config, 'email-fake')

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/users/search')
    expectedUrl.searchParams.set('q', 'email-fake')

    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'GET',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
