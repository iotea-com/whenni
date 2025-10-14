import type { ChannelConfig } from '@iotea/libs/engine/channels/index'
import { ClientConfig } from '..'
import { Channel } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = Channel

const updateChannelConfig = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  channelId: string,
  channelConfig: ChannelConfig,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/channels/${channelId}/config`

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)

  return await apiRequest<ResponseData>(url, 'PATCH', channelConfig, {
    authorization: key,
    query: queryParams,
  })
}

export default updateChannelConfig
