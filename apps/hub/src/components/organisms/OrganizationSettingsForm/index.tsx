'use client'

import useAuth from '@iotea/hub/hooks/useAuth'
import ioteaClient from '@iotea/hub/lib/iotea'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { Organization } from '@prisma/client'
import { useMutation } from '@tanstack/react-query'
import { FC, FormEventHandler, useCallback, useState } from 'react'

type Props = {
  organization: Organization
}

const OrganizationSettingsForm: FC<Props> = ({ organization: initialOrganization }) => {
  const [organization, setOrganization] = useState<Organization>(initialOrganization)

  const { accessToken } = useAuth()

  const updateOrganizationMutation = useMutation({
    mutationKey: ['updateOrganization', organization.id],
    mutationFn: async (organization: Organization) => {
      if (!accessToken) return

      const { errors } = await ioteaClient(accessToken).organizations.update(organization)

      if (errors && errors.length > 0) {
        addToast({
          title: 'Could not update the organization',
          body: errors[0],
          level: 'error',
        })

        return errors[0]
      }

      addToast({
        title: 'Successfully updated the organization',
        body: 'Your new settings have been saved.',
        level: 'success',
      })

      return
    },
  })

  const handleSubmit: FormEventHandler<HTMLFormElement> = useCallback(
    (e) => {
      e.preventDefault()
      updateOrganizationMutation.mutate(organization)
    },
    [updateOrganizationMutation, organization],
  )

  return (
    <form onSubmit={handleSubmit}>
      <FormFieldText
        label="Organization Name"
        name="organizationName"
        value={organization.name}
        onChange={(e) => {
          const name = e.target.value
          setOrganization((current) => ({ ...current, name }))
        }}
      />
      <Button type="submit" text="Save" disabled={organization.name === initialOrganization.name} />
    </form>
  )
}

export default OrganizationSettingsForm
