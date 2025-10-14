import { ClientConfig } from '..'
import { Tag } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = Tag

const removeTagById = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  tagId: string,
  subjectId: string,
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/tags/${tagId}/remove`

  const queryParams = new URLSearchParams({ spaceId, subjectId })

  return await apiRequest<ResponseData>(url, 'PATCH', null, {
    authorization: key,
    query: queryParams,
  })
}

export default removeTagById
