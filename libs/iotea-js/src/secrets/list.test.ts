import list from './list'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'

type ResponseData = string[]

describe('spaces/get', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return secrets', async () => {
    // Given

    const mockResponseCode = 200
    const mockResponseBody = { data: null, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      spaceId: 'space-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await list('tea_fake', config, mockRequestData.spaceId)

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/secrets')
    expectedUrl.searchParams.set('spaceId', mockRequestData.spaceId)

    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'GET',
    })
    expect(result.data).toEqual(mockResponseBody.data)
  })
})
