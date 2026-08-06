'use client'

import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import { Dispatch, FC, useEffect, useMemo, useState } from 'react'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import { RemixIcon, riCloseLine } from '@mwarnerdotme/react-remixicon'
import { Thing } from '@prisma/client'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'

type Props = {
  setFormData: Dispatch<Map<string, string | number | string[] | number[]>>
  secretOptions: { label: string; value: string }[]
  initial?: Partial<Thing>
}

const InfluxDbDatabase: FC<Props> = ({ setFormData, secretOptions, initial }) => {
  const attributes = useMemo(() => {
    if (initial && initial.attributes && typeof initial.attributes === 'string')
      return JSON.parse(initial.attributes)

    if (initial && initial.attributes && typeof initial.attributes === 'object')
      return initial.attributes

    return undefined
  }, [initial])

  const [protocol, setProtocol] = useState<string>(attributes?.protocol ?? 'https')
  const [host, setHost] = useState<string>(attributes?.host)
  const [port, setPort] = useState<number>(attributes?.port ?? 8086)
  const [token, setToken] = useState<string>(() => {
    if (secretOptions.length === 1 && !attributes?.token) {
      return secretOptions[0].value
    }
    return attributes?.token ?? ''
  })
  const [org, setOrg] = useState<string>(attributes?.org)
  const [buckets, setBuckets] = useState<string[]>(attributes?.buckets ?? [''])

  useEffect(() => {
    const updatedFormData = new Map()
    updatedFormData.set('protocol', protocol)
    updatedFormData.set('host', host)
    updatedFormData.set('port', port)
    updatedFormData.set('token', token)
    updatedFormData.set('org', org)
    updatedFormData.set('buckets', buckets)
    setFormData(updatedFormData)
  }, [protocol, host, port, token, org, buckets, setFormData])

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
        className="mt-4 mb-2"
        value={port.toString()}
        onChange={(e) => setPort(parseInt(e.target.value))}
      />
      <FormFieldSelect
        name="token"
        label="Database Token"
        className="mt-4 mb-2"
        value={token}
        options={secretOptions}
        onChange={(e) => setToken(e.target.value)}
      />
      <FormFieldText
        name="org"
        label="Organization Name"
        className="mt-4 mb-2"
        value={org}
        onChange={(e) => setOrg(e.target.value)}
      />
      {buckets.map((path, i) => {
        return (
          <div key={`bucket-${i}`} className="flex mt-4 mb-2 items-center">
            <FormFieldText
              name={`bucket-${i}`}
              label={`Bucket ${i + 1}`}
              className="grow"
              value={path}
              onChange={(e) =>
                setBuckets((current) => {
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
                setBuckets((current) => current.filter((_p, j) => i !== j))
              }}
            />
          </div>
        )
      })}
      <div className="flex">
        <Button
          text="Add Bucket"
          className="mb-4"
          onClick={() => setBuckets((current) => [...current, ''])}
        />
      </div>
    </>
  )
}

export default InfluxDbDatabase
