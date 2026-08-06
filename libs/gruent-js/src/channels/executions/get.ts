import { ClientConfig } from '../..'
import { ApiResponse } from '../../response'
import apiRequest from '../../request'

type LogAttributes = {
  execution_id: string
  node: string
  level: string
  time: string
  data?: string
  'base64-data'?: string
}

type NodeExecutionLog = {
  traceId: string
  spanId: string
  timestamp: string
  level: string
  body: string
  logAttributes: LogAttributes
  resourceAttributes: Record<string, string>
}

export type ChannelExecution = {
  executionId: string
  startTime: string
  endTime: string
  durationMs: number
  channelId: string
  status: string
  nodeExecutionLogs: NodeExecutionLog[]
}

type ResponseData = ChannelExecution

const getChannelExecution = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  channelExecutionId: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/channels/executions/${channelExecutionId}`

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)
  return await apiRequest(url, 'GET', null, { authorization: key, query: queryParams })
}

export default getChannelExecution
