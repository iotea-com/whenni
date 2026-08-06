import update from './update'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { Model } from '@prisma/client'

type ResponseData = Model

describe('models/update', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return updated model', async () => {
    // Given
    const mockResponseData: ResponseData = {
      id: 'model-id-fake',
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
      model: mockResponseData,
    }

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await update('tea_fake', config, mockRequestData.spaceId, mockRequestData.model)

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/models/model-id-fake')
    expectedUrl.searchParams.set('spaceId', mockRequestData.spaceId)

    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      body: JSON.stringify({ model: mockRequestData.model }),
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'PUT',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
