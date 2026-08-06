import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import { ChannelNode } from '@gruent/libs/engine/nodes/v1'
import { MinioActionSubnodeConfig } from '@gruent/libs/engine/nodes/v1/src/action/fileStorage/lib/minio'
import { Thing } from '@prisma/client'
import { FC, useEffect, useMemo, useState } from 'react'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import Button from '@gruent/libs/frontend/components/atoms/Button'

type Props = {
  things: Thing[]
  orgId: string
  spaceId: string
  selectedFileStorageThing: Thing
}

const MinioActionOptions: FC<Props> = ({ things, selectedFileStorageThing }) => {
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<MinioActionSubnodeConfig>,
  )

  const [minioKey, setMinioKey] = useState<string>()

  // load in default values
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!things) return

    setMinioKey(currentConfig.key)
  }, [currentNode, things])

  // update config and dependencies when options have been updated
  useEffect(() => {
    if (!minioKey) return
    if (!selectedFileStorageThing) return

    // update config
    const updatedNode = structuredClone(currentNode)

    updatedNode.metadata.config = {
      'minioBucket::thing': selectedFileStorageThing.id,
      key: minioKey,
    }

    // update dependencies array
    let minioBucketThingDependencyIndex = updatedNode.metadata.dependencies.things.findIndex(
      (d) => d.fieldName && d.fieldName === 'minioBucket::thing',
    )
    if (minioBucketThingDependencyIndex < 0)
      minioBucketThingDependencyIndex = updatedNode.metadata.dependencies.things.length

    updatedNode.metadata.dependencies.things[minioBucketThingDependencyIndex] = {
      fieldName: 'minioBucket::thing',
      internal: false,
      thingId: selectedFileStorageThing.id,
    }

    // apply updates to the nodes map
    upsertNode(currentNode.id, updatedNode)

    // eslint-disable-next-line -- currentNode in the dependencies array causes node configs to sometimes be copied
  }, [selectedFileStorageThing, minioKey, upsertNode])

  const minioBuckets = useMemo(() => {
    if (!things) return []

    return things.filter((thing) => {
      if (thing.thingCategory == 'MINIO_BUCKET') return true
      return false
    })
  }, [things])

  if (!minioBuckets || minioBuckets.length <= 0)
    return (
      <p>
        <Button variant="underline" onClick={() => openModal('createThing')}>
          Add a Minio Bucket
        </Button>{' '}
        to your space before using this node.
      </p>
    )

  return (
    <>
      <FormFieldText
        name="minioKey"
        label="Minio Key"
        value={minioKey || ''}
        onChange={(e) => setMinioKey(e.target.value)}
      />
    </>
  )
}

export default MinioActionOptions
