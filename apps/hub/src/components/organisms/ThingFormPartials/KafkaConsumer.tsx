'use client'

import useAuth from '@gruent/hub/hooks/useAuth'
import gruentClient from '@gruent/hub/lib/gruent'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import { Thing } from '@prisma/client'
import { useQuery } from '@tanstack/react-query'
import { Dispatch, FC, useEffect, useMemo, useState } from 'react'

type Props = {
  setFormData: Dispatch<Map<string, string | number | string[] | number[]>>
  spaceId: string
  initial?: Partial<Thing>
}

const KafkaConsumer: FC<Props> = ({ setFormData, spaceId, initial }) => {
  const attributes = useMemo(() => {
    if (initial && initial.attributes && typeof initial.attributes === 'string')
      return JSON.parse(initial.attributes)

    if (initial && initial.attributes && typeof initial.attributes === 'object')
      return initial.attributes

    return undefined
  }, [initial])

  const [clusterId, setClusterId] = useState<string>(attributes?.cluster)
  const [groupId, setGroupId] = useState<string>(attributes?.groupId)
  const [autoOffsetReset, setAutoOffsetReset] = useState<'latest' | 'earliest' | 'none'>(
    attributes?.autoOffsetReset ?? 'earliest',
  )

  const { accessToken } = useAuth()

  const {
    status: getKafkaClustersStatus,
    data: kafkaClusters,
    error: getKafkaClustersError,
  } = useQuery({
    queryKey: ['things', 'kafkaClusters', spaceId],
    queryFn: async () => {
      if (!accessToken) throw new Error('Invalid auth session.')

      const { data: things, errors } = await gruentClient(accessToken).things.list(spaceId, {
        category: 'KAFKA_CLUSTER',
      })

      if (errors && errors.length > 0) throw new Error(errors[0])

      return things
    },
  })

  const kafkaClusterOptions = useMemo(() => {
    if (getKafkaClustersError) return []
    if (getKafkaClustersStatus !== 'success') return []
    if (!kafkaClusters) return []

    return kafkaClusters.map((kafkaCluster) => {
      return {
        label: kafkaCluster.name,
        value: kafkaCluster.id,
      }
    })
  }, [kafkaClusters, getKafkaClustersError, getKafkaClustersStatus])

  useEffect(() => {
    const updatedFormData = new Map()
    updatedFormData.set('cluster', clusterId)
    updatedFormData.set('groupId', groupId)
    updatedFormData.set('autoOffsetReset', autoOffsetReset)
    setFormData(updatedFormData)
  }, [clusterId, groupId, autoOffsetReset, setFormData])

  if (!kafkaClusters || getKafkaClustersStatus !== 'success') return

  if (kafkaClusters.length <= 0) {
    return (
      <p>
        No Kafka clusters were found in this space. You'll have to&nbsp;
        <span className="text-blue-600" onClick={() => openModal('createThing')}>
          create an Kafka cluster
        </span>
        &nbsp;in this space before creating an Kafka consumer.
      </p>
    )
  }

  return (
    <>
      <FormFieldSelect
        name="cluster"
        label="Cluster"
        className="mt-4 mb-2"
        options={kafkaClusterOptions}
        value={clusterId}
        onChange={(e) => setClusterId(e.target.value)}
      />
      <FormFieldText
        name="groupId"
        label="Group ID"
        className="mt-4 mb-2"
        value={groupId}
        onChange={(e) => setGroupId(e.target.value)}
      />
      <FormFieldSelect
        name="autoOffsetReset"
        label="Auto Offset Reset"
        className="mt-4 mb-2"
        options={[
          {
            label: 'latest',
            value: 'latest',
          },
          {
            label: 'earliest',
            value: 'earliest',
          },
          {
            label: 'none',
            value: 'none',
          },
        ]}
        value={autoOffsetReset}
        onChange={(e) => setAutoOffsetReset(e.target.value as 'latest' | 'earliest' | 'none')}
      />
    </>
  )
}

export default KafkaConsumer
