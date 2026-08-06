import create from './create'
import { createApiCallMock } from '../testHelpers'
import { ClientConfig } from '../index'

type ResponseData = null

describe('secrets/create', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })

  it('should return CREATED', async () => {
    // Given
    const mockResponseCode = 201
    const mockResponseBody = { data: null, errors: null }
    createApiCallMock<ResponseData>(mockResponseCode, mockResponseBody)

    // When
    const mockRequestData = {
      name: 'test',
      value: 'test',
    }

    const config: ClientConfig = { url: new URL('https://api.gruent.com') }
    const result = await create(
      'tea_fake',
      config,
      'space-id-fake',
      mockRequestData.name,
      mockRequestData.value,
    )

    // Then
    const expectedUrl = new URL('https://api.gruent.com/v1/secrets')
    expectedUrl.searchParams.set('spaceId', 'space-id-fake')
    expect(global.fetch).toHaveBeenCalledTimes(1)
    expect(global.fetch).toHaveBeenCalledWith(expectedUrl.toString(), {
      body: JSON.stringify(mockRequestData),
      headers: {
        Authorization: 'Bearer tea_fake',
      },
      method: 'POST',
    })
    expect(result.data).toEqual(null)
  })
})
