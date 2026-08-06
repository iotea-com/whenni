'use client'

import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import { Thing } from '@prisma/client'
import { Dispatch, FC, useEffect, useMemo, useState } from 'react'

type Props = {
  setFormData: Dispatch<Map<string, string | number | string[] | number[]>>
  initial?: Partial<Thing>
}

const MqttBroker: FC<Props> = ({ setFormData, initial }) => {
  const attributes = useMemo(() => {
    if (initial && initial.attributes && typeof initial.attributes === 'string')
      return JSON.parse(initial.attributes)

    if (initial && initial.attributes && typeof initial.attributes === 'object')
      return initial.attributes

    return undefined
  }, [initial])

  const [protocol, setProtocol] = useState<'mqtt' | 'mqtts'>(attributes?.protocol ?? 'mqtt')
  const [host, setHost] = useState<string>(attributes?.host)
  const [port, setPort] = useState<number>(attributes?.port ?? 1883)

  useEffect(() => {
    if (protocol === 'mqtt') setPort(1883)
    if (protocol === 'mqtts') setPort(8883)
  }, [protocol])

  useEffect(() => {
    const updatedFormData = new Map()
    updatedFormData.set('protocol', protocol)
    updatedFormData.set('host', host)
    updatedFormData.set('port', port)
    setFormData(updatedFormData)
  }, [protocol, host, port, setFormData])

  return (
    <>
      <FormFieldSelect
        name="protocol"
        label="Protocol"
        options={[
          {
            label: 'mqtt',
            value: 'mqtt',
          },
          {
            label: 'mqtts',
            value: 'mqtts',
          },
        ]}
        className="mt-4 mb-2"
        value={protocol}
        onChange={(e) => setProtocol(e.target.value as 'mqtt' | 'mqtts')}
      />
      <FormFieldText
        name="host"
        label="Host"
        className="mt-4 mb-2"
        value={host}
        onChange={(e) => setHost(e.target.value)}
      />
      <FormFieldText
        name="port"
        label="Port"
        className="mt-4 mb-2"
        value={port.toString()}
        onChange={(e) => setPort(parseInt(e.target.value))}
      />
    </>
  )
}

export default MqttBroker
