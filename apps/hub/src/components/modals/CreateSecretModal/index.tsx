'use client'

import { FC, useCallback, useState } from 'react'
import Modal from '@iotea/libs/frontend/components/organisms/Modal'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { closeModal } from '@iotea/libs/frontend/hooks/useModal'
import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import { handleAddSecret } from '@iotea/hub/actions/secrets'

type Props = {
  spaceId: string
}

const CreateSecretModal: FC<Props> = ({ spaceId }) => {
  const [name, setName] = useState('')
  const [value, setValue] = useState('')

  const handleClose = useCallback(() => {
    closeModal()
  }, [])

  const handleSubmit = useCallback(async () => {
    const { error } = await handleAddSecret(spaceId, name, value)
    if (error) {
      addToast({
        title: 'Could not add secret to space',
        body: `Could not add secret to space: ${error}`,
        level: 'warning',
      })
    } else {
      handleClose()
      addToast({
        title: 'Successfully added secret to space',
        body: 'You may now use this secret within thing configurations.',
        level: 'success',
      })
    }
  }, [spaceId, name, value, handleClose])

  return (
    <Modal
      id="addSecret"
      onAccept={handleSubmit}
      onClose={handleClose}
      acceptText="Create"
      className="min-w-[400px] max-w-[80vw] overflow-y-auto"
    >
      <h2>Add a secret</h2>
      <p className="text-gray-500">
        <small>Add a secret to your space to use within thing configurations.</small>
      </p>
      <form onSubmit={handleSubmit}>
        <FormFieldText
          name="name"
          label="Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <FormFieldText
          name="value"
          label="Value"
          value={value}
          onChange={(e) => setValue(e.target.value)}
        />
      </form>
    </Modal>
  )
}

export default CreateSecretModal
