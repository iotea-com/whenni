'use client'

import { FC, useCallback } from 'react'
import Modal from '@iotea/libs/frontend/components/organisms/Modal'
import { useMutation } from '@tanstack/react-query'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { closeModal } from '@iotea/libs/frontend/hooks/useModal'
import ioteaClient from '@iotea/hub/lib/iotea'
import { useRouter } from 'next/navigation'
import useAuth from '@iotea/hub/hooks/useAuth'

type Props = {
  orgId: string
  spaceId: string
  thingId: string
}

const DeleteThingModal: FC<Props> = ({ orgId, spaceId, thingId }) => {
  const { accessToken } = useAuth()
  const router = useRouter()

  const deleteThingMutation = useMutation({
    mutationKey: ['deleteThing'],
    mutationFn: async () => {
      if (!accessToken) return

      const { errors } = await ioteaClient(accessToken).things.delete(spaceId, thingId)

      if (errors && errors.length > 0) {
        addToast({
          title: 'Could not delete thing',
          body: `The thing could not be deleted: ${errors[0]}`,
          level: 'error',
        })
        return false
      }

      closeModal()
      // TODO: revalidatePath()
      router.push(`/organizations/${orgId}/spaces/${spaceId}/things`)

      return true
    },
  })

  const handleClose = () => {
    closeModal()
  }

  const handleSubmit = useCallback(() => {
    deleteThingMutation.mutate()
  }, [deleteThingMutation])

  return (
    <Modal id="deleteThing" onAccept={handleSubmit} onClose={handleClose} acceptText="Delete">
      <h2>Delete a thing</h2>
      <p className="text-gray-500">
        <small>Delete a thing from your space.</small>
      </p>
      <div className="mb-2">
        <strong>Are you sure you want to delete this thing?</strong>
      </div>
    </Modal>
  )
}

export default DeleteThingModal
