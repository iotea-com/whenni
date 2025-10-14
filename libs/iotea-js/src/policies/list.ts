import { Certificate } from '@prisma/client'
import { ClientConfig } from '..'
import { ApiResponse } from '../response'
import apiRequest, { ListRequestOptions } from '../request'

type ResponseData = Certificate[]

const listPolicies = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  options?: ListRequestOptions,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/policies`

  const queryParams = new URLSearchParams()
  queryParams.set('spaceId', spaceId)
  queryParams.set('page', options?.page?.toString() ?? '1')
  queryParams.set('resultsPerPage', options?.resultsPerPage?.toString() ?? '10')

  return await apiRequest<ResponseData>(url, 'GET', null, {
    authorization: key,
    query: queryParams,
  })
}

export default listPolicies
