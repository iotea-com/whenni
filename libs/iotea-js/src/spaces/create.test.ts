import create from './create'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'
import { Space } from '@prisma/client'

type ResponseData = Space

describe('spaces/create', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return created space', async () => {
    // Given
    const mockResponseData: ResponseData = {
      id: 'space-id-fake',
      name: 'test',
      organizationId: 'organization-id-fake',
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
    }

    const config: ClientConfig = { url: new URL('https://api.iotea.com') }
    const result = await create('tea_fake', config, 'org-id-fake', mockRequestData.name)

    // Then
    const expectedUrl = new URL('https://api.iotea.com/v1/spaces')
    expectedUrl.searchParams.set('orgId', 'org-id-fake')

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
