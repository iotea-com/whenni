import { ClientConfig } from '..'
import { ApiKey } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest, { ListRequestOptions } from '../request'

type ResponseData = ApiKey[]

const listApiKeys = async (
  key: string,
  config: ClientConfig,
  orgId: string,
  options?: ListRequestOptions,
  spaceId?: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/api-keys`
  url.searchParams.set('page', options?.page?.toString() ?? '1')
  url.searchParams.set('resultsPerPage', options?.resultsPerPage?.toString() ?? '10')

  const queryParams = new URLSearchParams({ orgId })
  if (spaceId) queryParams.set('spaceId', spaceId)

  return await apiRequest(url, 'GET', null, { authorization: key, query: queryParams })
}

export default listApiKeys
