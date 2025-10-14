'use client'

import { FC, createRef, useCallback } from 'react'
import Modal from '@iotea/libs/frontend/components/organisms/Modal'
import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { useRouter } from 'next/navigation'
import noop from '@iotea/libs/frontend/util/noop'
import { closeModal } from '@iotea/libs/frontend/hooks/useModal'
import { handleAddSpace } from '@iotea/hub/actions/spaces'
import useAuth from '@iotea/hub/hooks/useAuth'

type Props = {
  orgId: string
}

const CreateSpaceModal: FC<Props> = ({ orgId }) => {
  const { userId } = useAuth()
  const router = useRouter()

  const formRef = createRef<HTMLFormElement>()

  const handleAccept = useCallback(async () => {
    const formData = new FormData(formRef.current as HTMLFormElement)
    const name = formData.get('name')

    if (!userId) return null
    if (!name) return

    const { error } = await handleAddSpace(orgId, name.toString())

    if (error) {
      addToast({
        title: 'Could not create a new space',
        body: `The space could not be created: ${error}`,
        level: 'error',
      })
      return
    }

    addToast({
      title: 'Successfully created a new space!',
      body: `Your new space is ready! Happy connecting.`,
      level: 'success',
    })

    formRef.current?.reset()
    closeModal()
    router.refresh()
  }, [router, formRef, userId, orgId])

  const handleClose = useCallback(() => {
    formRef.current?.reset()
    closeModal()
  }, [formRef])

  return (
    <Modal id="createSpace" onAccept={handleAccept} onClose={handleClose} acceptText="Create">
      <div className="mb-6">
        <h2>Create a new space</h2>
        <p className="text-gray-500">
          <small>Use a space to organize and manage channels, things, and data.</small>
        </p>
      </div>
      <form onSubmit={noop} ref={formRef} className="mb-4">
        <FormFieldText label="Name" name="name" className="w-full" autocomplete="off" />
      </form>
    </Modal>
  )
}

export default CreateSpaceModal
