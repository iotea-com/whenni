import useChannelEditorStore from '@iotea/hub/stores/channelEditorStore'
import FormFieldSelect from '@iotea/libs/frontend/components/atoms/FormFieldSelect'
import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import { ChannelNode, MqttSourceNodeConfig } from '@iotea/libs/engine/nodes/v1'
import { Thing } from '@prisma/client'
import { FC, useEffect, useMemo, useState } from 'react'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import { openModal } from '@iotea/libs/frontend/hooks/useModal'
import { RemixIcon, riAddLine } from '@mwarnerdotme/react-remixicon'

type Props = {
  things: Thing[]
  orgId: string
  spaceId: string
}

const MqttSourceOptions: FC<Props> = ({ things }) => {
  const upsertNode = useChannelEditorStore((state) => state.upsertNode)
  const currentNode = useChannelEditorStore(
    (state) => state.currentNode as ChannelNode<MqttSourceNodeConfig>,
  )

  const [selectedMqttClient, setSelectedMqttClient] = useState<Thing>()
  const [qos, setQos] = useState<string>()
  const [topic, setTopic] = useState<string>()

  // load in default values
  useEffect(() => {
    const currentConfig = { ...currentNode.metadata.config }
    if (!things) return

    const thing = things.find((thing) => {
      if (thing.id === currentConfig['mqttClient::thing']) return true
      return false
    })

    setSelectedMqttClient(thing)
    setQos(currentConfig.qos.toString() || '0')
    setTopic(currentConfig.topic || '')
  }, [currentNode, things])

  const mqttClients = useMemo(() => {
    if (!things) return []

    return things.filter((thing) => {
      if (thing.thingCategory == 'MQTT_CLIENT') return true
      return false
    })
  }, [things])

  // update config and dependencies when options have been updated
  useEffect(() => {
    if (!upsertNode) return
    if (!selectedMqttClient) return
    if (!qos) return
    if (!topic) return

    // update config
    const updatedNode = structuredClone(currentNode)
    updatedNode.metadata.config = {
      'mqttClient::thing': selectedMqttClient.id,
      qos: parseInt(qos) as 0 | 1 | 2,
      topic,
    }

    // update dependencies array
    let mqttClientThingDependencyIndex = updatedNode.metadata.dependencies.things.findIndex(
      (d) => d.fieldName && d.fieldName === 'mqttClient::thing',
    )
    if (mqttClientThingDependencyIndex < 0)
      mqttClientThingDependencyIndex = updatedNode.metadata.dependencies.things.length

    updatedNode.metadata.dependencies.things[mqttClientThingDependencyIndex] = {
      fieldName: 'mqttClient::thing',
      internal: false,
      thingId: selectedMqttClient.id,
    }

    upsertNode(currentNode.id, updatedNode)

    // eslint-disable-next-line -- currentNode in the dependencies array causes node configs to sometimes be copied
  }, [upsertNode, selectedMqttClient, qos, topic])

  const mqttClientOptions = useMemo(() => {
    return mqttClients.map((mqttClient) => ({
      label: mqttClient.name,
      value: mqttClient.id,
    }))
  }, [mqttClients])

  if (!mqttClients || mqttClients.length <= 0)
    return (
      <p className="mb-2">
        <Button variant="underline" onClick={() => openModal('createThing')}>
          Add an MQTT client
        </Button>{' '}
        to your space before using this node.
      </p>
    )

  return (
    <>
      <div className="flex flex-row gap-2">
        <FormFieldSelect
          name="mqttClient::thing"
          label="MQTT Client"
          className="grow"
          options={mqttClientOptions}
          value={selectedMqttClient?.id ?? '__IOTEA_IGNORE__'}
          onChange={(e) =>
            setSelectedMqttClient(
              mqttClients.find((c) => {
                if (c.id === e.target.value) return true
                return false
              }),
            )
          }
        />
        <Button text="Create MQTT client" className="my-3" onClick={() => openModal('createThing')}>
          <RemixIcon icon={riAddLine} />
        </Button>
      </div>
      <FormFieldText
        name="topic"
        label="Topic"
        value={topic ?? ''}
        onChange={(e) => setTopic(e.target.value)}
      />
      <FormFieldSelect
        name="qos"
        label="QoS"
        options={[
          {
            label: '0 - At most once',
            value: '0',
          },
          {
            label: '1 - At least once',
            value: '1',
          },
          {
            label: '2 - Exactly once',
            value: '2',
          },
        ]}
        value={qos ?? '0'}
        onChange={(e) => setQos(e.target.value)}
      />
    </>
  )
}

export default MqttSourceOptions
