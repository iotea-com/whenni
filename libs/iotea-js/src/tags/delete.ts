import { ClientConfig } from '..'
import apiRequest from '../request'
import { ApiResponse } from '../response'

type ResponseData = null

const deleteTagById = async (
  key: string,
  config: ClientConfig,
  orgId: string,
  tagId: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/tags/${tagId}`

  const queryParams = new URLSearchParams({ orgId })

  return await apiRequest<ResponseData>(url, 'DELETE', null, {
    authorization: key,
    query: queryParams,
  })
}

export default deleteTagById
