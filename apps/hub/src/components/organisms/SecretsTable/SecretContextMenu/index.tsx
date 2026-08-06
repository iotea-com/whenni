'use client'

import { FC, useState } from 'react'
import { addToast } from '@gruent/libs/frontend/hooks/useToast'
import ContextMenu from '@gruent/hub/components/atoms/ContextMenu'
import { handleDeleteSecret } from '@gruent/hub/actions/secrets'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import EditSecretModal from '@gruent/hub/components/modals/EditSecretModal'
import { riDeleteBin7Line, riPencilLine } from '@mwarnerdotme/react-remixicon'

type Props = {
  secretName: string
  spaceId: string
  className?: string
}

const SecretContextMenu: FC<Props> = ({ secretName, spaceId }) => {
  const [showModal, setShowModal] = useState<boolean>(false)

  const handleEditPermissions = () => {
    openModal('editSecret')
    setShowModal(true)
  }

  const contextMenuButtons = [
    {
      icon: riPencilLine,
      title: 'Edit',
      onClick: handleEditPermissions,
    },
    {
      icon: riDeleteBin7Line,
      title: 'Delete',
      onClick: async () => {
        const { error } = await handleDeleteSecret(spaceId, secretName)
        if (error) {
          addToast({
            title: 'Could not delete the secret',
            body: `The secret that you selected could not be deleted: ${error}`,
            level: 'error',
          })
        } else {
          addToast({
            title: 'Successfully deleted secret',
            body: 'Successfully deleted secret from your space.',
            level: 'success',
          })
        }
      },
    },
  ]

  return (
    <>
      <EditSecretModal
        spaceId={spaceId}
        name={secretName}
        showModal={showModal}
        onClose={() => setShowModal(false)}
      />
      <ContextMenu contextMenuButtons={contextMenuButtons} />
    </>
  )
}

export default SecretContextMenu
