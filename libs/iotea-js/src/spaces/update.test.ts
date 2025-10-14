import update from './update'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { Space } from '@prisma/client'

type ResponseData = Space

describe('spaces/update', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return OK', async () => {
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
    const mockRequestData = mockResponseData

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await update(
      'tea_fake',
      config,
      mockRequestData.organizationId,
      mockRequestData.id,
      mockRequestData,
    )

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/spaces/id-fake')
    expectedUrl.searchParams.set('orgId', mockRequestData.organizationId)

    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      body: JSON.stringify({ space: mockRequestData }),
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'PUT',
    })
    expect(result.data).toEqual(mockRequestData)
  })
})
