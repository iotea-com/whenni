'use server'

import gruentClient from '@gruent/hub/lib/gruent'
import { revalidatePath } from 'next/cache'
import getAccessToken from '@gruent/hub/util/getAccessToken'

export const handleAddApiKey = async (
  orgId: string,
  permissionSetId: string,
  name?: string,
  spaceId?: string,
) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to add the API key.',
    }

  const { errors } = await gruentClient(accessToken).apiKeys.add(
    orgId,
    permissionSetId,
    name,
    spaceId,
  )

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('organizations/[orgId]/team', 'page')
  return { error: null }
}

export const handleDeleteApiKey = async (orgId: string, apiKeyId: string, spaceId?: string) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to delete the API key.',
    }

  const { errors } = await gruentClient(accessToken).apiKeys.remove(orgId, apiKeyId, spaceId)

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('organizations/[orgId]/team', 'page')
  return { error: null }
}
