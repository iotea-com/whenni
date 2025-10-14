import list from './list'
import { createApiCallMock } from '../../testHelpers'
import { ClientConfig } from '../../index'
import { OrganizationMember, User } from '@prisma/client'

type ResponseData = (OrganizationMember & {
  user: User
})[]

describe('organizations/members/list', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return list of users in the space', async () => {
    // Given
    const mockResponseData: ResponseData = []

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
    const expectedUrl = new URL('https://api.iotea.com/v1/organizations/members')
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
