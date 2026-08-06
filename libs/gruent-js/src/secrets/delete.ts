import { ClientConfig } from '..'
import apiRequest from '../request'
import { ApiResponse } from '../response'

type ResponseData = null

const deleteSecret = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  name: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/secrets/${name}`

  const queryParams = new URLSearchParams({ spaceId })

  return await apiRequest<ResponseData>(url, 'DELETE', null, {
    authorization: key,
    query: queryParams,
  })
}

export default deleteSecret
