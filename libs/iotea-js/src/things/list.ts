import { ClientConfig } from '..'
import { AppliedTag, Tag, Thing } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest, { ListRequestOptions } from '../request'

type ResponseData = (Thing & {
  tags: (AppliedTag & { tag: Tag })[]
})[]

const listThings = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  options?: ListRequestOptions & {
    category?: string
    tagFilter?: string[]
  },
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/things`

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)
  queryParams.set('page', options?.page?.toString() ?? '1')
  queryParams.set('resultsPerPage', options?.resultsPerPage?.toString() ?? '10')
  if (options && options.category) queryParams.set('category', options.category)
  if (options && options.filter) queryParams.set('q', options.filter)
  if (options && options.tagFilter) queryParams.set('tagFilter', options.tagFilter.join(','))

  return await apiRequest<ResponseData>(url, 'GET', null, {
    authorization: key,
    query: queryParams,
  })
}

export default listThings
