'use client'

import { FC, createRef, useCallback } from 'react'
import Modal from '@gruent/libs/frontend/components/organisms/Modal'
import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import { addToast } from '@gruent/libs/frontend/hooks/useToast'
import noop from '@gruent/libs/frontend/util/noop'
import { closeModal } from '@gruent/libs/frontend/hooks/useModal'
import { CreateChannelInput } from '@gruent/libs/gruent-js/src/channels/create'
import { handleCreateChannel } from '@gruent/hub/actions/channels'

type Props = {
  spaceId: string
}

const CreateChannelModal: FC<Props> = ({ spaceId }) => {
  const formRef = createRef<HTMLFormElement>()

  const handleAccept = useCallback(async () => {
    if (!formRef.current) return

    const formData = new FormData(formRef.current)
    const name = formData.get('name') as string

    const createChannelInput: CreateChannelInput = {
      name,
      edges: [],
      nodes: [],
    }

    const { error } = await handleCreateChannel(spaceId, createChannelInput)

    if (error) {
      addToast({
        title: 'Could not create a new channel',
        body: `Could not create a new channel: ${error}`,
        level: 'error',
      })
      return
    } else {
      addToast({
        title: 'Successfully created a new channel!',
        body: `Your new channel is ready! Happy connecting.`,
        level: 'success',
      })

      formRef.current.reset()
      closeModal()
    }
  }, [spaceId, formRef])

  const handleClose = useCallback(() => {
    formRef.current?.reset()
    closeModal()
  }, [formRef])

  return (
    <Modal id="createChannel" onAccept={handleAccept} onClose={handleClose} acceptText="Create">
      <div className="mb-6">
        <h2>Create a new channel</h2>
        <p className="text-gray-500">
          <small>Combine operations together to connect devices, services, and more.</small>
        </p>
      </div>
      <form onSubmit={noop} ref={formRef} className="mb-4">
        <FormFieldText label="Name" name="name" className="w-full" />
      </form>
    </Modal>
  )
}

export default CreateChannelModal
