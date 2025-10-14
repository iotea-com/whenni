import list from './list'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { Channel } from '@prisma/client'

type ResponseData = Channel[]

describe('channels/list', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return list of channels in the space', async () => {
    // Given
    const mockResponseData: ResponseData = [
      {
        id: 'channel-execution-id-fake',
        name: 'test',
        spaceId: 'channel-id-fake',
        config: {},
        createdAt: new Date(),
        createdBy: 'user-id-fake',
        updatedAt: new Date(),
        updatedBy: 'user-id-fake',
        publishedAt: null,
        publishedBy: null,
      },
    ]

    const mockResponseCode = 200
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      spaceId: 'space-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await list('tea_fake', config, mockRequestData.spaceId)

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/channels')
    expectedUrl.searchParams.set('page', '1')
    expectedUrl.searchParams.set('resultsPerPage', '10')
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
