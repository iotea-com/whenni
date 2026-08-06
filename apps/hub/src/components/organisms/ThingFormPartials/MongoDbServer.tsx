'use client'

import Button from '@gruent/libs/frontend/components/atoms/Button'
import FormFieldSelect from '@gruent/libs/frontend/components/atoms/FormFieldSelect'
import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import { RemixIcon, riCloseLine } from '@mwarnerdotme/react-remixicon'
import { Thing } from '@prisma/client'
import { Dispatch, FC, useEffect, useMemo, useState } from 'react'

type Props = {
  setFormData: Dispatch<Map<string, string | number | string[] | number[]>>
  secretOptions: { label: string; value: string }[]
  initial?: Partial<Thing>
}

const MongoDbServer: FC<Props> = ({ setFormData, secretOptions, initial }) => {
  const attributes = useMemo(() => {
    if (initial && initial.attributes && typeof initial.attributes === 'string')
      return JSON.parse(initial.attributes)

    if (initial && initial.attributes && typeof initial.attributes === 'object')
      return initial.attributes

    return undefined
  }, [initial])

  const [protocol, setProtocol] = useState<'mongodb' | 'mongodb+srv'>(
    attributes?.protocol ?? 'mongodb',
  )
  const [host, setHost] = useState<string>(attributes?.host ?? '')
  const [port, setPort] = useState<number>(attributes?.port ?? 27017)
  const [username, setUsername] = useState<string>(attributes?.username ?? '')
  const [password, setPassword] = useState<string>(attributes?.password ?? '')
  const [databases, setDatabases] = useState<string[]>(attributes?.databases ?? [])
  const [collections, setCollections] = useState<Record<string, string[] | undefined>>(
    attributes?.collections ?? {},
  )

  useEffect(() => {
    const updatedFormData = new Map()
    updatedFormData.set('protocol', protocol)
    updatedFormData.set('host', host)
    if (protocol === 'mongodb') updatedFormData.set('port', port)
    updatedFormData.set('username', username)
    updatedFormData.set('password', password)
    updatedFormData.set('databases', databases)
    updatedFormData.set('collections', collections)
    setFormData(updatedFormData)
  }, [protocol, host, port, username, password, databases, collections, setFormData])

  return (
    <>
      <FormFieldSelect
        name="protocol"
        label="Protocol"
        options={[
          {
            label: 'mongodb',
            value: 'mongodb',
          },
          {
            label: 'mongodb+srv',
            value: 'mongodb+srv',
          },
        ]}
        className="mt-4 mb-2"
        value={protocol}
        onChange={(e) => setProtocol(e.target.value as 'mongodb' | 'mongodb+srv')}
      />
      <FormFieldText
        name="host"
        label="Host"
        className="mt-4 mb-2"
        value={host}
        onChange={(e) => setHost(e.target.value)}
      />
      {protocol === 'mongodb' && (
        <FormFieldText
          name="port"
          inputType="number"
          min="1"
          max="65535"
          label="Port"
          className="mt-4 mb-2"
          value={port.toString()}
          onChange={(e) => setPort(parseInt(e.target.value))}
        />
      )}
      <FormFieldText
        name="username"
        label="Username"
        className="mt-4 mb-2"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <FormFieldSelect
        name="password"
        label="Password"
        className="mt-4 mb-2"
        value={password}
        options={secretOptions}
        onChange={(e) => setPassword(e.target.value)}
      />
      {databases.map((database, i) => {
        return (
          <div key={`database-${i}`}>
            <div className="flex mt-4 items-center">
              <FormFieldText
                name={`database-${i}`}
                label={`Database ${i + 1}`}
                className="grow"
                value={database}
                onChange={(e) =>
                  setDatabases((current) => {
                    const updated = [...current]
                    updated[i] = e.target.value
                    return updated
                  })
                }
              />
              <RemixIcon
                icon={riCloseLine}
                className="-mr-4 text-gray-500 hover:text-red-500 transition cursor-pointer"
                onClick={() => {
                  setDatabases((current) => current.filter((_p, j) => i !== j))
                  setCollections((current) => {
                    const updated = structuredClone(current)
                    updated[database] = undefined
                    return updated
                  })
                }}
              />
            </div>
            {collections[database] &&
              collections[database].map((collection, j) => {
                return (
                  <div key={`collection-${j}`} className="flex mt-0 mb-0 ml-4 items-center">
                    <FormFieldText
                      name={`collection-${j}`}
                      label={`Collection ${j + 1}`}
                      className="grow"
                      value={collection}
                      onChange={(e) =>
                        setCollections((current) => {
                          const updated = structuredClone(current)
                          if (updated[database]) updated[database][j] = e.target.value
                          return updated
                        })
                      }
                    />
                    <RemixIcon
                      icon={riCloseLine}
                      className="-mr-4 text-gray-500 hover:text-red-500 transition cursor-pointer"
                      onClick={() => {
                        setCollections((current) => {
                          const updated = structuredClone(current)
                          updated[database] = updated[database]?.filter((_p, k) => j !== k) || []
                          return updated
                        })
                      }}
                    />
                  </div>
                )
              })}
            <Button
              text="Add Collection"
              className="mb-1 ml-4"
              onClick={() =>
                setCollections((current) => {
                  const updated = structuredClone(current)
                  if (!updated[database]) {
                    updated[database] = ['']
                  } else {
                    updated[database] = [...updated[database], '']
                  }

                  return updated
                })
              }
            />
          </div>
        )
      })}
      <div className="flex">
        <Button
          text="Add Database"
          className="mt-2 mb-4"
          onClick={() => setDatabases((current) => [...current, `database${databases.length + 1}`])}
        />
      </div>
    </>
  )
}

export default MongoDbServer
