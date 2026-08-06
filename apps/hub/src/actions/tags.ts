'use server'

import gruentClient from '@gruent/hub/lib/gruent'
import { revalidatePath } from 'next/cache'
import getAccessToken from '@gruent/hub/util/getAccessToken'

export const handleCreateTag = async (spaceId: string, name: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to add the tag.',
    }

  const { data: tag, errors } = await gruentClient(accessToken).tags.create(spaceId, name)

  if (errors && errors.length > 0) return { data: null, error: errors[0] }

  revalidatePath('spaces/[spaceId]/things', 'page')
  revalidatePath('spaces/[spaceId]/models', 'page')
  revalidatePath('spaces/[spaceId]/channels', 'page')
  revalidatePath('spaces/[spaceId]/tags', 'page')
  return { data: tag, error: null }
}

export const handleApplyTag = async (spaceId: string, tagId: string, subjectId: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to add the tag.',
    }

  const { errors } = await gruentClient(accessToken).tags.apply(spaceId, tagId, subjectId)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('spaces/[spaceId]/things', 'page')
  revalidatePath('spaces/[spaceId]/models', 'page')
  revalidatePath('spaces/[spaceId]/channels', 'page')
  return { error: null }
}

export const handleRemoveTag = async (spaceId: string, tagId: string, subjectId: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to add the tag.',
    }

  const { errors } = await gruentClient(accessToken).tags.remove(spaceId, tagId, subjectId)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('spaces/[spaceId]/things', 'page')
  revalidatePath('spaces/[spaceId]/models', 'page')
  revalidatePath('spaces/[spaceId]/channels', 'page')
  return { error: null }
}

export const handleDeleteTag = async (spaceId: string, tagId: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to delete the tag.',
    }

  const { errors } = await gruentClient(accessToken).models.delete(spaceId, tagId)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('spaces/[spaceId]/things', 'page')
  revalidatePath('spaces/[spaceId]/models', 'page')
  revalidatePath('spaces/[spaceId]/channels', 'page')
  revalidatePath('spaces/[spaceId]/tags', 'page')
  return { error: null }
}
