import { ClientConfig } from '..'
import { Channel } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = Channel

const getChannel = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  channelId: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/channels/${channelId}`

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)

  return await apiRequest(url, 'GET', null, { authorization: key, query: queryParams })
}

export default getChannel
