'use client'

import { FC, useMemo, useState } from 'react'
import FormFieldTextArea from '@iotea/libs/frontend/components/atoms/FormFieldTextArea'
import { ModelAttributes } from '@iotea/libs/engine/dependencies/models'
import ModelSettingsForm from '../ModelSettingsForm'
import { handleAddModel } from '@iotea/hub/actions/models'
import { useMutation } from '@tanstack/react-query'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { useRouter } from 'next/navigation'

type Props = {
  orgId: string
  spaceId: string
}

type ModelImport = {
  name: string
  attributes: ModelAttributes
}

const ModelImportForm: FC<Props> = ({ orgId, spaceId }: Props) => {
  const router = useRouter()

  const [modelJson, setModelJson] = useState<string>()
  const [modelImportError, setModelImportError] = useState<string>()

  const modelImport = useMemo<ModelImport | null>(() => {
    try {
      if (!modelJson) return null

      const modelImport = JSON.parse(modelJson)

      if (!modelImport.attributes) {
        setModelImportError('attributes are required')
        return null
      }

      setModelImportError('')

      return modelImport as ModelImport
    } catch (_err) {
      setModelImportError('could not parse import object')
      return null
    }
  }, [modelJson])

  const { mutate: createModel } = useMutation({
    mutationFn: async ({ attributes, name }: ModelImport) => {
      const { error } = await handleAddModel(spaceId, { attributes, name })
      if (error) addToast({ title: 'Error', body: error, level: 'error' })
      else {
        addToast({
          title: 'Model imported',
          body: 'The model has been successfully imported into the space.',
          level: 'success',
        })

        router.push(`/organizations/${orgId}/spaces/${spaceId}/models`)
      }
    },
  })

  const handleSubmit = async (attributes: ModelAttributes, name: string) => {
    createModel({ attributes, name })
  }

  return (
    <>
      <section className="mb-8">
        <FormFieldTextArea
          label="Import object"
          name="modelJson"
          value={modelJson}
          onChange={(e) => setModelJson(e.target.value)}
          className="h-40"
        />
        {modelImportError && (
          <p className="text-red-400">Invalid model import object: {modelImportError}</p>
        )}
      </section>
      {modelImport && (
        <ModelSettingsForm
          spaceId={spaceId}
          initialModel={{
            name: modelImport.name,
            attributes: modelImport.attributes,
          }}
          onSubmit={handleSubmit}
        />
      )}
    </>
  )
}

export default ModelImportForm
