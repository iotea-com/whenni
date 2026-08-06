import deleteModel from './delete'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'

type ResponseData = null

describe('models/delete', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return OK', async () => {
    // Given
    const mockResponseData: ResponseData = null

    const mockResponseCode = 200
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      spaceId: 'space-id-fake',
      modelId: 'model-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await deleteModel(
      'tea_fake',
      config,
      mockRequestData.spaceId,
      mockRequestData.modelId,
    )

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/models/model-id-fake')
    expectedUrl.searchParams.set('spaceId', mockRequestData.spaceId)

    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'DELETE',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
