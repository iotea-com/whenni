import create from './create'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { Model } from '@prisma/client'

type ResponseData = Model

describe('models/create', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return created model', async () => {
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

    const mockResponseCode = 201
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      spaceId: 'space-id-fake',
      name: 'test',
      attributes: {},
    }

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await create('tea_fake', config, mockRequestData.spaceId, {
      name: mockRequestData.name,
      attributes: mockRequestData.attributes,
    })

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/models')
    expectedUrl.searchParams.set('spaceId', mockRequestData.spaceId)

    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      body: JSON.stringify({ name: mockRequestData.name, attributes: mockRequestData.attributes }),
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'POST',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
