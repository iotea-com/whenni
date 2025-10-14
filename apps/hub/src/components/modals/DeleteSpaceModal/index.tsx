'use client'

import { handleDeleteSpace } from '@iotea/hub/actions/spaces'
import Modal from '@iotea/libs/frontend/components/organisms/Modal'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { Space } from '@prisma/client'
import { useRouter } from 'next/navigation'
import { FC } from 'react'

type Props = {
  space: Space
}

const DeleteSpaceModal: FC<Props> = ({ space }) => {
  const router = useRouter()

  const handleAccept = async () => {
    const { error } = await handleDeleteSpace(space.organizationId, space.id)
    if (error) {
      addToast({
        title: 'Could not delete the space',
        body: `The space could not be deleted: ${error}`,
        level: 'error',
      })
      return
    }

    router.push('/dashboard')

    addToast({
      title: 'Successfully deleted the space',
      body: `Your space has been permanently deleted.`,
      level: 'success',
    })
  }

  return (
    <Modal id="deleteSpace" onAccept={handleAccept} acceptText="Delete">
      <h2 className="text-lg">Delete space</h2>
      <p className="mb-6">
        Once deleted, spaces cannot be recovered. Are you sure you want to delete {space.name}?
      </p>
    </Modal>
  )
}

export default DeleteSpaceModal
