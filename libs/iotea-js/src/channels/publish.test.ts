import publish from './publish'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'

type ResponseData = null

describe('channels/publish', () => {
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
      spaceId: 'space-id-fake',
      channelId: 'channel-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await publish(
      'tea_fake',
      config,
      mockRequestData.spaceId,
      mockRequestData.channelId,
    )

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/channels/channel-id-fake/publish')
    expectedUrl.searchParams.set('spaceId', mockRequestData.spaceId)

    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'PATCH',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
