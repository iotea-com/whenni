import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import { ChannelNode } from '@gruent/libs/engine/nodes/v1'
import { DocumentDbActionNodeConfig } from '@gruent/libs/engine/nodes/v1/src/action/documentDb'
import { Thing, Model } from '@prisma/client'
import { FC, useEffect, useMemo, useState } from 'react'
import MongoDbActionOptions from './MongoDbActionOptions'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import { RemixIcon, riAddLine } from '@mwarnerdotme/react-remixicon'

type Props = {
  things: Thing[]
  orgId: string
  spaceId: string
  models: Model[]
}

export enum DocumentDB {
  MONGODB_SERVER = 'MONGODB_SERVER',
}

const DocumentDbActionOptions: FC<Props> = ({ things, orgId, spaceId, models }) => {
  // Accessing the current node in the editor store
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<DocumentDbActionNodeConfig>,
  )
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const setCurrentNode = useChannelEditorStore((state) => state.setCurrentNode)

  // State to hold the selected document DB thing
  const [selectedDocumentDatabaseThing, setSelectedDocumentDbThing] = useState<Thing>()

  // Load initial config
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!things || !models) return

    // Check if a MongoDB thing is already selected in the config
    if (Object.keys(currentConfig).includes('mongoDbServer::thing')) {
      const documentDb = things.find((s) => {
        if (s.id === currentConfig['mongoDbServer::thing']) return true
        return false
      })
      if (documentDb) setSelectedDocumentDbThing(documentDb)
      return
    }
  }, [currentNode, things, models])

  // Reset database and collection when selected Document DB thing changes
  useEffect(() => {
    if (!selectedDocumentDatabaseThing) return
    if (selectedDocumentDatabaseThing.id === currentNode.metadata.config['mongoDbServer::thing'])
      return

    const updatedNode = structuredClone(currentNode)
    updatedNode.metadata.config = {
      ...currentNode.metadata.config,
      'mongoDbServer::thing': selectedDocumentDatabaseThing.id,
      database: undefined,
      collection: undefined,
    }
    updatedNode.metadata.dependencies.things = []
    updatedNode.metadata.dependencies.models = currentNode?.metadata?.dependencies?.models

    // Apply the updates to the node using the upsertNode function
    upsertNode(currentNode.id, updatedNode)
    setCurrentNode(updatedNode)
  }, [selectedDocumentDatabaseThing, currentNode, upsertNode, setCurrentNode])

  // Determine the database type based on the selected thing
  const selectedDocumentDbType = useMemo(() => {
    if (!selectedDocumentDatabaseThing) return

    if (selectedDocumentDatabaseThing.thingCategory === 'MONGODB_SERVER')
      return DocumentDB.MONGODB_SERVER
  }, [selectedDocumentDatabaseThing])

  // Generate a filtered list of available database types
  const documentDbThings = useMemo(() => {
    if (!things) return []

    const mongodbServers = things.filter((thing) => {
      if (thing.thingCategory == 'MONGODB_SERVER') return true
      return false
    })

    return [...mongodbServers]
  }, [things])

  // Generate options for the document database dropdown
  const documentDbOptions = useMemo(() => {
    return documentDbThings.map((thing) => ({
      label: thing.name,
      value: thing.id,
    }))
  }, [documentDbThings])

  // Conditionally rendering subnode options based on the selected type (SES/SNS)
  const selectedSubnodeOptions = useMemo(() => {
    if (!selectedDocumentDatabaseThing) return

    switch (selectedDocumentDbType) {
      case DocumentDB.MONGODB_SERVER:
        return (
          <MongoDbActionOptions
            things={things}
            orgId={orgId}
            spaceId={spaceId}
            selectedDocumentDbThing={selectedDocumentDatabaseThing}
            models={models}
          />
        )
      default:
        return
    }
  }, [models, selectedDocumentDbType, selectedDocumentDatabaseThing, things, orgId, spaceId])

  if (!documentDbThings || documentDbThings.length <= 0)
    return (
      <p className="mb-2">
        <Button variant="underline" onClick={() => openModal('createThing')}>
          Add a document database
        </Button>{' '}
        to your space before using this node.
      </p>
    )

  return (
    <>
      <div className="flex flex-row gap-2">
        <FormFieldSelect
          name="DocumentDbThing"
          label="Document Database Thing"
          className="grow"
          options={documentDbOptions}
          value={selectedDocumentDatabaseThing?.id ?? '__GRUENT_IGNORE__'}
          onChange={(e) =>
            setSelectedDocumentDbThing(
              documentDbThings.find((s) => {
                if (s.id === e.target.value) return true
                return false
              }),
            )
          }
        />
        <Button className="my-3" onClick={() => openModal('createThing')}>
          <RemixIcon icon={riAddLine} />
        </Button>
      </div>
      {selectedSubnodeOptions}
    </>
  )
}

export default DocumentDbActionOptions
