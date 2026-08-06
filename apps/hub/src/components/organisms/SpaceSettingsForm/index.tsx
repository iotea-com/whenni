'use client'

import useAuth from '@gruent/hub/hooks/useAuth'
import gruentClient from '@gruent/hub/lib/gruent'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import FormFieldText from '@gruent/libs/frontend/components/atoms/FormFieldText'
import { addToast } from '@gruent/libs/frontend/hooks/useToast'
import { Space } from '@prisma/client'
import { useMutation } from '@tanstack/react-query'
import { FC, FormEventHandler, useCallback, useState } from 'react'

type Props = {
  orgId: string
  space: Space
}

const SpaceSettingsForm: FC<Props> = ({ orgId, space: initialSpace }) => {
  const [space, setSpace] = useState(initialSpace)

  const { accessToken } = useAuth()

  const updateSpaceMutation = useMutation({
    mutationKey: ['updateSpace', space.id],
    mutationFn: async (space: Space) => {
      if (!accessToken) return

      const { errors } = await gruentClient(accessToken).spaces.update(orgId, space.id, space)

      if (errors && errors.length > 0) {
        addToast({
          title: 'Could not update the space',
          body: errors[0],
          level: 'error',
        })

        return errors[0]
      }

      addToast({
        title: 'Successfully updated the space',
        body: 'Your new settings have been saved.',
        level: 'success',
      })

      return
    },
  })

  const handleSubmit: FormEventHandler<HTMLFormElement> = useCallback(
    (e) => {
      e.preventDefault()
      updateSpaceMutation.mutate(space)
    },
    [updateSpaceMutation, space],
  )

  return (
    <form onSubmit={handleSubmit}>
      <FormFieldText
        label="Space Name"
        name="spaceName"
        value={space.name}
        onChange={(e) => {
          const name = e.target.value
          setSpace((current) => ({ ...current, name }))
        }}
      />
      <Button type="submit" text="Save" disabled={space.name === initialSpace.name} />
    </form>
  )
}

export default SpaceSettingsForm
