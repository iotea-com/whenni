import useChannelEditorStore from '@iotea/hub/stores/channelEditorStore'
import FormFieldSelect from '@iotea/libs/frontend/components/atoms/FormFieldSelect'
import { ChannelNode } from '@iotea/libs/engine/nodes/v1'
import { MessageQueueActionNodeConfig } from '@iotea/libs/engine/nodes/v1/src/action/messageQueue'
import { Thing } from '@prisma/client'
import { FC, useEffect, useMemo, useState } from 'react'
import KafkaActionOptions from './KafkaActionOptions'
import NatsActionOptions from './NatsActionOptions'
import { openModal } from '@iotea/libs/frontend/hooks/useModal'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import { RemixIcon, riAddLine } from '@mwarnerdotme/react-remixicon'

type Props = {
  things: Thing[]
  orgId: string
  spaceId: string
}

export enum MessageQueue {
  Kafka = 'Kafka',
  NATS = 'NATS',
}

const MessageQueueActionOptions: FC<Props> = ({ things, orgId, spaceId }) => {
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<MessageQueueActionNodeConfig>,
  )

  const [selectedMessageQueueClient, setSelectedMessageQueueClient] = useState<Thing>()

  // load in the default message queue client from the current node config
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!things) return

    if (Object.keys(currentConfig).includes('kafkaProducer::thing')) {
      const kafkaProducer = things.find((s) => {
        if (s.id === currentConfig['kafkaProducer::thing']) return true
        return false
      })
      if (kafkaProducer) setSelectedMessageQueueClient(kafkaProducer)
      return
    }

    if (Object.keys(currentConfig).includes('natsClient::thing')) {
      const natsClient = things.find((s) => {
        if (s.id === currentConfig['natsClient::thing']) return true
        return false
      })
      if (natsClient) setSelectedMessageQueueClient(natsClient)
      return
    }
  }, [currentNode, things])

  const selectedMessageQueueType = useMemo(() => {
    if (!selectedMessageQueueClient) return

    if (selectedMessageQueueClient.thingCategory === 'KAFKA_PRODUCER') return MessageQueue.Kafka

    if (selectedMessageQueueClient.thingCategory === 'NATS_CLIENT') return MessageQueue.NATS
  }, [selectedMessageQueueClient])

  const messageQueueClients = useMemo(() => {
    if (!things) return []

    const kafkaProducers = things.filter((thing) => {
      if (thing.thingCategory == 'KAFKA_PRODUCER') return true
      return false
    })

    const natsClients = things.filter((thing) => {
      if (thing.thingCategory == 'NATS_CLIENT') return true
      return false
    })

    return [...kafkaProducers, ...natsClients]
  }, [things])

  const messageQueueClientOptions = useMemo(() => {
    return messageQueueClients.map((thing) => ({
      label: thing.name,
      value: thing.id,
    }))
  }, [messageQueueClients])

  const selectedSubnodeOptions = useMemo(() => {
    if (!selectedMessageQueueType) return

    switch (selectedMessageQueueType) {
      case MessageQueue.Kafka:
        return (
          <KafkaActionOptions
            things={things}
            orgId={orgId}
            spaceId={spaceId}
            messageQueueClientOptions={messageQueueClientOptions}
            selectedMessageQueueClient={selectedMessageQueueClient}
          />
        )
      case MessageQueue.NATS:
        return (
          <NatsActionOptions
            things={things}
            orgId={orgId}
            spaceId={spaceId}
            messageQueueClientOptions={messageQueueClientOptions}
            selectedMessageQueueClient={selectedMessageQueueClient}
          />
        )
    }
  }, [
    selectedMessageQueueType,
    messageQueueClientOptions,
    selectedMessageQueueClient,
    orgId,
    spaceId,
    things,
  ])

  if (!messageQueueClients || messageQueueClients.length <= 0)
    return (
      <p className="mb-2">
        <Button variant="underline" className="mr-1" onClick={() => openModal('createThing')}>
          Add a message queue client
        </Button>
        to your space before using this node.
      </p>
    )

  return (
    <>
      <div className="flex flex-row gap-2">
        <FormFieldSelect
          name="messageQueueClient"
          label="Message Queue Client"
          className="grow"
          options={messageQueueClientOptions}
          value={selectedMessageQueueClient?.id ?? '__IOTEA_IGNORE__'}
          onChange={(e) =>
            setSelectedMessageQueueClient(
              messageQueueClients.find((s) => {
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

export default MessageQueueActionOptions
