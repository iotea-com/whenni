import changeRole from './changeRole'
import { createApiCallMock } from '../../testHelpers'
import { ClientConfig } from '../../index'

type ResponseData = null

describe('organizations/members/changeRole', () => {
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
      userId: 'user-id-fake',
      role: 'ADMIN' as 'ADMIN' | 'MEMBER',
    }

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await changeRole(
      'tea_fake',
      config,
      mockRequestData.organizationId,
      mockRequestData.userId,
      mockRequestData.role,
    )

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/organizations/members/changeRole')
    expectedUrl.searchParams.set('orgId', mockRequestData.organizationId)
    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      body: JSON.stringify({ userId: mockRequestData.userId, role: mockRequestData.role }),
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'POST',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
