import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import { ChannelNode, MetricActionNodeConfig } from '@gruent/libs/engine/nodes/v1'
import { Model, Thing } from '@prisma/client'
import { FC, useEffect, useMemo, useState } from 'react'
import { RemixIcon, riAddLine } from '@mwarnerdotme/react-remixicon'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import { ModelAttributes } from '@gruent/libs/engine/dependencies/models'

type Props = {
  models: Model[]
  things: Thing[]
  orgId: string
  spaceId: string
}

const MetricAction: FC<Props> = ({ models, things }) => {
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<MetricActionNodeConfig>,
  )

  const [selectedModel, setSelectedModel] = useState<Model>()
  const [selectedClickhouse, setSelectedClickhouse] = useState<Thing>()
  const [selectedAttribute, setSelectedAttribute] = useState<string>()
  const [selectedMetadata, setSelectedMetadata] = useState<string>()

  // Debug log for things array
  useEffect(() => {
    console.log('All things:', things)
    console.log(
      'Things categories:',
      things?.map((t) => ({ id: t.id, name: t.name, category: t.thingCategory })),
    )
  }, [things])

  // Filter for Clickhouse databases
  const clickhouseOptions = useMemo(() => {
    if (!things) return []

    const filtered = things.filter((thing) => {
      console.log('Checking thing:', {
        id: thing.id,
        name: thing.name,
        category: thing.thingCategory,
      })
      return thing.thingCategory === 'CLICKHOUSE_DATABASE'
    })

    console.log('Filtered Clickhouse things:', filtered)

    return filtered.map((thing) => ({
      label: thing.name,
      value: thing.id,
    }))
  }, [things])

  // Get numeric fields from selected model
  const numericAttributeOptions = useMemo(() => {
    if (!selectedModel?.attributes) return []

    const attributes = JSON.parse(selectedModel.attributes as string) as ModelAttributes
    console.log('Model attributes:', attributes)

    return Object.entries(attributes)
      .filter(([_, attr]) => attr.type === 'number')
      .map(([id, attr]) => ({
        label: attr.key,
        value: id,
      }))
  }, [selectedModel])

  // Get metadata fields (string, number, or boolean) from selected model
  const metadataOptions = useMemo(() => {
    if (!selectedModel?.attributes) return []

    const attributes = JSON.parse(selectedModel.attributes as string) as ModelAttributes
    return Object.entries(attributes)
      .filter(([_, attr]) => ['string', 'number', 'boolean'].includes(attr.type))
      .map(([id, attr]) => ({
        label: attr.key,
        value: id,
      }))
  }, [selectedModel])

  // load in default values
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    console.log('Current config:', currentConfig)

    if (!models || models.length <= 0) return

    const model = models.find((model) => {
      if (model.id === currentConfig.model) return true
      return false
    })

    const clickhouse = things?.find(
      (thing) => thing.id === currentConfig['clickhouseDatabase::thing'],
    )
    console.log('Found clickhouse:', clickhouse)

    setSelectedModel(model)
    setSelectedClickhouse(clickhouse)
    setSelectedAttribute(currentConfig.attribute)
    setSelectedMetadata(currentConfig.metadata)
  }, [currentNode, models, things])

  // update config and dependencies when options have been updated
  useEffect(() => {
    if (!selectedModel || !selectedClickhouse || !selectedAttribute || !selectedMetadata) return

    // update config
    const updatedNode = structuredClone(currentNode)
    updatedNode.metadata.config = {
      model: selectedModel.id,
      'clickhouseDatabase::thing': selectedClickhouse.id,
      attribute: selectedAttribute,
      metadata: selectedMetadata,
    }

    // update dependencies
    updatedNode.metadata.dependencies = {
      models: [
        {
          modelId: selectedModel.id,
          fieldName: 'model',
        },
      ],
      things: [
        {
          thingId: selectedClickhouse.id,
          fieldName: 'clickhouseDatabase::thing',
          internal: false,
        },
      ],
    }

    upsertNode(currentNode.id, updatedNode)
  }, [selectedModel, selectedClickhouse, selectedAttribute, selectedMetadata, upsertNode])

  const modelOptions = useMemo(() => {
    if (!models) return []

    return models.map((model) => ({
      label: model.name,
      value: model.id,
    }))
  }, [models])

  if (!models || models.length < 1)
    return (
      <p className="mb-2">
        <Button variant="underline" onClick={() => openModal('createModel')}>
          Add at least one model
        </Button>{' '}
        to your space before using this node.
      </p>
    )

  if (!clickhouseOptions || clickhouseOptions.length <= 0)
    return (
      <p className="mb-2">
        <Button variant="underline" onClick={() => openModal('createThing')}>
          Add a Clickhouse database
        </Button>{' '}
        to your space before using this node.
      </p>
    )

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-row gap-2">
        <FormFieldSelect
          name="model"
          label="Model"
          className="grow"
          options={modelOptions}
          value={selectedModel ? selectedModel.id : '__GRUENT_IGNORE__'}
          onChange={(e) => setSelectedModel(models.find((m) => m.id === e.target.value))}
        />
        <Button text="Create model" className="my-3" onClick={() => openModal('createModel')}>
          <RemixIcon icon={riAddLine} />
        </Button>
      </div>
      {selectedModel && numericAttributeOptions.length > 0 && (
        <FormFieldSelect
          name="attribute"
          label="Value"
          options={numericAttributeOptions}
          value={selectedAttribute || '__GRUENT_IGNORE__'}
          onChange={(e) => setSelectedAttribute(e.target.value)}
        />
      )}
      {selectedModel && numericAttributeOptions.length === 0 && (
        <p className="text-sm text-red-500">
          No numeric attributes found in this model. Add at least one numeric attribute to use it
          for metrics.
        </p>
      )}
      {selectedModel && !selectedAttribute && numericAttributeOptions.length > 0 && (
        <p className="text-sm text-red-500">Please select a value attribute for the metric.</p>
      )}
      {selectedModel && metadataOptions.length > 0 && (
        <FormFieldSelect
          name="metadata"
          label="Metadata"
          options={metadataOptions}
          value={selectedMetadata || '__GRUENT_IGNORE__'}
          onChange={(e) => setSelectedMetadata(e.target.value)}
        />
      )}
      {selectedModel && metadataOptions.length === 0 && (
        <p className="text-sm text-red-500">
          No suitable attributes found for metadata. Add at least one string, number, or boolean
          attribute.
        </p>
      )}
      {selectedModel && !selectedMetadata && metadataOptions.length > 0 && (
        <p className="text-sm text-red-500">Please select a metadata attribute.</p>
      )}
      <FormFieldSelect
        name="clickhouseDatabase"
        label="Clickhouse Database"
        options={clickhouseOptions}
        value={selectedClickhouse ? selectedClickhouse.id : '__GRUENT_IGNORE__'}
        onChange={(e) =>
          setSelectedClickhouse(
            things.find((s) => {
              if (s.id === e.target.value) return true
              return false
            }),
          )
        }
      />
      {!selectedClickhouse && (
        <p className="text-sm text-red-500">Please select a Clickhouse database.</p>
      )}
    </div>
  )
}

export default MetricAction
