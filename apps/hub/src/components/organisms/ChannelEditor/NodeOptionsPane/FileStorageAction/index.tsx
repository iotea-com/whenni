import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import { ChannelNode } from '@gruent/libs/engine/nodes/v1'
import { FileStorageActionNodeConfig } from '@gruent/libs/engine/nodes/v1/src/action/fileStorage'
import { Thing } from '@prisma/client'
import { FC, useEffect, useMemo, useState } from 'react'
import S3ActionOptions from './S3ActionOptions'
import MinioActionOptions from './MinioActionOptions'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import { RemixIcon, riAddLine } from '@mwarnerdotme/react-remixicon'

type Props = {
  things: Thing[]
  orgId: string
  spaceId: string
}

export enum FileStorage {
  S3 = 'S3',
  Minio = 'Minio',
}

const FileStorageActionOptions: FC<Props> = ({ things, orgId, spaceId }) => {
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<FileStorageActionNodeConfig>,
  )

  const [selectedFileStorageThing, setSelectedFileStorageThing] = useState<Thing>()

  // Load the initial storage type from the currentNode config if available
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!things) return

    if (Object.keys(currentConfig).includes('s3bucket::thing')) {
      const s3Bucket = things.find((s) => {
        if (s.id === currentConfig['s3bucket::thing']) return true
        return false
      })
      if (s3Bucket) setSelectedFileStorageThing(s3Bucket)
      return
    }

    if (Object.keys(currentConfig).includes('minioBucket::thing')) {
      const minioBucket = things.find((s) => {
        if (s.id === currentConfig['minioBucket::thing']) return true
        return false
      })
      if (minioBucket) setSelectedFileStorageThing(minioBucket)
      return
    }
  }, [currentNode, things])

  const selectedFileStorageType = useMemo(() => {
    if (!selectedFileStorageThing) return

    if (selectedFileStorageThing.thingCategory === 'S3_BUCKET') return FileStorage.S3

    if (selectedFileStorageThing.thingCategory === 'MINIO_BUCKET') return FileStorage.Minio
  }, [selectedFileStorageThing])

  const fileStorageThings = useMemo(() => {
    if (!things) return []

    const s3Buckets = things.filter((thing) => {
      if (thing.thingCategory == 'S3_BUCKET') return true
      return false
    })

    const minioBuckets = things.filter((thing) => {
      if (thing.thingCategory == 'MINIO_BUCKET') return true
      return false
    })

    return [...s3Buckets, ...minioBuckets]
  }, [things])

  const fileStorageThingOptions = useMemo(() => {
    return fileStorageThings.map((thing) => ({
      label: thing.name,
      value: thing.id,
    }))
  }, [fileStorageThings])

  // Conditionally rendering subnode options based on the selected type (S3/Minio)
  const selectedSubnodeOptions = useMemo(() => {
    if (!selectedFileStorageThing) return

    switch (selectedFileStorageType) {
      case FileStorage.S3:
        return (
          <S3ActionOptions
            things={things}
            orgId={orgId}
            spaceId={spaceId}
            selectedFileStorageThing={selectedFileStorageThing}
          />
        )
      case FileStorage.Minio:
        return (
          <MinioActionOptions
            things={things}
            orgId={orgId}
            spaceId={spaceId}
            selectedFileStorageThing={selectedFileStorageThing}
          />
        )
      default:
        return
    }
  }, [selectedFileStorageType, selectedFileStorageThing, things, orgId, spaceId])

  if (!fileStorageThings || fileStorageThings.length <= 0)
    return (
      <p className="mb-2">
        <Button variant="underline" onClick={() => openModal('createThing')}>
          Add a file storage server
        </Button>{' '}
        to your space before using this node.
      </p>
    )

  return (
    <>
      <div className="flex flex-row gap-2">
        <FormFieldSelect
          name="fileStorageThing"
          label="File Storage Thing"
          className="grow"
          options={fileStorageThingOptions}
          value={selectedFileStorageThing?.id ?? '__GRUENT_IGNORE__'}
          onChange={(e) =>
            setSelectedFileStorageThing(
              fileStorageThings.find((s) => {
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

export default FileStorageActionOptions
