import useChannelEditorStore from '@iotea/hub/stores/channelEditorStore'
import { ChannelNode } from '@iotea/libs/engine/nodes/v1'
import { TimeSeriesDbNodeConfig } from '@iotea/libs/engine/nodes/v1/src/action/timeSeriesDb'
import { Thing } from '@prisma/client'
import { FC, useEffect, useMemo, useState } from 'react'
import FormFieldSelect from '@iotea/libs/frontend/components/atoms/FormFieldSelect'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import { openModal } from '@iotea/libs/frontend/hooks/useModal'

type Props = {
  things: Thing[]
  orgId: string
  spaceId: string
  timeSeriesDatabaseOptions: {
    label: string
    value: string
  }[]
  selectedTimeSeriesDatabase?: Thing
}

const InfluxDbActionOptions: FC<Props> = ({
  things,
  orgId,
  spaceId,
  timeSeriesDatabaseOptions: timeSeriesDatabaseOptions,
  selectedTimeSeriesDatabase: selectedTimeSeriesDatabase,
}) => {
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<TimeSeriesDbNodeConfig>,
  )

  const [bucket, setBucket] = useState<string>('')

  // load in default values
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!things) return

    setBucket(currentConfig?.bucket || '')
  }, [currentNode, things])

  // update config and dependencies when options have been updated
  useEffect(() => {
    if (!selectedTimeSeriesDatabase) return

    // update config
    const updatedNode = structuredClone(currentNode)

    updatedNode.metadata.config = {
      'influxdbDatabase::thing': selectedTimeSeriesDatabase.id,
      bucket: bucket,
    }

    // update things dependency array
    updatedNode.metadata.dependencies.things = [
      {
        fieldName: 'influxdbDatabase::thing',
        internal: false,
        thingId: selectedTimeSeriesDatabase.id,
      },
    ]

    // apply updates to the nodes map
    upsertNode(currentNode.id, updatedNode)

    // eslint-disable-next-line -- currentNode in the dependencies array causes node configs to sometimes be copied
  }, [bucket, selectedTimeSeriesDatabase, upsertNode])

  // Memoize database options
  const bucketOptions = useMemo((): { label: string; value: string }[] => {
    if (!selectedTimeSeriesDatabase || !selectedTimeSeriesDatabase.attributes) return []

    try {
      const attributes = JSON.parse(selectedTimeSeriesDatabase.attributes as string)

      if (attributes.buckets && attributes.buckets.length && attributes.buckets.length > 0) {
        return attributes.buckets.map((bucket) => ({
          label: bucket,
          value: bucket,
        }))
      }
    } catch (e) {
      console.error('Could not create options for the buckets select field', e)
    }

    return []
  }, [selectedTimeSeriesDatabase])

  if (!timeSeriesDatabaseOptions || timeSeriesDatabaseOptions.length <= 0)
    return (
      <p>
        <Button variant="underline" onClick={() => openModal('createThing')}>
          Add a time series database
        </Button>{' '}
        to your space before using this node.
      </p>
    )

  return (
    <>
      {bucketOptions.length <= 0 && (
        <p className="mb-2">
          <a
            href={`/organizations/${orgId}/spaces/${spaceId}/things/${selectedTimeSeriesDatabase!.id}`}
            className="text-blue-500 underline"
          >
            Add at least one bucket
          </a>{' '}
          to this InfluxDB server before using it.
        </p>
      )}
      {bucketOptions.length > 0 && (
        <FormFieldSelect
          options={bucketOptions}
          name="bucket"
          label="Bucket"
          value={bucket}
          onChange={(e) => setBucket(e.target.value)}
        />
      )}
    </>
  )
}

export default InfluxDbActionOptions
