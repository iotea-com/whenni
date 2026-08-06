import { ClientConfig } from '..'
import { Channel } from '@prisma/client'
import type { ChannelConfig } from '@gruent/libs/engine/channels/index'
import { ApiResponse } from '../response'
import apiRequest from '../request'

export type CreateChannelInput = Pick<ChannelConfig, 'name' | 'nodes' | 'edges'>

type ResponseData = Channel

const createChannel = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  channelInput: CreateChannelInput,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/channels`

  const requestBody = {
    name: channelInput.name,
    config: {
      nodes: channelInput.nodes,
      edges: channelInput.edges,
      runtime: {
        size: 'small',
      },
    },
  }

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)

  return await apiRequest(url, 'POST', requestBody, { authorization: key, query: queryParams })
}

export default createChannel
