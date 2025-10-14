import useChannelEditorStore from '@iotea/hub/stores/channelEditorStore'
import FormFieldSelect from '@iotea/libs/frontend/components/atoms/FormFieldSelect'
import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { ChannelNode } from '@iotea/libs/engine/nodes/v1'
import {
  MongoDbActionSubnodeConfig,
  MongoDbQueryMethod,
  MongoDbQueryMethodOptions,
} from '@iotea/libs/engine/nodes/v1/src/action/documentDb/lib/mongodb'
import { Thing, Model } from '@prisma/client'
import { FC, useEffect, useMemo, useState } from 'react'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import { openModal } from '@iotea/libs/frontend/hooks/useModal'

type Props = {
  things: Thing[]
  models: Model[]
  orgId: string
  spaceId: string
  selectedDocumentDbThing: Thing
}

const MongoDbActionOptions: FC<Props> = ({
  things,
  models,
  orgId,
  spaceId,
  selectedDocumentDbThing,
}) => {
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<MongoDbActionSubnodeConfig>,
  )

  // State variables to hold the form input values
  const [filter, setFilter] = useState<string>()
  const [database, setDatabase] = useState<string>()
  const [collection, setCollection] = useState<string>()
  const [queryMethod, setQueryMethod] = useState<MongoDbQueryMethod>()
  const [document, setDocument] = useState<string>()
  const [modelAttributes, setModelAttributes] = useState<{ [key: string]: string }>({})
  const [selectedTemplateModel, setSelectedTemplateModel] = useState<Model>()
  const [defaultValues, setDefaultValues] = useState<{ [key: string]: any }>({})

  // Load default values from the current node config
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!things) return

    setFilter(currentConfig.filter)
    setDocument(currentConfig.document)
    setQueryMethod(currentConfig.queryMethod)
    setDatabase(currentConfig.database)
    setCollection(currentConfig.collection)
    setDefaultValues(currentConfig.defaultValues || {})

    // Load initial model if available
    if (models && currentConfig['input::model']) {
      const model = models.find((m) => m.id === currentConfig['input::model'])
      if (model) setSelectedTemplateModel(model)
    }
  }, [currentNode, things, models])

  // Update selected attributes whenever the message template changes
  const templateVariables = useMemo(() => {
    if (!selectedTemplateModel || !selectedTemplateModel.attributes || !filter) return

    // Match the message to any attempts at using a variable
    const templateVariables = (() => {
      const matches = filter.match(/\{(\w+)\}/g) || []
      return matches.map((match) => match.slice(1, -1))
    })()

    try {
      const selectedModelAttributes: { [key: string]: string } = JSON.parse(
        selectedTemplateModel.attributes.toString(),
      )

      const filteredTemplateVariables = templateVariables.filter((field) => {
        if (selectedModelAttributes) {
          const modelAttributesKeys = Object.keys(selectedModelAttributes)
          if (modelAttributesKeys && modelAttributesKeys.includes(field)) return true
        }

        return false
      })

      return filteredTemplateVariables
    } catch (e) {
      console.error(
        `could not parse the current input model attributes: ${selectedTemplateModel.attributes}`,
        e,
      )
      return []
    }
  }, [selectedTemplateModel, filter])

  // Parse the attributes and update attributes whenever the selected model changes
  useEffect(() => {
    if (!selectedTemplateModel || !selectedTemplateModel.attributes) return

    setModelAttributes({})

    try {
      const parsedAttributes: { [key: string]: string } = JSON.parse(
        selectedTemplateModel.attributes.toString(),
      )
      setModelAttributes(parsedAttributes)
    } catch (err) {
      addToast({
        title: 'Could not update the message template model',
        body: `Received err: ${err}`,
      })
    }
  }, [selectedTemplateModel])

  // update config and dependencies when options have been updated
  useEffect(() => {
    if (!selectedDocumentDbThing) return

    // update config
    const updatedNode = structuredClone(currentNode)

    // Update the node's config with the current form and selected thing/model values
    updatedNode.metadata.config = {
      'mongoDbServer::thing': selectedDocumentDbThing.id,
      'input::model': selectedTemplateModel?.id,
      queryMethod,
      filter,
      document,
      database,
      collection,
      defaultValues,
    }

    // Update the dependencies for the MongoDB thing in the node
    let mongoDbServerThingDependencyIndex = updatedNode.metadata.dependencies.things.findIndex(
      (d) => d.fieldName && d.fieldName === 'mongoDbServer::thing',
    )
    if (mongoDbServerThingDependencyIndex < 0)
      mongoDbServerThingDependencyIndex = updatedNode.metadata.dependencies.things.length

    updatedNode.metadata.dependencies.things[mongoDbServerThingDependencyIndex] = {
      fieldName: 'mongoDbServer::thing',
      internal: false,
      thingId: selectedDocumentDbThing.id,
    }

    // Update the dependencies for the model in the node
    let modelDependencyIndex = updatedNode.metadata.dependencies.models.findIndex(
      (d) => d.fieldName && d.fieldName === 'input::model',
    )
    if (modelDependencyIndex < 0)
      modelDependencyIndex = updatedNode.metadata.dependencies.models.length

    if (selectedTemplateModel) {
      updatedNode.metadata.dependencies.models[modelDependencyIndex] = {
        fieldName: 'input::model',
        internal: false,
        modelId: selectedTemplateModel?.id,
      }
    }

    // Apply the updates to the node using the upsertNode function
    upsertNode(currentNode.id, updatedNode)

    // eslint-disable-next-line -- currentNode in the dependencies array causes node configs to sometimes be copied
  }, [
    selectedDocumentDbThing,
    filter,
    document,
    queryMethod,
    database,
    collection,
    selectedTemplateModel,
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
        // Check if the value is an incomplete float (ending with a period)
        if (value.match(/^\d+\.$/)) {
          // If value ends with a period (e.g., "55."), don't parse it yet
          typedValue = value
        } else {
          // Otherwise, parse it as a float
          typedValue = parseFloat(value)
          if (isNaN(typedValue)) {
            typedValue = value // Fallback to string if not a valid float
          }
        }
        break

      case 'int':
      case 'integer':
        typedValue = parseInt(value, 10)
        if (isNaN(typedValue)) {
          typedValue = value // Fallback to string if not a valid integer
        }
        break

      case 'bool':
      case 'boolean':
        // Handle boolean conversion by checking common string representations
        typedValue =
          value.toLowerCase() === 'true' ? true : value.toLowerCase() === 'false' ? false : value
        break

      case 'string':
      default:
        typedValue = value // Keep the value as a string
        break
    }

    // Update the defaultValues object with the correctly typed value
    setDefaultValues((prev) => ({
      ...prev,
      [attribute]: typedValue,
    }))
  }

  // Ensure that a MongoDB server is selected
  const mongoDbServers = useMemo(() => {
    if (!things) return []

    return things.filter((thing) => thing.thingCategory == 'MONGODB_SERVER')
  }, [things])

  // Memoize database options
  const databaseOptions = useMemo((): { label: string; value: string }[] => {
    if (!selectedDocumentDbThing || !selectedDocumentDbThing.attributes) return []

    try {
      const attributes = JSON.parse(selectedDocumentDbThing.attributes as string)

      if (attributes.databases && attributes.databases.length && attributes.databases.length > 0) {
        return attributes.databases.map((database) => ({
          label: database,
          value: database,
        }))
      }
    } catch (e) {
      console.error('Could not create options for the database select field', e)
    }

    return []
  }, [selectedDocumentDbThing])

  // Memoize database collection options
  const collectionOptions = useMemo((): { label: string; value: string }[] => {
    if (!selectedDocumentDbThing || !selectedDocumentDbThing.attributes) return []
    if (!database) return []

    try {
      const attributes = JSON.parse(selectedDocumentDbThing.attributes as string)
      if (
        !attributes.collections ||
        typeof attributes.collections !== 'object' ||
        !attributes.collections[database]
      )
        return []
      if (!attributes.collections[database].length) return []

      return attributes.collections[database].map((collection) => ({
        label: collection,
        value: collection,
      }))
    } catch (e) {
      console.error('Could not create options for the database select field', e)
    }

    return []
  }, [selectedDocumentDbThing, database])

  if (!mongoDbServers || mongoDbServers.length <= 0)
    return (
      <p>
        <Button variant="underline" onClick={() => openModal('createThing')}>
          Add a MongoDB server
        </Button>{' '}
        to your space before using this node.
      </p>
    )

  // Render form fields for the MongoDB node configuration
  return (
    <>
      {databaseOptions.length <= 0 && (
        <p>
          <a
            href={`/organizations/${orgId}/spaces/${spaceId}/things/${selectedDocumentDbThing.id}`}
            className="text-blue-500 underline"
          >
            Add at least one database
          </a>{' '}
          to this MongoDB server before using it.
        </p>
      )}
      {databaseOptions.length > 0 && (
        <FormFieldSelect
          options={databaseOptions}
          name="database"
          label="Database"
          value={database || ''}
          onChange={(e) => setDatabase(e.target.value)}
        />
      )}
      {database && collectionOptions.length <= 0 && (
        <p>
          <a
            href={`/organizations/${orgId}/spaces/${spaceId}/things/${selectedDocumentDbThing.id}`}
            className="text-blue-500 underline"
          >
            Add at least one collection
          </a>{' '}
          to this database before using it.
        </p>
      )}
      {database && collectionOptions.length > 0 && (
        <FormFieldSelect
          options={collectionOptions}
          name="collection"
          label="Collection"
          value={collection || ''}
          onChange={(e) => setCollection(e.target.value)}
        />
      )}
      <FormFieldSelect
        name="method"
        label="Query Method"
        options={MongoDbQueryMethodOptions}
        value={queryMethod || ''}
        onChange={(e) => setQueryMethod(e.target.value as MongoDbQueryMethod | undefined)}
      />
      {(queryMethod === MongoDbQueryMethod.Find ||
        queryMethod === MongoDbQueryMethod.FindOne ||
        queryMethod === MongoDbQueryMethod.FindOneAndUpdate ||
        queryMethod === MongoDbQueryMethod.FindOneAndReplace ||
        queryMethod === MongoDbQueryMethod.FindOneAndDelete ||
        queryMethod === MongoDbQueryMethod.DeleteOne ||
        queryMethod === MongoDbQueryMethod.DeleteMany ||
        queryMethod === MongoDbQueryMethod.ReplaceOne ||
        queryMethod === MongoDbQueryMethod.UpdateOne) && (
        // TODO:
        // queryMethod === MongoDbQueryMethod.UpdateMany
        <>
          <label>Query Filter</label>
          <textarea
            className="w-full min-h-36 p-2"
            name="filter"
            value={filter || ''}
            onChange={(e) => setFilter(e.target.value)}
          />
        </>
      )}
      {(queryMethod === MongoDbQueryMethod.FindOneAndUpdate ||
        queryMethod === MongoDbQueryMethod.FindOneAndReplace ||
        queryMethod === MongoDbQueryMethod.InsertOne ||
        queryMethod === MongoDbQueryMethod.ReplaceOne ||
        queryMethod === MongoDbQueryMethod.UpdateOne) && (
        // TODO:
        // queryMethod === MongoDbQueryMethod.InsertMany
        // queryMethod === MongoDbQueryMethod.UpdateMany
        <>
          <label>Document</label>
          <textarea
            className="w-full min-h-36 p-2"
            name="document"
            value={document || ''}
            onChange={(e) => setDocument(e.target.value)}
          />
        </>
      )}
      {selectedDocumentDbThing && models && (
        <>
          <FormFieldSelect
            name="templateModel"
            label="Template model (optional)"
            options={modelOptions} // Populate the model dropdown
            value={selectedTemplateModel?.id || '__IOTEA_IGNORE__'}
            onChange={(e) => setSelectedTemplateModel(models.find((m) => m.id === e.target.value))}
          />
        </>
      )}
      {selectedTemplateModel && templateVariables && templateVariables.length > 0 && (
        <>
          <h3>Message Template Variables</h3>
          {templateVariables.map((field) => (
            <FormFieldText
              name={field}
              key={field}
              label={field}
              placeholder={`Default value for ${field} (optional)`}
              value={defaultValues[field] || ''}
              onChange={(e) => handleDefaultValueChange(field, e.target.value)}
            />
          ))}
        </>
      )}
    </>
  )
}

export default MongoDbActionOptions
