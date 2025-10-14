'use client'

import useAuth from '@iotea/hub/hooks/useAuth'
import ioteaClient from '@iotea/hub/lib/iotea'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import FormFieldSelect from '@iotea/libs/frontend/components/atoms/FormFieldSelect'
import { openModal } from '@iotea/libs/frontend/hooks/useModal'
import { Thing } from '@prisma/client'
import { useQuery } from '@tanstack/react-query'
import { Dispatch, FC, useEffect, useMemo, useState } from 'react'

type Props = {
  setFormData: Dispatch<Map<string, string | number | string[] | number[]>>
  spaceId: string
  orgId: string
  initial?: Partial<Thing>
}

const NatsClient: FC<Props> = ({ setFormData, spaceId, initial }) => {
  const attributes = useMemo(() => {
    if (initial && initial.attributes && typeof initial.attributes === 'string')
      return JSON.parse(initial.attributes)

    if (initial && initial.attributes && typeof initial.attributes === 'object')
      return initial.attributes

    return undefined
  }, [initial])

  const [server, setServer] = useState<string>(attributes?.server)

  const { accessToken } = useAuth()

  const {
    status: getNatsServersStatus,
    data: natsServers,
    error: getNatsServersError,
  } = useQuery({
    queryKey: ['things', 'natsServers', spaceId],
    queryFn: async () => {
      if (!accessToken) throw new Error('Invalid auth session.')

      const { data: things, errors } = await ioteaClient(accessToken).things.list(spaceId, {
        category: 'NATS_SERVER',
      })

      if (errors && errors.length > 0) throw new Error(errors[0])

      return things
    },
  })

  const natsServerOptions = useMemo(() => {
    if (getNatsServersError) return []
    if (getNatsServersStatus !== 'success') return []
    if (!natsServers) return []

    return natsServers.map((thing) => {
      return {
        label: thing.name,
        value: thing.id,
      }
    })
  }, [natsServers, getNatsServersError, getNatsServersStatus])

  useEffect(() => {
    const updatedFormData = new Map()
    updatedFormData.set('server', server)
    setFormData(updatedFormData)
  }, [server, setFormData])

  if (!natsServers || getNatsServersStatus !== 'success') return

  if (natsServers.length <= 0) {
    return (
      <p>
        No NATS servers were found in this space. You'll have to&nbsp;
        <Button variant="underline" onClick={() => openModal('createThing')}>
          create a NATS server
        </Button>
        &nbsp;in this space before creating a NATS client.
      </p>
    )
  }

  return (
    <>
      <FormFieldSelect
        name="server"
        label="Server"
        className="mt-4 mb-2"
        options={natsServerOptions}
        value={server}
        onChange={(e) => setServer(e.target.value)}
      />
    </>
  )
}

export default NatsClient
