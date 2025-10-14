'use server'

import ioteaClient from '@iotea/hub/lib/iotea'
import { CreateChannelInput } from '@iotea/libs/iotea-js/src/channels'
import { revalidatePath } from 'next/cache'
import getAccessToken from '@iotea/hub/util/getAccessToken'

export const handleCreateChannel = async (spaceId: string, channelInput: CreateChannelInput) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to create the channel.',
    }

  const { errors } = await ioteaClient(accessToken).channels.create(spaceId, channelInput)

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

  const { errors } = await ioteaClient(accessToken).channels.delete(spaceId, channelId)

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

  const { errors } = await ioteaClient(accessToken).channels.publish(spaceId, channelId)

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

  const { errors } = await ioteaClient(accessToken).channels.unpublish(spaceId, channelId)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('/organizations/[orgId]/spaces/[spaceId]', 'page')
  return { error: null }
}
