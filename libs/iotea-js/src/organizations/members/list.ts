import { ClientConfig } from '../..'
import { OrganizationMember, User } from '@prisma/client'
import { ApiResponse } from '../../response'
import apiRequest, { ListRequestOptions } from '../../request'

type ResponseData = (OrganizationMember & {
  user: User
})[]

const listMembers = async (
  key: string,
  config: ClientConfig,
  orgId: string,
  options?: ListRequestOptions,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/organizations/members`
  url.searchParams.set('page', options?.page?.toString() ?? '1')
  url.searchParams.set('resultsPerPage', options?.resultsPerPage?.toString() ?? '10')
  if (options?.filter) url.searchParams.set('q', options.filter)

  const queryParams = new URLSearchParams({ orgId })

  return await apiRequest<ResponseData>(url, 'GET', null, {
    authorization: key,
    query: queryParams,
  })
}

export default listMembers
