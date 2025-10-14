'use client'

import useAuth from '@iotea/hub/hooks/useAuth'
import ioteaClient from '@iotea/hub/lib/iotea'
import Modal from '@iotea/libs/frontend/components/organisms/Modal'
import { closeModal } from '@iotea/libs/frontend/hooks/useModal'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { Model } from '@prisma/client'
import { useMutation } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import { FC } from 'react'

type Props = {
  orgId: string
  spaceId: string
  model: Model
}

const DeleteModelModal: FC<Props> = ({ orgId, spaceId, model }) => {
  const { accessToken } = useAuth()
  const router = useRouter()

  const deleteModelMutation = useMutation({
    mutationKey: ['deleteModel', spaceId, model],
    mutationFn: async () => {
      if (!accessToken) return

      const { errors } = await ioteaClient(accessToken).models.delete(spaceId, model.id)

      if (errors && errors.length > 0) {
        addToast({
          title: 'Could not delete model',
          body: errors[0],
          level: 'error',
        })

        return errors[0]
      }

      closeModal()
      // TODO: revalidatePath()
      router.push(`/organizations/${orgId}/spaces/${spaceId}/models`)
    },
  })

  const handleAccept = () => {
    deleteModelMutation.mutate()
  }

  return (
    <Modal id="deleteModel" onAccept={handleAccept} acceptText="Delete">
      <h2 className="text-lg">Delete model</h2>
      <p className="mb-6">
        Once deleted, spaces cannot be recovered. Are you sure you want to delete {model.name}?
      </p>
    </Modal>
  )
}

export default DeleteModelModal
