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

const KafkaCluster: FC<Props> = ({ setFormData, initial }) => {
  const attributes = useMemo(() => {
    if (initial && initial.attributes && typeof initial.attributes === 'string')
      return JSON.parse(initial.attributes)

    if (initial && initial.attributes && typeof initial.attributes === 'object')
      return initial.attributes

    return undefined
  }, [initial])

  const [bootstrapServers, setBootstrapServers] = useState<string[]>(
    attributes?.bootstrapServers ?? [],
  )
  const [topics, setTopics] = useState<string[]>(attributes?.topics ?? [])

  useEffect(() => {
    const updatedFormData = new Map()
    updatedFormData.set('bootstrapServers', bootstrapServers)
    updatedFormData.set('topics', topics)
    setFormData(updatedFormData)
  }, [bootstrapServers, topics, setFormData])

  return (
    <>
      {bootstrapServers.map((bootstrapServer, i) => {
        return (
          <div key={`bootstrapServer-${i}`} className="flex mt-4 mb-2 items-center">
            <FormFieldText
              name={`bootstrapServer-${i}`}
              label={`Bootstrap Server ${i + 1}`}
              className="grow"
              value={bootstrapServer}
              onChange={(e) =>
                setBootstrapServers((current) => {
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
                setBootstrapServers((current) => current.filter((_p, j) => i !== j))
              }}
            />
          </div>
        )
      })}
      <div className="flex">
        <Button
          text="Add Bootstrap Server"
          className="mb-4"
          onClick={() => setBootstrapServers((current) => [...current, ''])}
        />
      </div>
      {topics.map((topic, i) => {
        return (
          <div key={`topic-${i}`} className="flex mt-4 mb-2 items-center">
            <FormFieldText
              name={`topic-${i}`}
              label={`Path ${i + 1}`}
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

export default KafkaCluster
