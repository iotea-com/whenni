'use client'

import { FC, useCallback, useState } from 'react'
import Modal from '@iotea/libs/frontend/components/organisms/Modal'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { closeModal } from '@iotea/libs/frontend/hooks/useModal'
import { CreateModelInput } from '@iotea/libs/iotea-js/src/models/create'
import { handleAddModel } from '@iotea/hub/actions/models'
import { ModelAttributes } from '@iotea/libs/engine/dependencies/models'
import ModelSettingsForm from '../../organisms/ModelSettingsForm'

type Props = {
  spaceId: string
  onClose?: () => void
  onSubmit?: () => void
}

const CreateModelModal: FC<Props> = ({ spaceId, onClose, onSubmit }) => {
  const [name, setName] = useState('')
  const [attributes, setAttributes] = useState<ModelAttributes>({})

  const handleClose = useCallback(() => {
    setName('')
    setAttributes({})
    closeModal()

    if (onClose) onClose()
  }, [onClose])

  const handleSubmit = useCallback(async () => {
    const modelInput: CreateModelInput = {
      name,
      attributes,
    }

    const { error } = await handleAddModel(spaceId, modelInput)
    if (error) {
      addToast({
        title: 'Could not add model to space',
        body: `Could not add model to space: ${error}`,
        level: 'warning',
      })
    } else {
      handleClose()
      addToast({
        title: 'Successfully added model to space',
        body: 'Successfully added model to space.',
        level: 'success',
      })

      if (onSubmit) onSubmit()
    }
  }, [name, attributes, spaceId, handleClose, onSubmit])

  return (
    <Modal
      id="createModel"
      onAccept={handleSubmit}
      onClose={handleClose}
      acceptText="Create"
      className="min-w-[400px] max-w-[80vw] overflow-y-auto"
    >
      <h2>Add a model</h2>
      <p className="text-gray-500">
        <small>Add data formations to define specific inputs and outputs for nodes.</small>
      </p>
      <ModelSettingsForm
        spaceId={spaceId}
        showSaveButton={false}
        attributes={attributes}
        setAttributes={setAttributes}
        name={name}
        setName={setName}
      />
    </Modal>
  )
}

export default CreateModelModal
