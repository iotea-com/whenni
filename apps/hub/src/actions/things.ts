'use server'

import gruentClient from '@gruent/hub/lib/gruent'
import { Thing } from '@prisma/client'
import { revalidatePath } from 'next/cache'
import getAccessToken from '@gruent/hub/util/getAccessToken'

export const handleAddThing = async (
  spaceId: string,
  name: string,
  category: string,
  attributes: Object,
) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to add the thing.',
    }

  const { errors } = await gruentClient(accessToken).things.create(
    spaceId,
    name,
    category,
    attributes,
  )

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('spaces/[spaceId]/things', 'page')
  return { error: null }
}

export const handleUpdateThing = async (
  spaceId: string,
  thingId: string,
  thingInput: Pick<Thing, 'name' | 'attributes'>,
) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to add the thing.',
    }

  const { errors } = await gruentClient(accessToken).things.update(spaceId, thingId, thingInput)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('spaces/[spaceId]/things/[thingId]', 'page')
  return { error: null }
}

export const handleDeleteThing = async (spaceId: string, thingId: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to delete the thing.',
    }

  const { errors } = await gruentClient(accessToken).things.delete(spaceId, thingId)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('spaces/[spaceId]/things', 'page')
  return { error: null }
}

export const handleHealthcheckThing = async (
  spaceId: string,
  thingCategory: string,
  attributes: Record<string, any>,
) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to delete the thing.',
    }

  const { errors } = await gruentClient(accessToken).things.healthcheck(
    spaceId,
    thingCategory,
    attributes,
  )

  if (errors && errors.length > 0) return { error: errors[0] }

  return { error: null }
}
