import { ClientConfig } from '..'
import { Model } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest, { ListRequestOptions } from '../request'

type ResponseData = Model[]

const listModels = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  options?: ListRequestOptions & {
    tagFilter?: string[]
  },
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/models`

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)
  queryParams.set('page', options?.page?.toString() ?? '1')
  queryParams.set('resultsPerPage', options?.resultsPerPage?.toString() ?? '10')
  if (options?.filter) queryParams.set('q', options.filter)
  if (options?.tagFilter) queryParams.set('tagFilter', options.tagFilter.join(','))

  return await apiRequest<ResponseData>(url, 'GET', null, {
    authorization: key,
    query: queryParams,
  })
}

export default listModels
