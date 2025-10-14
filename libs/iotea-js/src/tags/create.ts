import { ClientConfig } from '..'
import { ApiResponse } from '../response'
import apiRequest from '../request'
import { Tag } from '@prisma/client'

type ResponseData = {
  tag: Tag
}

const createTag = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  name: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/tags`

  const requestBody = { name }

  const queryParams = new URLSearchParams({ spaceId })

  return await apiRequest<ResponseData>(url, 'POST', requestBody, {
    authorization: key,
    query: queryParams,
  })
}

export default createTag
