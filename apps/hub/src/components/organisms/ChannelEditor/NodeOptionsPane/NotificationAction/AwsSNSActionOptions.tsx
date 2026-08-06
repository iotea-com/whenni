import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import { addToast } from '@gruent/libs/frontend/hooks/useToast'
import { ChannelNode } from '@gruent/libs/engine/nodes/v1'
import { AwsSNSActionSubnodeConfig } from '@gruent/libs/engine/nodes/v1/src/action/notification/lib/awsSNS'
import { Thing, Model } from '@prisma/client'
import { FC, useEffect, useMemo, useState } from 'react'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import Button from '@gruent/libs/frontend/components/atoms/Button'

type Props = {
  things: Thing[]
  models: Model[]
  orgId: string
  spaceId: string
  selectedNotificationThing: Thing
}

const AwsSNSOptions: FC<Props> = ({ things, models, selectedNotificationThing }) => {
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<AwsSNSActionSubnodeConfig>,
  )

  // State variables to hold the form input values
  const [topicARN, setTopicARN] = useState<string>()
  const [templateMessage, setTemplateMessage] = useState<string>()
  const [selectedTemplateModel, setSelectedTemplateModel] = useState<Model>()
  const [modelAttributes, setModelAttributes] = useState<{ [key: string]: string }>({})
  const [defaultValues, setDefaultValues] = useState<{ [key: string]: any }>({})

  // Load default values from the current node config
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!things) return

    setTopicARN(currentConfig.topicARN)
    setTemplateMessage(currentConfig.templateMessage) // Load the template
    setDefaultValues(currentConfig.defaultValues || {}) // Load default values

    // Load initial model if available
    if (models && currentConfig['input::model']) {
      const model = models.find((m) => m.id === currentConfig['input::model'])
      if (model) setSelectedTemplateModel(model)
    }
  }, [currentNode, things, models])

  // Update selected attributes whenever the message template changes
  const templateVariables = useMemo(() => {
    if (!selectedTemplateModel || !selectedTemplateModel.attributes || !templateMessage) return

    // Match the message to any attempts at using a variable
    const templateVariables = (() => {
      const matches = templateMessage.match(/\{([^}]+)\}/g) || []
      return matches.map((match) => match.slice(1, -1))
    })()

    try {
      const selectedModelAttributes: { [key: string]: string } = JSON.parse(
        selectedTemplateModel.attributes.toString(),
      )

      const filteredTemplateVariables = templateVariables.filter((attribute) => {
        if (selectedModelAttributes) {
          const schemaKeys = Object.keys(selectedModelAttributes)
          if (schemaKeys && schemaKeys.includes(attribute)) return true
        }

        return false
      })

      return filteredTemplateVariables
    } catch (e) {
      console.error(
        `could not parse the current input model schema: ${selectedTemplateModel.attributes}`,
        e,
      )
      return []
    }
  }, [selectedTemplateModel, templateMessage])

  // Parse the schema and update attributes whenever the selected model changes
  useEffect(() => {
    if (!selectedTemplateModel || !selectedTemplateModel.attributes) return

    setModelAttributes({})

    try {
      const parsedSchema: { [key: string]: string } = JSON.parse(
        selectedTemplateModel.attributes.toString(),
      )
      setModelAttributes(parsedSchema)
    } catch (err) {
      addToast({
        title: 'Could not update the message template model',
        body: `Received err: ${err}`,
      })
    }
  }, [selectedTemplateModel])

  // update config and dependencies when options have been updated
  useEffect(() => {
    if (!topicARN || !templateMessage || !selectedNotificationThing) return

    // update config
    const updatedNode = structuredClone(currentNode)

    // Update the node's config with the current form and selected thing/model values
    updatedNode.metadata.config = {
      'awsSNS::thing': selectedNotificationThing.id,
      'input::model': selectedTemplateModel?.id,
      topicARN: topicARN,
      templateMessage: templateMessage,
      defaultValues: defaultValues,
    }

    // Update the dependencies for the SNS thing in the node
    let awsSNSThingDependencyIndex = updatedNode.metadata.dependencies.things.findIndex(
      (d) => d.fieldName && d.fieldName === 'awsSNS::thing',
    )
    if (awsSNSThingDependencyIndex < 0)
      awsSNSThingDependencyIndex = updatedNode.metadata.dependencies.things.length

    updatedNode.metadata.dependencies.things[awsSNSThingDependencyIndex] = {
      fieldName: 'awsSNS::thing',
      internal: false,
      thingId: selectedNotificationThing.id,
    }

    // Update the dependencies for the model in the node
    let awsSNSModelDependencyIndex = updatedNode.metadata.dependencies.models.findIndex(
      (d) => d.fieldName && d.fieldName === 'input::model',
    )
    if (awsSNSModelDependencyIndex < 0)
      awsSNSModelDependencyIndex = updatedNode.metadata.dependencies.models.length

    if (selectedTemplateModel) {
      updatedNode.metadata.dependencies.models[awsSNSModelDependencyIndex] = {
        fieldName: 'input::model',
        internal: false,
        modelId: selectedTemplateModel?.id,
      }
    } else {
      updatedNode.metadata.dependencies.models.splice(awsSNSModelDependencyIndex, 1)
    }

    // Apply the updates to the node using the upsertNode function
    upsertNode(currentNode.id, updatedNode)

    // eslint-disable-next-line -- currentNode in the dependencies array causes node configs to sometimes be copied
  }, [
    selectedNotificationThing,
    selectedTemplateModel,
    topicARN,
    templateMessage,
    templateVariables,
    defaultValues,
    upsertNode,
  ])

  // Generate options for the Model dropdown
  const modelOptions = useMemo(() => {
    if (!models || models.length === 0) return []

    return models.map((model) => ({
      label: model.name,
      value: model.id,
    }))
  }, [models])

  const handleDefaultValueChange = (attribute: string, value: string) => {
    const expectedType = modelAttributes[attribute]
    let typedValue: any = value

    // Switch statement for type conversion based on the expected type
    switch (expectedType) {
      case 'float':
        // Parse as float unless it's incomplete (e.g., "55.")
        typedValue = value.match(/^\d+\.$/) ? value : parseFloat(value)
        if (isNaN(typedValue)) typedValue = value
        break

      case 'int':
      case 'integer':
        typedValue = parseInt(value, 10)
        if (isNaN(typedValue)) typedValue = value
        break

      case 'bool':
      case 'boolean':
        typedValue =
          value.toLowerCase() === 'true' ? true : value.toLowerCase() === 'false' ? false : value
        break

      case 'string':
      default:
        typedValue = value // Keep as string
        break
    }

    // Update the defaultValues object with the correctly typed value
    setDefaultValues((prev) => ({
      ...prev,
      [attribute]: typedValue,
    }))
  }

  // Ensure that AWS SNS is selected
  const awsSNS = useMemo(() => {
    if (!things) return []

    return things.filter((thing) => thing.thingCategory == 'AWS_SNS_ENDPOINT')
  }, [things])

  if (!awsSNS || awsSNS.length <= 0)
    return (
      <p>
        <Button variant="underline" onClick={() => openModal('createThing')}>
          Add an AWS SNS Endpoint
        </Button>{' '}
        to your space before using this node.
      </p>
    )

  // Render form fields for the SNS node configuration
  return (
    <>
      <FormFieldText
        name="topicARN"
        label="Topic ARN"
        value={topicARN || ''}
        onChange={(e) => setTopicARN(e.target.value)}
      />
      <label>Template Message</label>
      <textarea
        className="w-full min-h-36 p-2"
        name="templateMessage"
        value={templateMessage || ''}
        onChange={(e) => setTemplateMessage(e.target.value)}
      />
      {selectedNotificationThing && models && (
        <>
          {/* Add model dropdown after selecting SES or SNS */}
          <FormFieldSelect
            name="templateModel"
            label="Template model (optional)"
            options={modelOptions} // Populate the model dropdown
            value={selectedTemplateModel?.id || '__GRUENT_IGNORE__'}
            onChange={(e) => setSelectedTemplateModel(models.find((s) => s.id === e.target.value))}
          />
        </>
      )}
      {selectedTemplateModel && templateVariables && templateVariables.length > 0 && (
        <>
          <h3>Message Template Variables</h3>
          {templateVariables.map((attribute) => (
            <FormFieldText
              name={attribute}
              key={attribute}
              label={attribute}
              placeholder={`Default value for ${attribute} (optional)`}
              value={defaultValues[attribute] || ''}
              onChange={(e) => handleDefaultValueChange(attribute, e.target.value)}
            />
          ))}
        </>
      )}
    </>
  )
}

export default AwsSNSOptions
