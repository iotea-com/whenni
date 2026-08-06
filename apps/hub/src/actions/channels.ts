'use server'

import gruentClient from '@gruent/hub/lib/gruent'
import { CreateChannelInput } from '@gruent/libs/gruent-js/src/channels'
import { revalidatePath } from 'next/cache'
import getAccessToken from '@gruent/hub/util/getAccessToken'

export const handleCreateChannel = async (spaceId: string, channelInput: CreateChannelInput) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to create the channel.',
    }

  const { errors } = await gruentClient(accessToken).channels.create(spaceId, channelInput)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('/organizations/[orgId]/spaces/[spaceId]', 'page')
  return { error: null }
}

export const handleDeleteChannel = async (spaceId: string, channelId: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to delete the channel.',
    }

  const { errors } = await gruentClient(accessToken).channels.delete(spaceId, channelId)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('/organizations/[orgId]/spaces/[spaceId]', 'page')
  return { error: null }
}

export const handlePublishChannel = async (spaceId: string, channelId: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to delete the channel.',
    }

  const { errors } = await gruentClient(accessToken).channels.publish(spaceId, channelId)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('/organizations/[orgId]/spaces/[spaceId]', 'page')
  return { error: null }
}

export const handleUnpublishChannel = async (spaceId: string, channelId: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to delete the channel.',
    }

  const { errors } = await gruentClient(accessToken).channels.unpublish(spaceId, channelId)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('/organizations/[orgId]/spaces/[spaceId]', 'page')
  return { error: null }
}
