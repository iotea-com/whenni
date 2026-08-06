'use client'

import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import { Dispatch, FC, useEffect, useMemo, useState } from 'react'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import { RemixIcon, riCloseLine } from '@mwarnerdotme/react-remixicon'
import { Thing } from '@prisma/client'

type Props = {
  setFormData: Dispatch<Map<string, string | number | string[] | number[]>>
  initial?: Partial<Thing>
}

const HttpServer: FC<Props> = ({ setFormData, initial }) => {
  const attributes = useMemo(() => {
    if (initial && initial.attributes && typeof initial.attributes === 'string')
      return JSON.parse(initial.attributes)

    if (initial && initial.attributes && typeof initial.attributes === 'object')
      return initial.attributes

    return undefined
  }, [initial])

  const [protocol, setProtocol] = useState<'http' | 'https'>(attributes?.['protocol'] ?? 'https')
  const [host, setHost] = useState<string>(attributes?.host)
  const [port, setPort] = useState<number>(attributes?.port ?? 443)
  const [paths, setPaths] = useState<string[]>(attributes?.paths ?? [''])

  useEffect(() => {
    if (protocol === 'http') setPort(80)
    if (protocol === 'https') setPort(443)
  }, [protocol])

  useEffect(() => {
    const updatedFormData = new Map()
    updatedFormData.set('protocol', protocol)
    updatedFormData.set('host', host)
    updatedFormData.set('port', port)
    updatedFormData.set('paths', paths)
    setFormData(updatedFormData)
  }, [protocol, host, port, paths, setFormData])

  return (
    <>
      <FormFieldSelect
        name="protocol"
        label="Protocol"
        options={[
          {
            label: 'http',
            value: 'http',
          },
          {
            label: 'https',
            value: 'https',
          },
        ]}
        className="mt-4 mb-2"
        value={protocol}
        onChange={(e) => setProtocol(e.target.value as 'http' | 'https')}
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
        inputType="number"
        className="mt-4 mb-2"
        value={port.toString()}
        onChange={(e) => setPort(parseInt(e.target.value))}
      />
      {paths.map((path, i) => {
        return (
          <div key={`path-${i}`} className="flex mt-4 mb-2 items-center">
            <FormFieldText
              name={`path-${i}`}
              label={`Path ${i + 1}`}
              className="grow"
              value={path}
              onChange={(e) =>
                setPaths((current) => {
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
                setPaths((current) => current.filter((_p, j) => i !== j))
              }}
            />
          </div>
        )
      })}
      <div className="flex">
        <Button
          text="Add Path"
          className="mb-4"
          onClick={() => setPaths((current) => [...current, ''])}
        />
      </div>
    </>
  )
}

export default HttpServer
