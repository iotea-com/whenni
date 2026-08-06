import { ClientConfig } from '..'
import apiRequest from '../request'
import { ApiResponse } from '../response'

type ResponseData = null

const unpublishChannel = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  channelId: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/channels/${channelId}/unpublish`

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)

  return await apiRequest<ResponseData>(url, 'PATCH', null, {
    authorization: key,
    query: queryParams,
  })
}

export default unpublishChannel
