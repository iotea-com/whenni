import add from './add'
import { createApiCallMock } from '../../testHelpers'
import { ClientConfig } from '../../index'

type ResponseData = null

describe('organizations/members/add', () => {
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
    }

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await add(
      'tea_fake',
      config,
      mockRequestData.organizationId,
      mockRequestData.userId,
    )

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/organizations/members')
    expectedUrl.searchParams.set('orgId', mockRequestData.organizationId)
    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      body: JSON.stringify({ userId: mockRequestData.userId }),
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'POST',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
