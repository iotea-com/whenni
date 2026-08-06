import get from './get'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { Model } from '@prisma/client'

type ResponseData = Model

describe('models/get', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return model', async () => {
    // Given
    const mockResponseData: ResponseData = {
      id: 'thing-id-fake',
      spaceId: 'space-id-fake',
      name: 'test',
      attributes: {},
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
      modelId: 'model-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await get('tea_fake', config, mockRequestData.spaceId, mockRequestData.modelId)

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/models/model-id-fake')
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
