import useChannelEditorStore from '@iotea/hub/stores/channelEditorStore'
import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import { ChannelNode } from '@iotea/libs/engine/nodes/v1'
import { MessageQueueSourceNodeConfig } from '@iotea/libs/engine/nodes/v1/src/source/messageQueue'
import { Thing } from '@prisma/client'
import { FC, useEffect, useState } from 'react'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import { openModal } from '@iotea/libs/frontend/hooks/useModal'

type Props = {
  things: Thing[]
  orgId: string
  spaceId: string
  messageQueueClientOptions: {
    label: string
    value: string
  }[]
  selectedMessageQueueClient?: Thing
}

const KafkaSourceOptions: FC<Props> = ({
  things,
  messageQueueClientOptions,
  selectedMessageQueueClient,
}) => {
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<MessageQueueSourceNodeConfig>,
  )

  const [topic, setTopic] = useState<string>('')

  // load in default values
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!things) return

    setTopic(currentConfig?.topic || '')
  }, [currentNode, things])

  // update config and dependencies when options have been updated
  useEffect(() => {
    if (!selectedMessageQueueClient) return

    // update config
    const updatedNode = structuredClone(currentNode)

    updatedNode.metadata.config = {
      'kafkaConsumer::thing': selectedMessageQueueClient.id,
      topic,
    }

    // update things dependency array
    updatedNode.metadata.dependencies.things = [
      {
        fieldName: 'kafkaConsumer::thing',
        internal: false,
        thingId: selectedMessageQueueClient.id,
      },
    ]

    // apply updates to the nodes map
    upsertNode(currentNode.id, updatedNode)

    // eslint-disable-next-line -- currentNode in the dependencies array causes node configs to sometimes be copied
  }, [topic, selectedMessageQueueClient, upsertNode])

  if (!messageQueueClientOptions || messageQueueClientOptions.length <= 0)
    return (
      <p>
        <Button variant="underline" onClick={() => openModal('createThing')}>
          Add a message queue client
        </Button>{' '}
        to your space before using this node.
      </p>
    )

  return (
    <>
      <FormFieldText
        name="topic"
        label="Topic"
        value={topic}
        onChange={(e) => setTopic(e.target.value)}
      />
    </>
  )
}

export default KafkaSourceOptions
