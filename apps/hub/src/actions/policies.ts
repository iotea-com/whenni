'use server'

import ioteaClient from '@iotea/hub/lib/iotea'
import { Policy } from '@iotea/libs/iotea-js/src/policies'
import { revalidatePath } from 'next/cache'
import getAccessToken from '@iotea/hub/util/getAccessToken'

export const handleUpdatePolicy = async (
  spaceId: string,
  certificateId: string,
  updatedPolicy: Policy,
  revoke: boolean,
) => {
  const accessToken = await getAccessToken()

  if (!accessToken)
    return {
      error: 'Your session is currently inactive. Sign in again to update the policy.',
    }

  const { errors: updateErrors } = await ioteaClient(accessToken).policies.update(
    spaceId,
    certificateId,
    updatedPolicy,
    revoke,
  )

  if (updateErrors && updateErrors.length > 0)
    return {
      error: `An error occurred while updating the policy: ${updateErrors[0]}`,
    }

  revalidatePath('spaces/[spaceId]/things/[thingId]', 'page')
  return { error: null }
}
