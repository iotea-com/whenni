import { ClientConfig } from '..'
import apiRequest from '../request'
import { ApiResponse } from '../response'

type ResponseData = null

const deleteSpaceById = async (
  key: string,
  config: ClientConfig,
  orgId: string,
  spaceId: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/spaces/${spaceId}`

  const queryParams = new URLSearchParams({ orgId })

  return await apiRequest<ResponseData>(url, 'DELETE', null, {
    authorization: key,
    query: queryParams,
  })
}

export default deleteSpaceById
