import get, { ChannelExecution } from './get'
import { createApiCallMock } from '../../testHelpers'
import { ClientConfig } from '../../index'

type ResponseData = ChannelExecution

describe('channels/executions/get', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return channel execution', async () => {
    // Given
    const mockResponseData: ResponseData = {
      executionId: 'channel-execution-id-fake',
      channelId: 'channel-id-fake',
      startTime: new Date().toISOString(),
      endTime: new Date().toISOString(),
      durationMs: 1000,
      status: 'COMPLETED',
      nodeExecutionLogs: [],
    }

    const mockResponseCode = 200
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      spaceId: 'space-id-fake',
      channelId: 'channel-id-fake',
      channelExecutionId: 'channel-execution-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await get(
      'tea_fake',
      config,
      mockRequestData.spaceId,
      mockRequestData.channelExecutionId,
    )

    // Then
    const expectedUrl = new URL(
      'https://api.gruent.com/v1/channels/executions/channel-execution-id-fake',
    )
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
