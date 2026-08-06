import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import Modal from '@gruent/libs/frontend/components/organisms/Modal'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import { ChannelNode, TransformNodeConfig } from '@gruent/libs/engine/nodes/v1'
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
import { ModelAttributes } from '@gruent/libs/engine/dependencies/models'
import { RemixIcon, riAddLine, riArrowLeftLine } from '@mwarnerdotme/react-remixicon'

type Props = {
  models: Model[]
  orgId: string
  spaceId: string
}

const TransformProcessingOptions: FC<Props> = ({ models, orgId, spaceId }) => {
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<TransformNodeConfig>,
  )

  const [inputModel, setInputModel] = useState<Model>()
  const [outputModel, setOutputModel] = useState<Model>()
  const [transformation, setTransformation] = useState<Map<string, string>>(new Map())

  // load in default values
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!models || models.length <= 0) return

    const im = models.find((model) => {
      if (model.id === currentConfig['inputModel::model']) return true
      return false
    })

    const om = models.find((model) => {
      if (model.id === currentConfig['outputModel::model']) return true
      return false
    })

    setInputModel(im)
    setOutputModel(om)
    setTransformation(new Map<string, string>(Object.entries(currentConfig.mapping || {})))
  }, [currentNode, models])

  // update config and dependencies when options have been updated
  useEffect(() => {
    if (!inputModel) return
    if (!outputModel) return

    // update config
    const updatedNode = structuredClone(currentNode)
    updatedNode.metadata.config = {
      'inputModel::model': inputModel.id,
      'outputModel::model': outputModel.id,
      mapping: Object.fromEntries(transformation.entries()),
    }

    // update dependencies
    updatedNode.metadata.dependencies.models = [
      {
        modelId: inputModel.id,
        fieldName: 'inputModel::model',
      },
      {
        modelId: outputModel.id,
        fieldName: 'outputModel::model',
      },
    ]

    upsertNode(currentNode.id, updatedNode)

    // eslint-disable-next-line -- currentNode in the dependencies array causes node configs to sometimes be copied
  }, [inputModel, outputModel, transformation, upsertNode])

  // reset the transformation when the models change
  useEffect(() => {
    if (!inputModel || !outputModel) return

    if (
      inputModel.id === currentNode.metadata.config['inputModel::model'] &&
      outputModel.id === currentNode.metadata.config['outputModel::model']
    )
      return

    setTransformation(new Map())
  }, [inputModel, outputModel, currentNode.metadata.config])

  const modelOptions = useMemo(() => {
    if (!models) return []

    return models.map((model) => ({
      label: model.name,
      value: model.id,
    }))
  }, [models])

  if (!models || models.length < 2)
    return (
      <p className="mb-2">
        <Button variant="underline" onClick={() => openModal('createModel')}>
          Add at least two models
        </Button>{' '}
        to your space before using this node.
      </p>
    )

  return (
    <>
      <TransformNodeModal
        transformation={transformation}
        currentNode={currentNode}
        inputModel={inputModel}
        outputModel={outputModel}
        orgId={orgId}
        spaceId={spaceId}
        setTransformation={setTransformation}
      />
      <div className="flex flex-row gap-2">
        <FormFieldSelect
          name="inputModel::model"
          label="Input"
          className="grow"
          options={modelOptions}
          value={inputModel ? inputModel.id : '__GRUENT_IGNORE__'}
          onChange={(e) => setInputModel(models.find((m) => m.id === e.target.value))}
        />
        <Button text="Create input model" className="my-3" onClick={() => openModal('createModel')}>
          <RemixIcon icon={riAddLine} />
        </Button>
      </div>
      <div className="flex flex-row gap-2">
        <FormFieldSelect
          name="outputModel::model"
          label="Output"
          className="grow"
          options={modelOptions}
          value={outputModel ? outputModel.id : '__GRUENT_IGNORE__'}
          onChange={(e) => setOutputModel(models.find((m) => m.id === e.target.value))}
        />
        <Button
          text="Create output model"
          className="my-3"
          onClick={() => openModal('createModel')}
        >
          <RemixIcon icon={riAddLine} />
        </Button>
      </div>
      <Button className="mt-2" text="Transform your data" onClick={() => openModal('transform')} />
    </>
  )
}

type TransformModalProps = {
  transformation: Map<string, string>
  currentNode: ChannelNode
  inputModel?: Model
  outputModel?: Model
  orgId: string
  spaceId: string
  setTransformation: Dispatch<SetStateAction<Map<string, string>>>
}

const TransformNodeModal: FC<TransformModalProps> = ({
  transformation,
  inputModel,
  outputModel,
  orgId,
  spaceId,
  setTransformation,
}) => {
  const inputModelSchema: ModelAttributes = useMemo(() => {
    if (inputModel && inputModel.attributes)
      return JSON.parse(inputModel.attributes as string) as ModelAttributes

    return {}
  }, [inputModel])

  const outputModelSchema: ModelAttributes = useMemo(() => {
    if (outputModel && outputModel.attributes)
      return JSON.parse(outputModel.attributes as string) as ModelAttributes

    return {}
  }, [outputModel])

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

  const inputModelSchemaOptions = useMemo(() => {
    if (!inputModelSchema) return []

    const flatAttributes = flattenModelAttributes(inputModelSchema)
    return flatAttributes.map(({ key, value }) => ({
      label: key,
      value: value,
    }))
  }, [inputModelSchema, flattenModelAttributes])

  const outputAttributes = useMemo(() => {
    if (!outputModelSchema) return []

    return flattenModelAttributes(outputModelSchema)
  }, [outputModelSchema, flattenModelAttributes])

  const handleTransformAttributeChange = (e: ChangeEvent<HTMLSelectElement>, inputKey: string) => {
    const outputKey = e.target.value

    const updatedTransformation = new Map(transformation)
    updatedTransformation.set(inputKey, outputKey)
    if (outputKey === '__GRUENT_IGNORE__') updatedTransformation.delete(inputKey)

    setTransformation(updatedTransformation)
  }

  if (!inputModel || !outputModel) {
    return (
      <Modal id="transform" showAccept={false}>
        Missing input or output model. Make sure you set these before trying to apply a
        transformation.
      </Modal>
    )
  }

  return (
    <Modal id="transform" showAccept={false}>
      <div className="space-y-4">
        <h2 className="text-lg">
          {inputModel.name} to {outputModel.name}
        </h2>

        {outputAttributes.length <= 0 && (
          <p>
            No attributes were found in the '{outputModel.name}' model.{' '}
            <a
              className="text-blue-500 underline"
              href={`/organizations/${orgId}/spaces/${spaceId}/models/${outputModel.id}`}
            >
              Add at least one
            </a>{' '}
            and then try again.
          </p>
        )}

        <div className="space-y-2">
          {outputAttributes.map(({ key, value }) => {
            const defaultValue = transformation.get(value) ?? '__GRUENT_IGNORE__'

            return (
              <div key={value} className="grid grid-cols-[1fr_auto_1fr] items-center gap-4">
                <div className="text-right font-medium">{key}</div>
                <RemixIcon icon={riArrowLeftLine} className="text-gray-500" />
                <FormFieldSelect
                  name={key}
                  label={''}
                  defaultValue={defaultValue}
                  options={inputModelSchemaOptions}
                  onChange={(e) => handleTransformAttributeChange(e, value)}
                />
              </div>
            )
          })}
        </div>
      </div>
    </Modal>
  )
}

export default TransformProcessingOptions
