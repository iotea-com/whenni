import create from './create'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { Channel } from '@prisma/client'

type ResponseData = Channel

describe('channels/create', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return created permission set', async () => {
    // Given
    const mockResponseData: ResponseData = {
      id: 'id-fake',
      name: 'test',
      spaceId: 'space-id-fake',
      config: {},
      createdAt: new Date(),
      createdBy: 'user-id-fake',
      updatedAt: new Date(),
      updatedBy: 'user-id-fake',
      publishedAt: null,
      publishedBy: null,
    }

    const mockResponseCode = 201
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      spaceId: 'space-id-fake',
      name: '',
      nodes: [],
      edges: [],
    }

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await create('tea_fake', config, mockRequestData.spaceId, mockRequestData)

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/channels')
    expectedUrl.searchParams.set('spaceId', mockRequestData.spaceId)

    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      body: JSON.stringify({
        name: mockRequestData.name,
        config: {
          nodes: mockRequestData.nodes,
          edges: mockRequestData.edges,
        },
      }),
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'POST',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
