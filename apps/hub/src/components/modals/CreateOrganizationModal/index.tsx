'use client'

import { FC, createRef, useCallback } from 'react'
import Modal from '@gruent/libs/frontend/components/organisms/Modal'
import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import { addToast } from '@gruent/libs/frontend/hooks/useToast'
import { useRouter } from 'next/navigation'
import noop from '@gruent/libs/frontend/util/noop'
import { closeModal } from '@gruent/libs/frontend/hooks/useModal'
import { handleAddOrganization } from '@gruent/hub/actions/organizations'
import useAuth from '@gruent/hub/hooks/useAuth'

type Props = {}

const CreateOrganizationModal: FC<Props> = () => {
  const { userId } = useAuth()
  const router = useRouter()

  const formRef = createRef<HTMLFormElement>()

  const handleAccept = useCallback(async () => {
    const formData = new FormData(formRef.current as HTMLFormElement)
    const name = formData.get('name')

    if (!userId) return null
    if (!name) return

    const { error } = await handleAddOrganization(userId, name.toString())

    if (error) {
      addToast({
        title: 'Could not create a new organization',
        body: `The organization could not be created: ${error}`,
        level: 'error',
      })
      return
    }

    addToast({
      title: 'Successfully created a new organization!',
      body: `Your new organization is ready! Happy collaborating.`,
      level: 'success',
    })

    formRef.current?.reset()
    closeModal()
    router.refresh()
  }, [router, formRef, userId])

  const handleClose = useCallback(() => {
    formRef.current?.reset()
    closeModal()
  }, [formRef])

  return (
    <Modal
      id="createOrganization"
      onAccept={handleAccept}
      onClose={handleClose}
      acceptText="Create"
    >
      <div className="mb-6">
        <h2>Create a new organization</h2>
        <p className="text-gray-500">
          <small>Use an organization to collaborate with your team across spaces.</small>
        </p>
      </div>
      <form onSubmit={noop} ref={formRef} className="mb-4">
        <FormFieldText label="Name" name="name" className="w-full" autocomplete="off" />
      </form>
    </Modal>
  )
}

export default CreateOrganizationModal
