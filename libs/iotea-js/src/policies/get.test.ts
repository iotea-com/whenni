import get from './get'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { Policy } from '.'

type ResponseData = {
  policy: Policy
  revoke: boolean
}

describe('policies/get', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return policy', async () => {
    // Given
    const mockResponseData: ResponseData = {
      policy: {
        allowedPublishTopics: [],
        allowedSubscriptionTopics: [],
      },
      revoke: false,
    }

    const mockResponseCode = 200
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      spaceId: 'space-id-fake',
      certificateId: 'certificate-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await get(
      'tea_fake',
      config,
      mockRequestData.spaceId,
      mockRequestData.certificateId,
    )

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/policies/certificate-id-fake')
    expectedUrl.searchParams.set('spaceId', mockRequestData.spaceId)

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
