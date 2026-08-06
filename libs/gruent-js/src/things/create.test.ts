import create from './create'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { Thing } from '@prisma/client'

type ResponseData = Thing

describe('things/create', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return created thing', async () => {
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

    const mockResponseCode = 201
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      spaceId: 'space-id-fake',
      name: 'test',
      category: 'CATEGORY_FAKE',
      attributes: {},
    }

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await create(
      'tea_fake',
      config,
      mockRequestData.spaceId,
      mockRequestData.name,
      mockRequestData.category,
      mockRequestData.attributes,
    )

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/things')
    expectedUrl.searchParams.set('spaceId', mockRequestData.spaceId)
    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      body: JSON.stringify({
        name: mockRequestData.name,
        category: mockRequestData.category,
        attributes: mockRequestData.attributes,
      }),
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'POST',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
