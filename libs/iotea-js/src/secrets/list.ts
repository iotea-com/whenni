import { ClientConfig } from '..'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = { name: string }[]

const getSpaceById = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/secrets`

  const queryParams = new URLSearchParams({ spaceId })

  return await apiRequest<ResponseData>(url, 'GET', null, {
    authorization: key,
    query: queryParams,
  })
}

export default getSpaceById
