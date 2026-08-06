'use client'

import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import { Dispatch, FC, useEffect, useMemo, useState } from 'react'
import { Thing } from '@prisma/client'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'

type Props = {
  setFormData: Dispatch<Map<string, string | number | string[] | number[]>>
  secretOptions: { label: string; value: string }[]
  initial?: Partial<Thing>
}

const ClickhouseDatabase: FC<Props> = ({ setFormData, secretOptions, initial }) => {
  const attributes = useMemo(() => {
    if (initial && initial.attributes && typeof initial.attributes === 'string')
      return JSON.parse(initial.attributes)

    if (initial && initial.attributes && typeof initial.attributes === 'object')
      return initial.attributes

    return undefined
  }, [initial])

  const [host, setHost] = useState<string>(attributes?.host)
  const [port, setPort] = useState<number>(attributes?.port ?? 9000)
  const [username, setUsername] = useState<string>(attributes?.username ?? '')
  const [password, setPassword] = useState<string>(() => {
    if (secretOptions.length === 1 && !attributes?.password) {
      return secretOptions[0].value
    }
    return attributes?.password ?? ''
  })
  const [database, setDatabase] = useState<string>(attributes?.database ?? '')

  useEffect(() => {
    const updatedFormData = new Map()
    updatedFormData.set('host', host)
    updatedFormData.set('port', port)
    updatedFormData.set('username', username)
    updatedFormData.set('password', password)
    updatedFormData.set('database', database)
    setFormData(updatedFormData)
  }, [host, port, username, password, database, setFormData])

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
        className="mt-4 mb-2"
        value={port.toString()}
        onChange={(e) => setPort(parseInt(e.target.value))}
      />
      <FormFieldText
        name="username"
        label="Username"
        className="mt-4 mb-2"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <FormFieldSelect
        name="password"
        label="Database Password"
        className="mt-4 mb-2"
        value={password}
        options={secretOptions}
        onChange={(e) => setPassword(e.target.value)}
      />
      <FormFieldText
        name="database"
        label="Database Name"
        className="mt-4 mb-2"
        value={database}
        onChange={(e) => setDatabase(e.target.value)}
      />
    </>
  )
}

export default ClickhouseDatabase
