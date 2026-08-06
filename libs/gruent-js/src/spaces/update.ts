import { ClientConfig } from '..'
import { Space } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = Space

const updateSpaceById = async (
  key: string,
  config: ClientConfig,
  orgId: string,
  spaceId: string,
  space: Space,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/spaces/${spaceId}`

  const queryParams = new URLSearchParams({ orgId })

  return await apiRequest<ResponseData>(
    url,
    'PUT',
    { space },
    {
      authorization: key,
      query: queryParams,
    },
  )
}

export default updateSpaceById
