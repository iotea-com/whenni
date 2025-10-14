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

const AwsSNSEndpoint: FC<Props> = ({ setFormData, secretOptions, initial }) => {
  const attributes = useMemo(() => {
    if (initial && initial.attributes && typeof initial.attributes === 'string')
      return JSON.parse(initial.attributes)

    if (initial && initial.attributes && typeof initial.attributes === 'object')
      return initial.attributes

    return undefined
  }, [initial])

  const [AwsAccessKeyId, setAwsAccessKeyId] = useState<string>(attributes?.AwsAccessKeyId ?? '')
  const [AwsSecretAccessKey, setAwsSecretAccessKey] = useState<string>(() => {
    if (secretOptions.length === 1 && !attributes?.AwsSecretAccessKey) {
      return secretOptions[0].value
    }
    return attributes?.AwsSecretAccessKey ?? ''
  })
  const [AwsRegion, setAwsRegion] = useState<string>(attributes?.AwsRegion ?? '')

  useEffect(() => {
    const updatedFormData = new Map()
    updatedFormData.set('AwsAccessKeyId', AwsAccessKeyId)
    updatedFormData.set('AwsSecretAccessKey', AwsSecretAccessKey)
    updatedFormData.set('AwsRegion', AwsRegion)
    setFormData(updatedFormData)
  }, [AwsAccessKeyId, AwsSecretAccessKey, AwsRegion, setFormData])

  return (
    <>
      <FormFieldText
        name="AwsAccessKeyId"
        label="AWS Access Key ID"
        className="mt-4 mb-2"
        value={AwsAccessKeyId}
        onChange={(e) => setAwsAccessKeyId(e.target.value)}
      />
      <FormFieldSelect
        name="AwsSecretAccessKey"
        label="AWS Secret Access Key"
        className="mt-4 mb-2"
        value={AwsSecretAccessKey}
        options={secretOptions}
        onChange={(e) => {
          setAwsSecretAccessKey(e.target.value)
        }}
      />
      <FormFieldText
        name="AwsRegion"
        label="AWS Region"
        className="mt-4 mb-2"
        value={AwsRegion}
        onChange={(e) => setAwsRegion(e.target.value)}
      />
    </>
  )
}

export default AwsSNSEndpoint
