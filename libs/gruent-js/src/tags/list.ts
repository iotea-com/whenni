import { ClientConfig } from '..'
import { Tag } from '@prisma/client'
import { ApiResponse } from '../response'
import apiRequest from '../request'

type ResponseData = {
  [key: string]: {
    tag: Tag
    things: string[]
    channels: string[]
    models: string[]
  }
}

const listTags = async (
  key: string,
  config: ClientConfig,
  spaceId: string,
  category?: 'things' | 'channels' | 'models',
): Promise<ApiResponse<ResponseData>> => {
  const url = new URL(config.url)
  url.pathname = `v1/tags`

  const queryParams = new URLSearchParams({ spaceId })
  if (category) queryParams.set('category', category)

  return await apiRequest<ResponseData>(url, 'GET', null, {
    authorization: key,
    query: queryParams,
  })
}

export default listTags
