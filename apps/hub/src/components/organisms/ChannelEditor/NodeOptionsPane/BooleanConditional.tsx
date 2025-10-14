import useChannelEditorStore from '@iotea/hub/stores/channelEditorStore'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import FormFieldSelect from '@iotea/libs/frontend/components/atoms/FormFieldSelect'
import Modal from '@iotea/libs/frontend/components/organisms/Modal'
import { openModal } from '@iotea/libs/frontend/hooks/useModal'
import { ChannelNode, BooleanConditionalNodeConfig } from '@iotea/libs/engine/nodes/v1'
import { BooleanCondition } from '@iotea/libs/engine/nodes/v1/src/conditional/boolean'
import { RemixIcon, riAddLine, riCloseLine } from '@mwarnerdotme/react-remixicon'
import {
  ChangeEvent,
  Dispatch,
  FC,
  SetStateAction,
  useEffect,
  useMemo,
  useState,
  useCallback,
} from 'react'
import { ModelAttributes } from '@iotea/libs/engine/dependencies/models'
import { Model } from '@prisma/client'

type Props = {
  models: Model[]
  orgId: string
  spaceId: string
}

const BooleanConditionalOptions: FC<Props> = ({ models, orgId, spaceId }) => {
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<BooleanConditionalNodeConfig>,
  )

  const [model, setModel] = useState<Model>()
  const [conditions, setConditions] = useState<BooleanCondition[]>([])
  const [logicalOperator, setLogicalOperator] = useState<'AND' | 'OR'>('OR')

  // load in default values
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!models || models.length <= 0) return

    const model = models.find((model) => {
      if (model.id === currentConfig['model::model']) return true
      return false
    })

    setModel(model)
    setConditions(currentConfig.conditions || [])
    setLogicalOperator(currentConfig.logicalOperator || 'OR')
  }, [currentNode, models])

  // update config and dependencies when options have been updated
  useEffect(() => {
    if (!model) return

    // update config
    const updatedNode = structuredClone(currentNode)
    updatedNode.metadata.config = {
      'model::model': model ? model.id : '',
      conditions,
      logicalOperator,
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
  }, [model, conditions, logicalOperator, upsertNode])

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
      <BooleanNodeModal
        model={model}
        conditions={conditions}
        orgId={orgId}
        spaceId={spaceId}
        setConditions={setConditions}
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
      <FormFieldSelect
        name="logicalOperator"
        label="Logical Operator"
        options={[
          {
            label: 'OR',
            value: 'OR',
          },
          {
            label: 'AND',
            value: 'AND',
          },
        ]}
        value={logicalOperator}
        onChange={(e) => setLogicalOperator(e.target.value as 'AND' | 'OR')}
      />
      <Button
        text="Edit conditions"
        className="w-full mb-2"
        onClick={() => openModal('booleanConditions')}
      />
    </>
  )
}

type BooleanModalProps = {
  conditions: BooleanCondition[]
  model?: Model
  orgId: string
  spaceId: string
  setConditions: Dispatch<SetStateAction<BooleanCondition[]>>
}

const BooleanNodeModal: FC<BooleanModalProps> = ({
  conditions,
  model,
  orgId,
  spaceId,
  setConditions,
}) => {
  const modelAttributes: ModelAttributes = useMemo(() => {
    if (model && model.attributes) return JSON.parse(model.attributes as string) as ModelAttributes
    return {}
  }, [model])

  const flattenSchemaFields = useCallback(
    (
      attributes: ModelAttributes,
      parentId = '',
      prefix = '',
    ): Array<{ key: string; value: string }> => {
      if (!attributes || Object.keys(attributes).length <= 0) return []

      const fieldsInObject = Object.entries(attributes).filter(([_attributeId, attribute]) => {
        if (parentId) return attribute.parentId === parentId
        return !attribute.parentId
      })

      return fieldsInObject.flatMap(([attributeId, attribute]) => {
        if (attribute.type === 'object') {
          const newPrefix = prefix ? `${prefix}.${attribute.key}` : attribute.key
          return flattenSchemaFields(attributes, attributeId, newPrefix)
        }

        // Only include boolean fields since we're doing boolean comparisons
        if (attribute.type === 'boolean') {
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

    const flatAttributes = flattenSchemaFields(modelAttributes)
    return flatAttributes.map(({ key, value }) => ({
      label: key,
      value: value,
    }))
  }, [modelAttributes, flattenSchemaFields])

  const handleAddCondition = () => {
    if (availableAttributes.length === 0) return

    const defaultCondition = {
      attributeId: availableAttributes[0].value,
      operator: 'equal-to',
      value: 0,
    }

    const updatedConditions = [...conditions, defaultCondition]
    setConditions(updatedConditions as BooleanCondition[])
  }

  const handleRemoveCondition = (conditionIndex: number) => {
    const updatedConditions = conditions.filter((_condition, i) => {
      if (conditionIndex === i) return false
      return true
    })

    setConditions(updatedConditions as BooleanCondition[])
  }

  const handleConditionFieldChange = (
    e: ChangeEvent<HTMLSelectElement>,
    conditionIndex: number,
  ) => {
    const attributeId = e.target.value

    const updatedConditions = [...conditions]
    updatedConditions[conditionIndex] = {
      ...updatedConditions[conditionIndex],
      attributeId,
    }

    setConditions(updatedConditions)
  }

  const handleConditionOperatorChange = (
    e: ChangeEvent<HTMLSelectElement>,
    conditionIndex: number,
  ) => {
    const operator = e.target.value

    const updatedConditions = [...conditions]

    updatedConditions[conditionIndex] = {
      ...updatedConditions[conditionIndex],
      operator: operator as 'equal-to',
    }

    setConditions(updatedConditions)
  }

  const handleConditionValueChange = (
    e: ChangeEvent<HTMLSelectElement>,
    conditionIndex: number,
  ) => {
    const value = e.target.value === 'true'

    const updatedConditions = [...conditions]

    updatedConditions[conditionIndex] = {
      ...updatedConditions[conditionIndex],
      value,
    }

    setConditions(updatedConditions)
  }

  if (!model) {
    return (
      <Modal id="booleanConditions" showAccept={false}>
        Missing model. Make sure you set a model before trying to apply a threshold.
      </Modal>
    )
  }

  return (
    <Modal id="booleanConditions" showAccept={false}>
      <div className="overflow-scroll">
        {availableAttributes.length <= 0 && (
          <p>
            No boolean attributes were found in the '{model.name}' model.{' '}
            <a
              className="text-blue-500 underline"
              href={`/organizations/${orgId}/spaces/${spaceId}/models/${model.id}`}
            >
              Add at least one boolean field
            </a>{' '}
            and then try again.
          </p>
        )}
        {availableAttributes.length > 0 && (
          <>
            {conditions.map((condition, index) => {
              // Find the display label for the current attributeId
              const attributeOption = availableAttributes.find(
                (a) => a.value === condition.attributeId,
              )
              const attributeLabel = attributeOption?.label || condition.attributeId

              return (
                <div
                  key={`${attributeLabel}-${index}`}
                  id={`${attributeLabel}-${index}`}
                  className="relative"
                >
                  <RemixIcon
                    icon={riCloseLine}
                    className="absolute top-1 right-0 cursor-pointer"
                    onClick={() => handleRemoveCondition(index)}
                  />
                  <h2 className="text-xs mb-1 mt-2">Condition {index + 1}</h2>
                  <div className="border border-green-500 rounded px-6 py-3">
                    <FormFieldSelect
                      name={`condition${index}.field`}
                      label="Field"
                      options={availableAttributes}
                      value={condition.attributeId}
                      onChange={(e) => handleConditionFieldChange(e, index)}
                    />
                    <FormFieldSelect
                      name={`condition${index}.operator`}
                      label="Operator"
                      options={[{ label: '==', value: 'equal-to' }]}
                      value={condition.operator}
                      onChange={(e) => handleConditionOperatorChange(e, index)}
                    />
                    <FormFieldSelect
                      name={`condition${index}.value`}
                      label="Value"
                      options={[
                        { label: 'True', value: 'true' },
                        { label: 'False', value: 'false' },
                      ]}
                      value={condition.value.toString()}
                      onChange={(e) => handleConditionValueChange(e, index)}
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
        onClick={handleAddCondition}
        className="mb-2 w-full"
        disabled={availableAttributes.length === 0}
      />
    </Modal>
  )
}

export default BooleanConditionalOptions
