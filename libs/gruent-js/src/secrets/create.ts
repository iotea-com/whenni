import { ClientConfig } from '..'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = null

const createSecret = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  name: string,
  value: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/secrets`

  const requestBody = { name, value }

  const queryParams = new URLSearchParams({ spaceId })

  return await apiRequest<ResponseData>(url, 'POST', requestBody, {
    authorization: key,
    query: queryParams,
  })
}

export default createSecret
