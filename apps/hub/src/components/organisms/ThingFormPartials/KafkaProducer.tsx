'use client'

import useAuth from '@iotea/hub/hooks/useAuth'
import ioteaClient from '@iotea/hub/lib/iotea'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import FormFieldSelect from '@iotea/libs/frontend/components/atoms/FormFieldSelect'
import { openModal } from '@iotea/libs/frontend/hooks/useModal'
import { Thing } from '@prisma/client'
import { useQuery } from '@tanstack/react-query'
import { Dispatch, FC, useEffect, useMemo, useState } from 'react'

type Props = {
  setFormData: Dispatch<Map<string, string | number | string[] | number[]>>
  spaceId: string
  orgId: string
  initial?: Partial<Thing>
}

const KafkaProducer: FC<Props> = ({ setFormData, spaceId, initial }) => {
  const attributes = useMemo(() => {
    if (initial && initial.attributes && typeof initial.attributes === 'string')
      return JSON.parse(initial.attributes)

    if (initial && initial.attributes && typeof initial.attributes === 'object')
      return initial.attributes

    return undefined
  }, [initial])

  const [cluster, setCluster] = useState<string>(attributes?.cluster)

  const { accessToken } = useAuth()

  const {
    status: getKafkaClustersStatus,
    data: kafkaClusters,
    error: getKafkaClustersError,
  } = useQuery({
    queryKey: ['things', 'kafkaClusters', spaceId],
    queryFn: async () => {
      if (!accessToken) throw new Error('Invalid auth session.')

      const { data: things, errors } = await ioteaClient(accessToken).things.list(spaceId, {
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
    updatedFormData.set('cluster', cluster)
    setFormData(updatedFormData)
  }, [cluster, setFormData])

  if (!kafkaClusters || getKafkaClustersStatus !== 'success') return

  if (kafkaClusters.length <= 0) {
    return (
      <p>
        No Kafka clusters were found in this space. You'll have to&nbsp;
        <Button variant="underline" onClick={() => openModal('createThing')}>
          create a Kafka cluster
        </Button>
        &nbsp;in this space before creating an Kafka producer.
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
        value={cluster}
        onChange={(e) => setCluster(e.target.value)}
      />
    </>
  )
}

export default KafkaProducer
