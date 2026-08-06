import remove from './remove'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'

type ResponseData = null

describe('permissions/remove', () => {
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
      orgId: 'org-id-fake',
      permissionSetId: 'permission-set-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await remove(
      'tea_fake',
      config,
      mockRequestData.orgId,
      mockRequestData.permissionSetId,
    )

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/permissions')
    expectedUrl.searchParams.set('orgId', mockRequestData.orgId)

    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      body: JSON.stringify({
        permissionSetId: mockRequestData.permissionSetId,
      }),
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'DELETE',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
