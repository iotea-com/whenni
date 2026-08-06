import remove from './remove'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'

type ResponseData = null

describe('apiKeys/remove', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return OK', async () => {
    // Given
    const mockResponseData: ResponseData = null

    const mockResponseCode = 200
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      orgId: 'org-id-fake',
      apiKeyId: 'api-key-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await remove('tea_fake', config, mockRequestData.orgId, mockRequestData.apiKeyId)

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/api-keys')
    expectedUrl.searchParams.set('orgId', mockRequestData.orgId)

    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      body: JSON.stringify({ apiKeyId: 'api-key-id-fake' }),
      method: 'DELETE',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
