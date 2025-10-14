'use client'

import { FC, useState } from 'react'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { useMutation } from '@tanstack/react-query'
import ioteaClient from '@iotea/hub/lib/iotea'
import { openModal } from '@iotea/libs/frontend/hooks/useModal'
import { PermissionSet } from '@prisma/client'
import ContextMenu from '@iotea/hub/components/atoms/ContextMenu'
import { handleRemovePermissionSet } from '@iotea/hub/actions/permissionSets'
import UpdateSpacePermissionSetModal from '@iotea/hub/components/modals/UpdateSpacePermissionSetModal'
import { riDeleteBin7Line, riEdit2Line } from '@mwarnerdotme/react-remixicon'
import useAuth from '@iotea/hub/hooks/useAuth'

type Props = {
  permissionSetId: string
  orgId: string
  spaceId: string
  className?: string
}

const SpacePermissionSetContextMenu: FC<Props> = ({ permissionSetId, spaceId, orgId }) => {
  const [currentPermissionSet, setCurrentPermissionSet] = useState<PermissionSet>()
  const [showModal, setShowModal] = useState(false)

  const { accessToken } = useAuth()

  const handleEditPermissionSet = useMutation({
    mutationKey: ['editPermissionSet', permissionSetId],
    mutationFn: async () => {
      if (!accessToken) return

      const { data: results, errors } = await ioteaClient(accessToken).permissions.list(
        spaceId ?? orgId,
      )

      if (errors && errors.length > 0) {
        addToast({
          title: 'Could not retrieve the permission sets',
          body: 'The permission sets could not be retrieved for this space. Please try again.',
          level: 'error',
        })

        throw new Error(`could not remove the permission set with ID ${permissionSetId}`)
      }

      const permissionSet = results?.find((p) => {
        if (p.id === permissionSetId) return true
        return false
      })

      setCurrentPermissionSet(permissionSet)
      openModal('updateSpacePermissionSet')

      return permissionSet
    },
  })

  const contextMenuButtons = [
    {
      icon: riEdit2Line,
      title: 'Edit',
      onClick: () => {
        handleEditPermissionSet.mutate()
        setShowModal(true)
      },
    },
    {
      icon: riDeleteBin7Line,
      title: 'Delete',
      onClick: async () => {
        const { error } = await handleRemovePermissionSet(spaceId ?? orgId, permissionSetId)
        if (error) {
          addToast({
            title: 'Could not delete the permission set',
            body: `The permission set that you selected could not be deleted: ${error}`,
            level: 'error',
          })
        } else {
          addToast({
            title: 'Successfully deleted permission set',
            body: 'Successfully deleted the permission set from your space.',
            level: 'success',
          })
        }
      },
    },
  ]

  return (
    <>
      <UpdateSpacePermissionSetModal
        orgId={orgId}
        spaceId={spaceId}
        showModal={showModal}
        initialPermissionSet={currentPermissionSet}
        onClose={() => setShowModal(false)}
      />
      <ContextMenu contextMenuButtons={contextMenuButtons} />
    </>
  )
}

export default SpacePermissionSetContextMenu
