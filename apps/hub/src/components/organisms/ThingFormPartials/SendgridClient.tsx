'use client'

import FormFieldSelect from '@iotea/libs/frontend/components/atoms/FormFieldSelect'
import { Thing } from '@prisma/client'
import { Dispatch, FC, useEffect, useMemo, useState } from 'react'

type Props = {
  setFormData: Dispatch<Map<string, string | number | string[] | number[]>>
  secretOptions: { label: string; value: string }[]
  initial?: Partial<Thing>
}

const SendgridClient: FC<Props> = ({ setFormData, secretOptions, initial }) => {
  const attributes = useMemo(() => {
    if (initial && initial.attributes && typeof initial.attributes === 'string')
      return JSON.parse(initial.attributes)

    if (initial && initial.attributes && typeof initial.attributes === 'object')
      return initial.attributes

    return undefined
  }, [initial])

  const [sendgridApiKey, setSendgridApiKey] = useState<string>(() => {
    if (secretOptions.length === 1 && !attributes?.apiKey) {
      return secretOptions[0].value
    }
    return attributes?.apiKey ?? ''
  })

  useEffect(() => {
    const updatedFormData = new Map()
    updatedFormData.set('apiKey', sendgridApiKey)
    setFormData(updatedFormData)
  }, [sendgridApiKey, setFormData])

  return (
    <>
      <FormFieldSelect
        name="apiKey"
        label="Sendgrid API Key"
        className="mt-4 mb-2"
        value={sendgridApiKey}
        options={secretOptions}
        onChange={(e) => setSendgridApiKey(e.target.value)}
      />
    </>
  )
}

export default SendgridClient
