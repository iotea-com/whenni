'use client'

import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import { Dispatch, FC, useEffect, useMemo, useState } from 'react'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import { RemixIcon, riCloseLine } from '@mwarnerdotme/react-remixicon'
import { Thing } from '@prisma/client'

type Props = {
  setFormData: Dispatch<Map<string, string | number | string[] | number[]>>
  initial?: Partial<Thing>
}

const NatsServer: FC<Props> = ({ setFormData, initial }) => {
  const attributes = useMemo(() => {
    if (initial && initial.attributes && typeof initial.attributes === 'string')
      return JSON.parse(initial.attributes)

    if (initial && initial.attributes && typeof initial.attributes === 'object')
      return initial.attributes

    return undefined
  }, [initial])

  const [host, setHost] = useState<string>(attributes?.host ?? '')
  const [port, setPort] = useState<number>(attributes?.port ?? 4222)
  const [topics, setTopics] = useState<string[]>(attributes?.topics ?? [])

  useEffect(() => {
    const updatedFormData = new Map()
    updatedFormData.set('host', host)
    updatedFormData.set('port', port)
    updatedFormData.set('topics', topics)
    setFormData(updatedFormData)
  }, [host, port, topics, setFormData])

  return (
    <>
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
        inputType="number"
        className="mt-4 mb-2"
        value={port.toString()}
        onChange={(e) => setPort(parseInt(e.target.value))}
      />
      {topics.map((topic, i) => {
        return (
          <div key={`topic-${i}`} className="flex mt-4 mb-2 items-center">
            <FormFieldText
              name={`topic-${i}`}
              label={`Topic ${i + 1}`}
              className="grow"
              value={topic}
              onChange={(e) =>
                setTopics((current) => {
                  const n = [...current]
                  n[i] = e.target.value
                  return n
                })
              }
            />
            <RemixIcon
              icon={riCloseLine}
              className="-mr-4 text-gray-500 hover:text-red-500 transition cursor-pointer"
              onClick={() => {
                setTopics((current) => current.filter((_p, j) => i !== j))
              }}
            />
          </div>
        )
      })}
      <div className="flex">
        <Button
          text="Add Topic"
          className="mb-4"
          onClick={() => setTopics((current) => [...current, ''])}
        />
      </div>
    </>
  )
}

export default NatsServer
