import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import { ChannelNode } from '@gruent/libs/engine/nodes/v1'
import { TimeSeriesDbNodeConfig } from '@gruent/libs/engine/nodes/v1/src/action/timeSeriesDb'
import { Thing } from '@prisma/client'
import { FC, useEffect, useMemo, useState } from 'react'
import InfluxDbActionOptions from './InfluxDbActionOptions'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import { RemixIcon, riAddLine } from '@mwarnerdotme/react-remixicon'

type Props = {
  things: Thing[]
  orgId: string
  spaceId: string
}

export enum TimeSeriesDb {
  InfluxDB = 'InfluxDB',
}

const TimeSeriesDbActionOptions: FC<Props> = ({ things, orgId, spaceId }) => {
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<TimeSeriesDbNodeConfig>,
  )

  const [selectedTimeSeriesDatabase, setSelectedTimeSeriesDatabase] = useState<Thing>()

  // load in the default message queue client from the current node config
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!things) return

    if (Object.keys(currentConfig).includes('influxdbDatabase::thing')) {
      const influxdbDatabase = things.find((s) => {
        if (s.id === currentConfig['influxdbDatabase::thing']) return true
        return false
      })
      if (influxdbDatabase) setSelectedTimeSeriesDatabase(influxdbDatabase)
      return
    }
  }, [currentNode, things])

  const selectedTimeSeriesDatabaseType = useMemo(() => {
    if (!selectedTimeSeriesDatabase) return

    if (selectedTimeSeriesDatabase.thingCategory === 'INFLUXDB_DATABASE')
      return TimeSeriesDb.InfluxDB

    return
  }, [selectedTimeSeriesDatabase])

  const timeSeriesDatabases = useMemo(() => {
    if (!things) return []

    const influxDb = things.filter((thing) => {
      if (thing.thingCategory == 'INFLUXDB_DATABASE') return true
      return false
    })

    return [...influxDb]
  }, [things])

  const timeSeriesDatabaseOptions = useMemo(() => {
    return timeSeriesDatabases.map((thing) => ({
      label: thing.name,
      value: thing.id,
    }))
  }, [timeSeriesDatabases])

  const selectedSubnodeOptions = useMemo(() => {
    if (!selectedTimeSeriesDatabaseType) return

    switch (selectedTimeSeriesDatabaseType) {
      case TimeSeriesDb.InfluxDB:
        return (
          <InfluxDbActionOptions
            things={things}
            orgId={orgId}
            spaceId={spaceId}
            timeSeriesDatabaseOptions={timeSeriesDatabaseOptions}
            selectedTimeSeriesDatabase={selectedTimeSeriesDatabase}
          />
        )
    }
  }, [
    selectedTimeSeriesDatabaseType,
    timeSeriesDatabaseOptions,
    selectedTimeSeriesDatabase,
    orgId,
    spaceId,
    things,
  ])

  if (!timeSeriesDatabases || timeSeriesDatabases.length <= 0)
    return (
      <p className="mb-2">
        <Button variant="underline" onClick={() => openModal('createThing')}>
          Add a time series database
        </Button>{' '}
        to your space before using this node.
      </p>
    )

  return (
    <>
      <div className="flex flex-row gap-2">
        <FormFieldSelect
          name="timeSeriesDatabase"
          label="Time Series Database"
          className="grow"
          options={timeSeriesDatabaseOptions}
          value={selectedTimeSeriesDatabase?.id ?? '__GRUENT_IGNORE__'}
          onChange={(e) =>
            setSelectedTimeSeriesDatabase(
              timeSeriesDatabases.find((s) => {
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

export default TimeSeriesDbActionOptions
