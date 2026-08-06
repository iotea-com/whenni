import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import { ChannelNode, ExistenceConditionalNodeConfig } from '@gruent/libs/engine/nodes/v1'
import { Model } from '@prisma/client'
import { FC, useCallback, useEffect, useMemo, useState } from 'react'
import { ModelAttribute, ModelAttributes } from '@gruent/libs/engine/dependencies/models'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import { RemixIcon, riAddLine } from '@mwarnerdotme/react-remixicon'

type Props = {
  models: Model[]
}

const ExistenceConditionalOptions: FC<Props> = ({ models }) => {
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<ExistenceConditionalNodeConfig>,
  )

  const [model, setModel] = useState<Model>()
  const [attributeIds, setAttributeIds] = useState<string[]>([])

  // load in default values
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!models || models.length <= 0) return

    const model = models.find((model) => {
      if (model.id === currentConfig['model::model']) return true
      return false
    })

    setModel(model)
    setAttributeIds(currentConfig.attributeIds || [])
  }, [currentNode, models])

  // update config and dependencies when options have been updated
  useEffect(() => {
    if (!model) return

    // update config
    const updatedNode = structuredClone(currentNode)
    updatedNode.metadata.config = {
      'model::model': model ? model.id : '',
      attributeIds,
    }

    // update dependencies
    updatedNode.metadata.dependencies.models = [
      {
        modelId: model.id,
        internal: false,
        fieldName: 'model::model',
      },
    ]

    upsertNode(currentNode.id, updatedNode)

    // eslint-disable-next-line -- currentNode in the dependencies array causes node configs to sometimes be copied
  }, [model, attributeIds, upsertNode])

  const modelOptions = useMemo(() => {
    if (!models) return []

    return models.map((model) => ({
      label: model.name,
      value: model.id,
    }))
  }, [models])

  const modelSchema: ModelAttributes = useMemo(() => {
    if (model && model.attributes) return JSON.parse(model.attributes as string) as ModelAttributes
    return {}
  }, [model])

  const flattenModelAttributes = useCallback(
    (
      attributes: Record<string, ModelAttribute>,
      parentId = '',
      prefix = '',
    ): Array<{ key: string; value: string }> => {
      if (!attributes || Object.keys(attributes).length <= 0) return []

      const attributesInObject = Object.entries(attributes).filter(([_attributeId, attribute]) => {
        if (parentId) return attribute.parentId === parentId
        return !attribute.parentId
      })

      return attributesInObject.flatMap(([attributeId, attribute]) => {
        if (attribute.type === 'object') {
          const newPrefix = prefix ? `${prefix}.${attribute.key}` : attribute.key
          return flattenModelAttributes(attributes, attributeId, newPrefix)
        }

        // Include only optional fields
        return [
          {
            key: prefix ? `${prefix}.${attribute.key}` : attribute.key,
            value: attributeId,
          },
        ]
      })
    },
    [],
  )

  const availableAttributes = useMemo(() => {
    if (!modelSchema) return []

    const flatAttributes = flattenModelAttributes(modelSchema)
    return flatAttributes.map(({ key, value }) => ({
      label: key,
      value: value,
    }))
  }, [modelSchema, flattenModelAttributes])

  if (!models || models.length < 1)
    return (
      <p className="mb-2">
        <Button variant="underline" onClick={() => openModal('createModel')}>
          Add at least one model
        </Button>{' '}
        to your space before using this node.
      </p>
    )

  return (
    <>
      <div className="flex flex-row gap-2">
        <FormFieldSelect
          name="model"
          label="Model"
          className="grow"
          options={modelOptions}
          value={model ? model.id : '__GRUENT_IGNORE__'}
          onChange={(e) => setModel(models.find((m) => m.id === e.target.value))}
        />
        <Button className="my-3" onClick={() => openModal('createModel')}>
          <RemixIcon icon={riAddLine} />
        </Button>
      </div>
      {model && (
        <FormFieldSelect
          name="modelAttributes"
          label="Attributes"
          options={availableAttributes}
          value={attributeIds}
          multiple
          onChange={(e) => {
            const selectedOptions = Array.from(e.target.selectedOptions, (option) => option.value)
            setAttributeIds(selectedOptions.filter((id) => id !== '__GRUENT_IGNORE__'))
          }}
        />
      )}
    </>
  )
}

export default ExistenceConditionalOptions
