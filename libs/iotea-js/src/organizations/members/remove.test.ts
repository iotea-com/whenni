import remove from './remove'
import { createApiCallMock } from '../../testHelpers'
import { ClientConfig } from '../../index'

type ResponseData = null

describe('organizations/members/remove', () => {
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
      userId: 'user-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await remove('tea_fake', config, mockRequestData.orgId, mockRequestData.userId)

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/organizations/members')
    expectedUrl.searchParams.set('orgId', mockRequestData.orgId)
    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      body: JSON.stringify({ userId: mockRequestData.userId }),
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'DELETE',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
