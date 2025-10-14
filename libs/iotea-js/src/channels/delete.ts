import { ClientConfig } from '..'
import apiRequest from '../request'
import { ApiResponse } from '../response'

type ResponseData = null

const deleteChannel = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  channelId: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/channels/${channelId}`

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)

  return await apiRequest(url, 'DELETE', null, { authorization: key, query: queryParams })
}

export default deleteChannel
