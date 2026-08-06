import add from './add'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { PermissionSet } from '@prisma/client'

type ResponseData = PermissionSet

describe('permissions/add', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return created permission set', async () => {
    // Given
    const mockResponseData: ResponseData = {
      id: 'id-fake',
      organizationId: 'organization-id-fake',
      spaceId: null,
      name: 'test',
      permissions: ['*'],
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
      name: 'test',
      permissions: ['*'],
    }

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await add(
      'tea_fake',
      config,
      mockResponseData.organizationId,
      mockRequestData.name,
      mockRequestData.permissions,
    )

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/permissions')
    expectedUrl.searchParams.set('orgId', mockResponseData.organizationId)

    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      body: JSON.stringify(mockRequestData),
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'POST',
    })
    expect(result.data).toEqual(mockResponseData)
  })
})
