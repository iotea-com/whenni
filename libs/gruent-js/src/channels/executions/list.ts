import { ClientConfig } from '../..'
import { ApiResponse } from '../../response'
import apiRequest, { ListRequestOptions } from '../../request'

export type ChannelExecutionListItem = {
  executionId: string
  startTime: string
  endTime: string
  status: string
  durationMs: number
  channelId: string
  logAttributes: Record<string, string>
}

type ResponseData = ChannelExecutionListItem[]

const listChannelExecutions = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  channelId: string,
  options?: ListRequestOptions & { statusFilter?: string },
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/channels/executions`

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)
  queryParams.set('channelId', channelId)
  queryParams.set('page', options?.page?.toString() ?? '1')
  queryParams.set('resultsPerPage', options?.resultsPerPage?.toString() ?? '10')
  if (options?.statusFilter) queryParams.set('statusFilter', options.statusFilter)

  return await apiRequest(url, 'GET', null, { authorization: key, query: queryParams })
}

export default listChannelExecutions
