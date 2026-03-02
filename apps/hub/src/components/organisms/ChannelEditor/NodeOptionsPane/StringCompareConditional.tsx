import useChannelEditorStore from '@iotea/hub/stores/channelEditorStore'
import FormFieldSelect from '@iotea/libs/frontend/components/atoms/FormFieldSelect'
import { ChannelNode } from '@iotea/libs/engine/nodes/v1'
import {
  StringCompareConditionalNodeConfig,
  StringComparison,
  StringComparisonOperator,
} from '@iotea/libs/engine/nodes/v1/src/conditional/stringCompare'
import { Model } from '@prisma/client'
import {
  ChangeEvent,
  Dispatch,
  FC,
  SetStateAction,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from 'react'
import { ModelAttributes } from '@iotea/libs/engine/dependencies/models'
import Modal from '@iotea/libs/frontend/components/organisms/Modal'
import { riAddLine, riCloseLine } from '@mwarnerdotme/react-remixicon'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import { RemixIcon } from '@mwarnerdotme/react-remixicon'
import { openModal } from '@iotea/libs/frontend/hooks/useModal'

type Props = {
  models: Model[]
  orgId: string
  spaceId: string
}

const StringCompareConditionalOptions: FC<Props> = ({ models, orgId, spaceId }) => {
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<StringCompareConditionalNodeConfig>,
  )

  const [model, setModel] = useState<Model>()
  const [comparisons, setComparisons] = useState<StringComparison[]>([])

  // load in default values
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!models || models.length <= 0) return

    const model = models.find((model) => {
      if (model.id === currentConfig['model::model']) return true
      return false
    })

    setModel(model)
    setComparisons(currentConfig.comparisons || [])
  }, [currentNode, models])

  // update config and dependencies when options have been updated
  useEffect(() => {
    if (!model) return

    // update config
    const updatedNode = structuredClone(currentNode)
    updatedNode.metadata.config = {
      'model::model': model ? model.id : '',
      comparisons,
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
  }, [model, comparisons, upsertNode])

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

  return (
    <>
      <StringCompareNodeModal
        comparisons={comparisons}
        model={model}
        orgId={orgId}
        spaceId={spaceId}
        setComparisons={setComparisons}
      />
      <div className="flex flex-row gap-2">
        <FormFieldSelect
          name="model"
          label="Model"
          className="grow"
          options={modelOptions}
          value={model ? model.id : '__IOTEA_IGNORE__'}
          onChange={(e) => setModel(models.find((m) => m.id === e.target.value))}
        />
        <Button className="my-3" onClick={() => openModal('createModel')}>
          <RemixIcon icon={riAddLine} />
        </Button>
      </div>
      <Button text="Edit conditions" onClick={() => openModal('stringCompare')} />
    </>
  )
}

type StringCompareModalProps = {
  comparisons: StringComparison[]
  model?: Model
  orgId: string
  spaceId: string
  setComparisons: Dispatch<SetStateAction<StringComparison[]>>
}

const StringCompareNodeModal: FC<StringCompareModalProps> = ({
  comparisons,
  model,
  orgId,
  spaceId,
  setComparisons,
}) => {
  const modelAttributes: ModelAttributes = useMemo(() => {
    if (model && model.attributes) return JSON.parse(model.attributes as string) as ModelAttributes
    return {}
  }, [model])

  const flattenModelAttributes = useCallback(
    (
      attributes: ModelAttributes,
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

        // Only include string fields since we're doing string comparisons
        if (attribute.type === 'string') {
          return [
            {
              key: prefix ? `${prefix}.${attribute.key}` : attribute.key,
              value: attributeId,
            },
          ]
        }
        return []
      })
    },
    [],
  )

  const availableAttributes = useMemo(() => {
    if (!modelAttributes) return []

    const flatAttributes = flattenModelAttributes(modelAttributes)
    return flatAttributes.map(({ key, value }) => ({
      label: key,
      value: value,
    }))
  }, [modelAttributes, flattenModelAttributes])

  const handleAddComparison = () => {
    if (availableAttributes.length === 0) return

    const defaultComparison = {
      attributeId: availableAttributes[0].value,
      operator: StringComparisonOperator.EqualTo,
      value: '',
    }

    const updatedComparisons = [...comparisons, defaultComparison]
    setComparisons(updatedComparisons as StringComparison[])
  }

  const handleRemoveComparison = (comparisonIndex: number) => {
    const updatedComparisons = comparisons.filter((_comparison, i) => {
      if (comparisonIndex === i) return false
      return true
    })

    setComparisons(updatedComparisons)
  }

  const handleComparisonAttributeChange = (
    e: ChangeEvent<HTMLSelectElement>,
    comparisonIndex: number,
  ) => {
    const attributeId = e.target.value

    const updatedComparisons = [...comparisons]
    updatedComparisons[comparisonIndex] = {
      ...updatedComparisons[comparisonIndex],
      attributeId,
    }

    setComparisons(updatedComparisons)
  }

  const handleComparisonOperatorChange = (
    e: ChangeEvent<HTMLSelectElement>,
    comparisonIndex: number,
  ) => {
    const operator = e.target.value

    const updatedComparisons = [...comparisons]

    updatedComparisons[comparisonIndex] = {
      ...updatedComparisons[comparisonIndex],
      operator: operator as StringComparisonOperator,
    }

    setComparisons(updatedComparisons)
  }

  const handleComparisonValueChange = (
    e: ChangeEvent<HTMLInputElement>,
    comparisonIndex: number,
  ) => {
    const updatedComparisons = [...comparisons]

    updatedComparisons[comparisonIndex] = {
      ...updatedComparisons[comparisonIndex],
      value: e.target.value,
    }

    setComparisons(updatedComparisons)
  }

  if (!model) {
    return (
      <Modal id="stringCompare" showAccept={false}>
        Missing model. Make sure you set a model before trying to apply a threshold.
      </Modal>
    )
  }

  return (
    <Modal id="stringCompare" showAccept={false}>
      <div className="overflow-scroll">
        {availableAttributes.length <= 0 && (
          <p>
            No string attributes were found in the '{model.name}' model.{' '}
            <a
              className="text-blue-500 underline"
              href={`/organizations/${orgId}/spaces/${spaceId}/models/${model.id}`}
            >
              Add at least one string attribute
            </a>{' '}
            and then try again.
          </p>
        )}
        {availableAttributes.length > 0 && (
          <>
            {comparisons.map((comparison, index) => {
              // Find the display label for the current attributeId
              const attributeOption = availableAttributes.find(
                (a) => a.value === comparison.attributeId,
              )
              const attributeLabel = attributeOption?.label || comparison.attributeId

              return (
                <div
                  key={`${attributeLabel}-${index}`}
                  id={`${attributeLabel}-${index}`}
                  className="relative"
                >
                  <RemixIcon
                    icon={riCloseLine}
                    className="absolute top-1 right-0 cursor-pointer"
                    onClick={() => handleRemoveComparison(index)}
                  />
                  <h2 className="text-xs mb-1 mt-2">Condition {index + 1}</h2>
                  <div className="border border-green-500 rounded-sm px-6 py-3">
                    <FormFieldSelect
                      name={`comparison${index}.attribute`}
                      label="Attribute"
                      options={availableAttributes}
                      value={comparison.attributeId}
                      onChange={(e) => handleComparisonAttributeChange(e, index)}
                    />
                    <FormFieldSelect
                      name={`condition${index}.operator`}
                      label="Operator"
                      options={[
                        { label: 'Equals', value: StringComparisonOperator.EqualTo },
                        { label: 'Contains', value: StringComparisonOperator.Contains },
                        { label: 'Starts with', value: StringComparisonOperator.StartsWith },
                        { label: 'Ends with', value: StringComparisonOperator.EndsWith },
                      ]}
                      value={comparison.operator}
                      onChange={(e) => handleComparisonOperatorChange(e, index)}
                    />
                    <FormFieldText
                      name={`comparison${index}.value`}
                      label="Value"
                      value={comparison.value.toString()}
                      onChange={(e) => handleComparisonValueChange(e, index)}
                    />
                  </div>
                </div>
              )
            })}
          </>
        )}
      </div>
      <Button
        text="Add condition"
        onClick={handleAddComparison}
        className="mb-2 w-full"
        disabled={availableAttributes.length === 0}
      />
    </Modal>
  )
}

export default StringCompareConditionalOptions
