import { ClientConfig } from '..'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = null

const updateSecret = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  name: string,
  value: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/secrets/${name}`

  const requestBody = { value }

  const queryParams = new URLSearchParams({ spaceId })

  return await apiRequest<ResponseData>(url, 'PUT', requestBody, {
    authorization: key,
    query: queryParams,
  })
}

export default updateSecret
