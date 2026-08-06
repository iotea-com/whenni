import validate from './validate'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { ChannelLifetime } from '@gruent/libs/engine/channels/channels'

type ResponseData = string[]

describe('channels/validate', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return OK', async () => {
    // Given
    const mockResponseData: ResponseData = []

    const mockResponseCode = 200
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      spaceId: 'space-id-fake',
      channelId: 'channel-id-fake',
      config: {
        id: 'channel-id-fake',
        name: 'test',
        lifetime: 'single' as ChannelLifetime,
        edges: [],
        nodes: [],
        notes: [],
      },
    }

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await validate(
      'tea_fake',
      config,
      mockRequestData.spaceId,
      mockRequestData.config,
    )

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/channels/validate')
    expectedUrl.searchParams.set('spaceId', mockRequestData.spaceId)

    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      body: JSON.stringify(mockRequestData.config),
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'PATCH',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
