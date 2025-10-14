import { ClientConfig } from '..'
import { Channel } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest, { ListRequestOptions } from '../request'

type ResponseData = Channel[]

const listChannels = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  options?: ListRequestOptions & { tagFilter?: string[] },
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/channels`

  const queryParams = new URLSearchParams()
  queryParams.set('page', options?.page?.toString() ?? '1')
  queryParams.set('resultsPerPage', options?.resultsPerPage?.toString() ?? '10')
  queryParams.set('spaceId', spaceId)
  if (options?.filter) queryParams.set('q', options.filter)
  if (options?.tagFilter) queryParams.set('tagFilter', options.tagFilter.join(','))

  return await apiRequest(url, 'GET', null, { authorization: key, query: queryParams })
}

export default listChannels
