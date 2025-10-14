import list from './list'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { ApiKey } from '@prisma/client'

type ResponseData = ApiKey[]

describe('apiKeys/list', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return list of API keys', async () => {
    // Given
    const mockResponseData: ResponseData = [
      {
        id: 'id-fake',
        organizationId: 'organization-id-fake',
        organizationPermissionSetId: 'organization-permission-set-id-fake',
        spaceId: null,
        spacePermissionSetId: null,
        name: 'test',
        createdAt: new Date(),
        createdBy: 'user-id-fake',
        expiresAt: null,
      },
    ]

    const mockResponseCode = 200
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      organizationId: 'organization-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await list('tea_fake', config, mockRequestData.organizationId)

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/api-keys')
    expectedUrl.searchParams.set('page', '1')
    expectedUrl.searchParams.set('resultsPerPage', '10')
    expectedUrl.searchParams.set('orgId', mockRequestData.organizationId)

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
