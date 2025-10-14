'use client'

import FormFieldSelect from '@iotea/libs/frontend/components/atoms/FormFieldSelect'
import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import { Thing } from '@prisma/client'
import { Dispatch, FC, useEffect, useMemo, useState } from 'react'

type Props = {
  setFormData: Dispatch<Map<string, string | number | string[] | number[]>>
  secretOptions: { label: string; value: string }[]
  initial?: Partial<Thing>
}

const MinioBucket: FC<Props> = ({ setFormData, secretOptions, initial }) => {
  const attributes = useMemo(() => {
    if (initial && initial.attributes && typeof initial.attributes === 'string')
      return JSON.parse(initial.attributes)

    if (initial && initial.attributes && typeof initial.attributes === 'object')
      return initial.attributes

    return undefined
  }, [initial])

  const [MinioEndpoint, setMinioEndpoint] = useState<string>(attributes?.MinioEndpoint ?? '')
  const [MinioAccessKeyId, setMinioAccessKeyId] = useState<string>(
    attributes?.MinioAccessKeyId ?? '',
  )
  const [MinioSecretAccessKey, setMinioSecretAccessKey] = useState<string>(() => {
    if (secretOptions.length === 1 && !attributes?.MinioSecretAccessKey) {
      return secretOptions[0].value
    }
    return attributes?.MinioSecretAccessKey ?? ''
  })
  const [bucketName, setbucketName] = useState<string>(attributes?.bucketName ?? '')
  const [UseSSL, setUseSSL] = useState<boolean>(attributes?.UseSSL ?? true)

  useEffect(() => {
    const updatedFormData = new Map()
    updatedFormData.set('MinioEndpoint', MinioEndpoint)
    updatedFormData.set('MinioAccessKeyId', MinioAccessKeyId)
    updatedFormData.set('MinioSecretAccessKey', MinioSecretAccessKey)
    updatedFormData.set('bucketName', bucketName)
    updatedFormData.set('UseSSL', UseSSL)
    setFormData(updatedFormData)
  }, [MinioEndpoint, MinioAccessKeyId, MinioSecretAccessKey, bucketName, UseSSL, setFormData])

  return (
    <>
      <FormFieldText
        name="MinioEndpoint"
        label="Minio Endpoint"
        className="mt-4 mb-2"
        value={MinioEndpoint}
        onChange={(e) => setMinioEndpoint(e.target.value)}
      />
      <FormFieldText
        name="MinioAccessKeyId"
        label="Minio Access Key ID"
        className="mt-4 mb-2"
        value={MinioAccessKeyId}
        onChange={(e) => setMinioAccessKeyId(e.target.value)}
      />
      <FormFieldSelect
        name="MinioSecretAccessKey"
        label="Minio Secret Access Key"
        className="mt-4 mb-2"
        value={MinioSecretAccessKey}
        options={secretOptions}
        onChange={(e) => setMinioSecretAccessKey(e.target.value)}
      />
      <FormFieldText
        name="bucketName"
        label="Bucket Name"
        className="mt-4 mb-2"
        value={bucketName}
        onChange={(e) => setbucketName(e.target.value)}
      />
      <FormFieldSelect
        name="UseSSL"
        label="Use SSL"
        className="mt-4 mb-2"
        options={[
          {
            label: 'True',
            value: 'true', // Changed to a string representation of true/false
          },
          {
            label: 'False',
            value: 'false',
          },
        ]}
        value={UseSSL ? 'true' : 'false'} // Convert boolean to string
        onChange={
          (e) => setUseSSL(e.target.value === 'true') // Convert string back to boolean
        }
      />
    </>
  )
}

export default MinioBucket
