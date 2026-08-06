import list from './list'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { PermissionSet } from '@prisma/client'

type ResponseData = PermissionSet[]

describe('permissions/list', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return list of permission sets', async () => {
    // Given
    const mockResponseData: ResponseData = [
      {
        id: 'id-fake',
        organizationId: 'organization-id-fake',
        spaceId: null,
        name: 'test',
        permissions: ['*'],
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

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await list('tea_fake', config, mockRequestData.spaceId)

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/permissions')
    expectedUrl.searchParams.set('page', '1')
    expectedUrl.searchParams.set('resultsPerPage', '10')
    expectedUrl.searchParams.set('spaceId', mockRequestData.spaceId)

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
