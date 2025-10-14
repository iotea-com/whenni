'use client'

import { FC, useCallback, useMemo, useState } from 'react'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { openModal } from '@iotea/libs/frontend/hooks/useModal'
import { PermissionSet } from '@prisma/client'
import ContextMenu from '@iotea/hub/components/atoms/ContextMenu'
import { handleChangeMemberRole, handleRemoveMember } from '@iotea/hub/actions/members'
import EditUserPermissionsModal from '@iotea/hub/components/modals/EditUserPermissionsModal'
import { riAdminLine, riDeleteBin7Line, riEdit2Line } from '@mwarnerdotme/react-remixicon'

type Props = {
  userId: string
  orgId: string
  admin: boolean
  initialPermissionSetId: string
  permissionSets: PermissionSet[]
  className?: string
}

const MemberContextMenu: FC<Props> = ({
  userId,
  orgId,
  admin = false,
  initialPermissionSetId,
  permissionSets,
}) => {
  const [showModal, setShowModal] = useState<boolean>(false)

  const handleEditPermissions = () => {
    openModal('editUserPermissionSet')
    setShowModal(true)
  }

  const handleChangeRole = useCallback(
    async (newRole: 'ADMIN' | 'MEMBER') => {
      const { error } = await handleChangeMemberRole(orgId, userId, newRole)
      if (error) {
        addToast({
          title: 'Could not change the member role',
          body: `The member role could not be changed: ${error}`,
          level: 'error',
        })
      } else {
        addToast({
          title: 'Successfully changed the member role',
          body: `The user is now a ${newRole === 'ADMIN' ? 'admin' : 'member'} in this organization. This change should be reflected immediately.`,
          level: 'success',
        })
      }
    },
    [orgId, userId],
  )

  const contextMenuButtons = useMemo(() => {
    const adminCmb = [
      {
        icon: riAdminLine,
        title: 'Change role',
        description: 'Change role to member',
        onClick: () => {
          handleChangeRole('MEMBER')
        },
      },
    ]

    const memberCmb = [
      {
        icon: riAdminLine,
        title: 'Change role',
        description: 'Change role to admin',
        onClick: () => {
          handleChangeRole('ADMIN')
        },
      },
      {
        icon: riEdit2Line,
        title: 'Edit',
        onClick: () => handleEditPermissions(),
      },
      {
        icon: riDeleteBin7Line,
        title: 'Remove',
        description: "Revoke this member's access to this organization.",
        onClick: async () => {
          const { error } = await handleRemoveMember(orgId, userId)
          if (error) {
            addToast({
              title: 'Could not remove the team member',
              body: `The team member that you selected could not be removed: ${error}`,
              level: 'error',
            })
          } else {
            addToast({
              title: 'Successfully removed team member',
              body: 'Successfully removed the team member from your organization.',
              level: 'success',
            })
          }
        },
      },
    ]

    if (admin) return adminCmb
    return memberCmb
  }, [admin, orgId, userId, handleChangeRole])

  return (
    <>
      <EditUserPermissionsModal
        orgId={orgId}
        userId={userId}
        showModal={showModal}
        initialPermissionSetId={initialPermissionSetId}
        permissionSets={permissionSets}
        onClose={() => setShowModal(false)}
      />
      <ContextMenu contextMenuButtons={contextMenuButtons} />
    </>
  )
}

export default MemberContextMenu
