import { ChannelConfig } from '@gruent/libs/engine/channels/channels'
import { ClientConfig } from '..'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = {
  channelErrors: string[] | null
  nodeErrors: Record<string, string[]> | null
} | null

const validateChannel = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  channelConfig: ChannelConfig,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/channels/validate`

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)

  return await apiRequest<ResponseData>(url, 'PATCH', channelConfig, {
    authorization: key,
    query: queryParams,
  })
}

export default validateChannel
