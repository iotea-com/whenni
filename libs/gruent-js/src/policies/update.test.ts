import update from './update'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { Policy } from '.'

type ResponseData = {
  policy: Policy
  revoke: boolean
}

describe('policies/update', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return updated policy', async () => {
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
      policy: {
        allowedPublishTopics: [],
        allowedSubscriptionTopics: [],
      },
      revoke: false,
    }

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await update(
      'tea_fake',
      config,
      mockRequestData.spaceId,
      mockRequestData.certificateId,
      mockRequestData.policy,
      mockRequestData.revoke,
    )

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/policies/certificate-id-fake')
    expectedUrl.searchParams.set('spaceId', mockRequestData.spaceId)

    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      body: JSON.stringify(mockResponseData),
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'PUT',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
