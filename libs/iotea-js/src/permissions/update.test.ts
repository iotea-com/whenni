import update from './update'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { PermissionSet } from '@prisma/client'

type ResponseData = null

describe('permissions/update', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return updated permission set', async () => {
    // Given
    const mockResponseData: ResponseData = null

    const mockResponseCode = 200
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData: PermissionSet = {
      id: 'permission-set-id-fake',
      organizationId: 'organization-id-fake',
      spaceId: null,
      name: 'test',
      permissions: [],
      createdAt: new Date(),
      createdBy: 'user-id-fake',
      updatedAt: new Date(),
      updatedBy: 'user-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await update(
      'tea_fake',
      config,
      mockRequestData.organizationId,
      mockRequestData.id,
      mockRequestData.name,
      mockRequestData.permissions,
    )

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/permissions/permission-set-id-fake')
    expectedUrl.searchParams.set('orgId', mockRequestData.organizationId)

    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      body: JSON.stringify({
        name: mockRequestData.name,
        permissions: mockRequestData.permissions,
      }),
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'PUT',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
