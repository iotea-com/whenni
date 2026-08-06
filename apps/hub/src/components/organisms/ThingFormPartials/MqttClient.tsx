'use client'

import useAuth from '@gruent/hub/hooks/useAuth'
import gruentClient from '@gruent/hub/lib/gruent'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import { openModal } from '@gruent/libs/frontend/hooks/useModal'
import { Thing } from '@prisma/client'
import { useQuery } from '@tanstack/react-query'
import { Dispatch, FC, useEffect, useMemo, useState } from 'react'

type Props = {
  setFormData: Dispatch<Map<string, string | number | string[] | number[]>>
  spaceId: string
  orgId: string
  secretOptions: { label: string; value: string }[]
  initial?: Partial<Thing>
}

const MqttClient: FC<Props> = ({ setFormData, spaceId, secretOptions, initial }) => {
  const attributes = useMemo(() => {
    if (initial && initial.attributes && typeof initial.attributes === 'string')
      return JSON.parse(initial.attributes)

    if (initial && initial.attributes && typeof initial.attributes === 'object')
      return initial.attributes

    return undefined
  }, [initial])

  const [broker, setBroker] = useState<string>(attributes?.broker)
  const [username, setUsername] = useState<string>(attributes?.username)
  const [password, setPassword] = useState<string>(() => {
    if (secretOptions.length === 1 && !attributes?.password) {
      return secretOptions[0].value
    }
    return attributes?.password ?? ''
  })
  const [clientId, setClientId] = useState<string>(attributes?.clientId)
  const [certificateId, _setCertificateId] = useState<string>(attributes?.certificateId)
  const [keypair, _setKeypair] = useState<string>(attributes?.keypair)

  const { accessToken } = useAuth()

  const {
    status: getMqttBrokersStatus,
    data: mqttBrokers,
    error: getMqttBrokersError,
  } = useQuery({
    queryKey: ['things', spaceId],
    queryFn: async () => {
      if (!accessToken) throw new Error('Invalid auth session.')

      const { data: things, errors } = await gruentClient(accessToken).things.list(spaceId, {
        category: 'MQTT_BROKER',
      })

      if (errors && errors.length > 0) throw new Error(errors[0])

      return things
    },
  })

  const mqttBrokerOptions = useMemo(() => {
    if (getMqttBrokersError) return []
    if (getMqttBrokersStatus !== 'success') return []
    if (!mqttBrokers) return []

    return mqttBrokers.map((mqttBroker) => {
      return {
        label: mqttBroker.name,
        value: mqttBroker.id,
      }
    })
  }, [mqttBrokers, getMqttBrokersError, getMqttBrokersStatus])

  useEffect(() => {
    const updatedFormData = new Map()
    updatedFormData.set('broker', broker)
    updatedFormData.set('username', username)
    updatedFormData.set('password', password)
    updatedFormData.set('clientId', clientId)
    updatedFormData.set('certificateId', certificateId)
    updatedFormData.set('keypair', keypair)
    setFormData(updatedFormData)
  }, [broker, username, password, clientId, certificateId, keypair, setFormData])

  if (!mqttBrokers || getMqttBrokersStatus !== 'success') return

  if (mqttBrokers.length <= 0) {
    return (
      <p>
        No MQTT brokers were found in this space. You'll have to&nbsp;
        <Button variant="underline" onClick={() => openModal('createThing')}>
          create an MQTT broker
        </Button>
        &nbsp;in this space before creating an MQTT client.
      </p>
    )
  }

  return (
    <>
      <FormFieldSelect
        name="broker"
        label="Broker"
        className="mt-4 mb-2"
        options={mqttBrokerOptions}
        value={broker}
        onChange={(e) => setBroker(e.target.value)}
      />
      <FormFieldText
        name="username"
        label="Username (optional)"
        className="mt-4 mb-2"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <FormFieldSelect
        name="password"
        label="Password (optional)"
        className="mt-4 mb-2"
        value={password}
        options={secretOptions}
        onChange={(e) => setPassword(e.target.value)}
      />
      <FormFieldText
        name="clientId"
        label="Client ID (optional)"
        className="mt-4 mb-2"
        value={clientId}
        onChange={(e) => setClientId(e.target.value)}
      />
    </>
  )
}

export default MqttClient
