'use client'

import { handleDeleteOrganization } from '@iotea/hub/actions/organizations'
import Modal from '@iotea/libs/frontend/components/organisms/Modal'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { Organization } from '@prisma/client'
import { useRouter } from 'next/navigation'
import { FC } from 'react'

type Props = {
  organization: Organization
}

const DeleteOrganizationModal: FC<Props> = ({ organization }) => {
  const router = useRouter()

  const handleAccept = async () => {
    const { error } = await handleDeleteOrganization(organization.id)
    if (error) {
      addToast({
        title: 'Could not delete the organization',
        body: `The organization could not be deleted: ${error}`,
        level: 'error',
      })
      return
    }

    router.push('/dashboard')

    addToast({
      title: 'Successfully deleted the organization',
      body: `Your organization has been permanently deleted.`,
      level: 'success',
    })
  }

  return (
    <Modal id="deleteOrganization" onAccept={handleAccept} acceptText="Delete">
      <h2 className="text-lg">Delete organization</h2>
      <p className="mb-6">
        Once deleted, organizations cannot be recovered. Are you sure you want to delete{' '}
        {organization.name}?
      </p>
    </Modal>
  )
}

export default DeleteOrganizationModal
