import { ClientConfig } from '..'
import apiRequest from '../request'
import { ApiResponse } from '../response'

type ResponseData = {
  status: 'PUBLISHED' | 'UNPUBLISHED' | 'INITIALIZING'
}

const channelStatus = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  channelId: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/channels/${channelId}/status`

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)

  return await apiRequest<ResponseData>(url, 'GET', null, {
    authorization: key,
    query: queryParams,
  })
}

export default channelStatus
