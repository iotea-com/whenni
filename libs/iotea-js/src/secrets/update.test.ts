import update from './update'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'

type ResponseData = null

describe('secrets/update', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return OK', async () => {
    // Given
    const mockResponseCode = 200
    const mockResponseBody = { data: null, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      value: 'test',
    }

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await update('tea_fake', config, 'space-id-fake', 'test', mockRequestData.value)

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/secrets/test')
    expectedUrl.searchParams.set('spaceId', 'space-id-fake')
    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      body: JSON.stringify(mockRequestData),
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'PUT',
    })
    expect(result.data).toEqual(null)
  })
})
