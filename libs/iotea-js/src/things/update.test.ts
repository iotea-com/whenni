import update from './update'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { Thing } from '@prisma/client'

type ResponseData = Thing

describe('things/update', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return updated thing', async () => {
    // Given
    const mockResponseData: ResponseData = {
      id: 'thing-id-fake',
      spaceId: 'space-id-fake',
      name: 'test',
      thingCategory: 'CATEGORY_FAKE',
      attributes: {},
      internal: false,
      createdAt: new Date(),
      createdBy: 'user-id-fake',
      updatedAt: new Date(),
      updatedBy: 'user-id-fake',
    }

    const mockResponseCode = 200
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      spaceId: 'space-id-fake',
      name: mockResponseData.name,
      attributes: mockResponseData.attributes,
    }

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await update(
      'tea_fake',
      config,
      mockRequestData.spaceId,
      mockResponseData.id,
      mockRequestData,
    )

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/things/thing-id-fake')
    expectedUrl.searchParams.set('spaceId', mockRequestData.spaceId)
    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      body: JSON.stringify({
        name: mockRequestData.name,
        attributes: mockRequestData.attributes,
      }),
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'PUT',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
