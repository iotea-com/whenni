import useChannelEditorStore from '@gruent/hub/stores/channelEditorStore'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import { ChannelNode } from '@gruent/libs/engine/nodes/v1'
import { MessageQueueSourceNodeConfig } from '@gruent/libs/engine/nodes/v1/src/source/messageQueue'
import { Thing } from '@prisma/client'
import { FC, useEffect, useMemo, useState } from 'react'
import KafkaActionOptions from './KafkaSourceOptions'
import NatsSourceOptions from './NatsSourceOptions'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import { RemixIcon, riAddLine } from '@mwarnerdotme/react-remixicon'

type Props = {
  things: Thing[]
  orgId: string
  spaceId: string
}

export enum MessageQueueSourceClientType {
  Kafka = 'Kafka',
  NATS = 'NATS',
}

const MessageQueueSourceOptions: FC<Props> = ({ things, orgId, spaceId }) => {
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<MessageQueueSourceNodeConfig>,
  )

  const [selectedMessageQueueClient, setSelectedMessageQueueClient] = useState<Thing>()

  // load in the default message queue client from the current node config
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!things) return

    if (Object.keys(currentConfig).includes('kafkaConsumer::thing')) {
      const kafkaConsumer = things.find((s) => {
        if (s.id === currentConfig['kafkaConsumer::thing']) return true
        return false
      })
      if (kafkaConsumer) setSelectedMessageQueueClient(kafkaConsumer)
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

    if (selectedMessageQueueClient.thingCategory === 'KAFKA_CONSUMER')
      return MessageQueueSourceClientType.Kafka

    if (selectedMessageQueueClient.thingCategory === 'NATS_CLIENT')
      return MessageQueueSourceClientType.NATS
  }, [selectedMessageQueueClient])

  const messageQueueClients = useMemo(() => {
    if (!things) return []

    const kafkaConsumers = things.filter((thing) => {
      if (thing.thingCategory == 'KAFKA_CONSUMER') return true
      return false
    })

    const natsClients = things.filter((thing) => {
      if (thing.thingCategory == 'NATS_CLIENT') return true
      return false
    })

    return [...kafkaConsumers, ...natsClients]
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
      case MessageQueueSourceClientType.Kafka:
        return (
          <KafkaActionOptions
            things={things}
            orgId={orgId}
            spaceId={spaceId}
            messageQueueClientOptions={messageQueueClientOptions}
            selectedMessageQueueClient={selectedMessageQueueClient}
          />
        )
      case MessageQueueSourceClientType.NATS:
        return (
          <NatsSourceOptions
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
        <Button variant="underline" onClick={() => openModal('createThing')}>
          Add a message queue client
        </Button>{' '}
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
          value={selectedMessageQueueClient?.id ?? '__GRUENT_IGNORE__'}
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

export default MessageQueueSourceOptions
