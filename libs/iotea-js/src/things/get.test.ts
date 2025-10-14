import get from './get'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { Thing } from '@prisma/client'

type ResponseData = Thing

describe('things/get', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return thing', async () => {
    // Given
    const mockResponseData: ResponseData = {
      id: 'thing-id-fake',
      name: 'test',
      spaceId: 'space-id-fake',
      attributes: {},
      internal: false,
      thingCategory: 'CATEGORY_FAKE',
      createdBy: 'user-id-fake',
      updatedBy: 'user-id-fake',
      createdAt: new Date(),
      updatedAt: new Date(),
    }

    const mockResponseCode = 200
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      spaceId: 'space-id-fake',
      thingId: 'thing-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await get('tea_fake', config, mockRequestData.spaceId, mockRequestData.thingId)

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/things/thing-id-fake')
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
