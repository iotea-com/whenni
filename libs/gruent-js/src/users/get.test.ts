import get from './get'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { User } from '@prisma/client'

type ResponseData = User & {
  spaces?: {
    spaceId: string
    userId: string
    space: {
      id: string
      name: string
    }
  }[]
}

describe('users/get', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return success when user exists', async () => {
    // Given
    const mockResponseData: ResponseData = {
      id: 'id-fake',
      email: 'email-fake',
      name: null,
      emailVerified: null,
      password: null,
      image: null,
    }

    const mockResponseCode = 200
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await get('tea_fake', config, 'user-id-fake')

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/users/user-id-fake')

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
