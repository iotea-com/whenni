'use server'

import ioteaClient from '@iotea/hub/lib/iotea'
import { revalidatePath } from 'next/cache'
import getAccessToken from '@iotea/hub/util/getAccessToken'

export const handleAddPermissionSet = async (
  orgId: string,
  name: string,
  permissions: string[],
  spaceId?: string,
) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to add the permission set.',
    }

  const { errors } = await ioteaClient(accessToken).permissions.add(
    orgId,
    name,
    permissions,
    spaceId,
  )

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('spaces/[spaceId]/team', 'page')
  return { error: null }
}

export const handleRemovePermissionSet = async (
  orgId: string,
  permissionSetId: string,
  spaceId?: string,
) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to remove the permission set.',
    }

  const { errors } = await ioteaClient(accessToken).permissions.remove(
    orgId,
    permissionSetId,
    spaceId,
  )

  if (errors && errors.length > 0) return { error: errors[0] }

  revalidatePath('spaces/[spaceId]/team', 'page')
  return { error: null }
}
