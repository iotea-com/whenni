import list, { ChannelExecutionListItem } from './list'
import { createApiCallMock } from '../../testHelpers'
import { ClientConfig } from '../../index'

type ResponseData = ChannelExecutionListItem[]

describe('channels/executions/list', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return list of channel executions for the channel', async () => {
    // Given
    const mockResponseData: ResponseData = [
      {
        executionId: 'channel-execution-id-fake',
        channelId: 'channel-id-fake',
        startTime: new Date().toISOString(),
        endTime: new Date().toISOString(),
        status: 'COMPLETED',
        durationMs: 1000,
        logAttributes: {},
      },
    ]

    const mockResponseCode = 200
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      spaceId: 'space-id-fake',
      channelId: 'channel-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await list(
      'tea_fake',
      config,
      mockRequestData.spaceId,
      mockRequestData.channelId,
    )

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/channels/executions')
    expectedUrl.searchParams.set('spaceId', mockRequestData.spaceId)
    expectedUrl.searchParams.set('channelId', mockRequestData.channelId)
    expectedUrl.searchParams.set('page', '1')
    expectedUrl.searchParams.set('resultsPerPage', '10')

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
