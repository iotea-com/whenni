import updateConfig from './updateConfig'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { ChannelLifetime } from '@gruent/libs/engine/channels/channels'
import { Channel } from '@prisma/client'

type ResponseData = Channel

describe('channels/updateConfig', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return OK', async () => {
    // Given
    const mockResponseData: ResponseData = {
      id: 'channel-id-fake',
      spaceId: 'space-id-fake',
      name: 'test',
      config: {},
      createdAt: new Date(),
      createdBy: 'user-id-fake',
      updatedAt: new Date(),
      updatedBy: 'user-id-fake',
      publishedAt: null,
      publishedBy: null,
    }

    const mockResponseCode = 200
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      channelId: 'channel-id-fake',
      spaceId: 'space-id-fake',
      config: {
        id: 'channel-id-fake',
        name: 'test',
        lifetime: 'single' as ChannelLifetime,
        nodes: [],
        edges: [],
        notes: [],
      },
    }

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await updateConfig(
      'tea_fake',
      config,
      mockRequestData.spaceId,
      mockRequestData.channelId,
      mockRequestData.config,
    )

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/channels/channel-id-fake/config')
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
