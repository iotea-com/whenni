import useChannelEditorStore from '@iotea/hub/stores/channelEditorStore'
import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import { ChannelNode } from '@iotea/libs/engine/nodes/v1'
import { S3ActionSubnodeConfig } from '@iotea/libs/engine/nodes/v1/src/action/fileStorage/lib/awsS3'
import { Thing } from '@prisma/client'
import { FC, useEffect, useMemo, useState } from 'react'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import { openModal } from '@iotea/libs/frontend/hooks/useModal'

type Props = {
  things: Thing[]
  orgId: string
  spaceId: string
  selectedFileStorageThing: Thing
}

const S3ActionOptions: FC<Props> = ({ things, selectedFileStorageThing }) => {
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<S3ActionSubnodeConfig>,
  )

  const [s3Key, setS3Key] = useState<string>()

  // load in default values
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!things) return

    setS3Key(currentConfig.key)
  }, [currentNode, things])

  // update config and dependencies when options have been updated
  useEffect(() => {
    if (!s3Key) return
    if (!selectedFileStorageThing) return

    // update config
    const updatedNode = structuredClone(currentNode)

    updatedNode.metadata.config = {
      's3bucket::thing': selectedFileStorageThing.id,
      key: s3Key,
    }

    // update dependencies array
    let s3BucketThingDependencyIndex = updatedNode.metadata.dependencies.things.findIndex(
      (d) => d.fieldName && d.fieldName === 's3bucket::thing',
    )
    if (s3BucketThingDependencyIndex < 0)
      s3BucketThingDependencyIndex = updatedNode.metadata.dependencies.things.length

    updatedNode.metadata.dependencies.things[s3BucketThingDependencyIndex] = {
      fieldName: 's3bucket::thing',
      internal: false,
      thingId: selectedFileStorageThing.id,
    }

    // apply updates to the nodes map
    upsertNode(currentNode.id, updatedNode)

    // eslint-disable-next-line -- currentNode in the dependencies array causes node configs to sometimes be copied
  }, [selectedFileStorageThing, s3Key, upsertNode])

  const s3Buckets = useMemo(() => {
    if (!things) return []

    return things.filter((thing) => {
      if (thing.thingCategory == 'S3_BUCKET') return true
      return false
    })
  }, [things])

  if (!s3Buckets || s3Buckets.length <= 0)
    return (
      <p>
        <Button variant="underline" onClick={() => openModal('createThing')}>
          Add an S3 Bucket
        </Button>{' '}
        to your space before using this node.
      </p>
    )

  return (
    <>
      <FormFieldText
        name="s3Key"
        label="S3 Key"
        value={s3Key || ''}
        onChange={(e) => setS3Key(e.target.value)}
      />
    </>
  )
}

export default S3ActionOptions
