import get from './get'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { Organization } from '@prisma/client'

type ResponseData = Organization

describe('organizations/get', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return organization', async () => {
    // Given
    const mockResponseData: ResponseData = {
      id: 'id-fake',
      name: 'test',
      createdAt: new Date(),
      createdBy: 'user-id-fake',
      updatedAt: new Date(),
      updatedBy: 'user-id-fake',
    }

    const mockResponseCode = 200
    const mockResponseBody = { data: mockResponseData, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      orgId: 'org-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await get('tea_fake', config, mockRequestData.orgId)

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/organizations/org-id-fake')

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
