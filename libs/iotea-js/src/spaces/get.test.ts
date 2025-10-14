import get from './get'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { Space } from '@prisma/client'

type ResponseData = Space

describe('spaces/get', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return space', async () => {
    // Given
    const mockResponseData: ResponseData = {
      id: 'id-fake',
      name: 'test',
      organizationId: 'organization-id-fake',
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
      organizationId: 'organization-id-fake',
      spaceId: 'space-id-fake',
    }

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await get(
      'tea_fake',
      config,
      mockRequestData.organizationId,
      mockRequestData.spaceId,
    )

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/spaces/space-id-fake')
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
