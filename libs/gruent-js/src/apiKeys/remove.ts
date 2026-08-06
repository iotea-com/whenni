import { ClientConfig } from '..'
import apiRequest from '../request'
import { ApiResponse } from '../response'

type ResponseData = null

const removeApiKey = async (
  key: string,
  config: ClientConfig,
  orgId: string,
  apiKeyId: string,
  spaceId?: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/api-keys`

  const queryParams = new URLSearchParams({ orgId })
  if (spaceId) queryParams.set('spaceId', spaceId)

  return await apiRequest(url, 'DELETE', { apiKeyId }, { authorization: key, query: queryParams })
}

export default removeApiKey
