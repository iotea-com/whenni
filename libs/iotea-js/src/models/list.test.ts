import list from './list'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { Model } from '@prisma/client'

type ResponseData = Model[]

describe('models/list', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return list of models in the space', async () => {
    // Given
    const mockResponseData: ResponseData = [
      {
        id: 'model-id-fake',
        spaceId: 'space-id-fake',
        name: 'test',
        attributes: {},
        createdAt: new Date(),
        createdBy: 'user-id-fake',
        updatedAt: new Date(),
        updatedBy: 'user-id-fake',
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
    const expectedUrl = new URL('https://api.iotea.com/v1/models')
    expectedUrl.searchParams.set('spaceId', mockRequestData.spaceId)
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
