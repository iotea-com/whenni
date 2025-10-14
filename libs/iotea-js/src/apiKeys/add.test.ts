import add from './add'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { ApiKey } from '@prisma/client'

type ResponseData = ApiKey

describe('apiKeys/add', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return API key when it has successfully been added', async () => {
    // Given
    const mockResponseData: ResponseData = {
      id: 'id-fake',
      organizationId: 'organization-id-fake',
      organizationPermissionSetId: 'organization-permission-set-id-fake',
      spaceId: null,
      spacePermissionSetId: null,
      name: 'test',
      createdAt: new Date(),
      createdBy: 'user-id-fake',
      expiresAt: null,
    }

    const mockResponseCode = 200
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      name: 'test',
      organizationId: 'organization-id-fake',
      permissionSetId: 'permission-set-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await add(
      'tea_fake',
      config,
      mockRequestData.organizationId,
      mockRequestData.permissionSetId,
      mockRequestData.name,
    )

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/api-keys')
    expectedUrl.searchParams.set('orgId', mockRequestData.organizationId)

    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      body: JSON.stringify({
        permissionSetId: mockRequestData.permissionSetId,
        name: mockRequestData.name,
      }),
      method: 'POST',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
