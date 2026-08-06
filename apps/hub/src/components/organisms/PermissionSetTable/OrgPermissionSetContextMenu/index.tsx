'use client'

import { FC, useState } from 'react'
import { addToast } from '@gruent/libs/frontend/hooks/useToast'
import { useMutation } from '@tanstack/react-query'
import gruentClient from '@gruent/hub/lib/gruent'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import UpdateOrgPermissionSetModal from '@gruent/hub/components/modals/UpdateOrgPermissionSetModal'
import { PermissionSet } from '@prisma/client'
import ContextMenu from '@gruent/hub/components/atoms/ContextMenu'
import { handleRemovePermissionSet } from '@gruent/hub/actions/permissionSets'
import { riDeleteBin7Line, riEdit2Line } from '@mwarnerdotme/react-remixicon'
import useAuth from '@gruent/hub/hooks/useAuth'

type Props = {
  permissionSetId: string
  orgId: string
  className?: string
}

const OrgPermissionSetContextMenu: FC<Props> = ({ permissionSetId, orgId }) => {
  const [currentPermissionSet, setCurrentPermissionSet] = useState<PermissionSet>()
  const [showModal, setShowModal] = useState(false)

  const { accessToken } = useAuth()

  const handleEditPermissionSet = useMutation({
    mutationKey: ['editPermissionSet', permissionSetId],
    mutationFn: async () => {
      if (!accessToken) return

      const { data: results, errors } = await gruentClient(accessToken).permissions.list(orgId)

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
      openModal('updateOrgPermissionSet')

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
        const { error } = await handleRemovePermissionSet(orgId, permissionSetId)
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
      <UpdateOrgPermissionSetModal
        orgId={orgId}
        showModal={showModal}
        initialPermissionSet={currentPermissionSet}
        onClose={() => setShowModal(false)}
      />
      <ContextMenu contextMenuButtons={contextMenuButtons} />
    </>
  )
}

export default OrgPermissionSetContextMenu
