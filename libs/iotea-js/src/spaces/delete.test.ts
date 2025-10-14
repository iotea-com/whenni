import spacesDelete from './delete'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'

type ResponseData = null

describe('spaces/delete', () => {
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
      organizationId: 'organization-id-fake',
      spaceId: 'space-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await spacesDelete(
      'tea_fake',
      config,
      mockRequestData.organizationId,
      mockRequestData.spaceId,
    )

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/spaces/space-id-fake')
    expectedUrl.searchParams.set('orgId', mockRequestData.organizationId)

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
