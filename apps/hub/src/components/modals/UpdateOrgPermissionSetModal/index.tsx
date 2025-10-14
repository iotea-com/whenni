'use client'

import { FC, useCallback, useRef } from 'react'
import Modal from '@iotea/libs/frontend/components/organisms/Modal'
import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import { useMutation } from '@tanstack/react-query'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { closeModal } from '@iotea/libs/frontend/hooks/useModal'
import ioteaClient from '@iotea/hub/lib/iotea'
import { useRouter } from 'next/navigation'
import { PermissionSet } from '@prisma/client'
import { defaultOrgPermissions } from '@iotea/libs/http/permissions'
import FormFieldSelect from '@iotea/libs/frontend/components/atoms/FormFieldSelect'
import useAuth from '@iotea/hub/hooks/useAuth'

type Props = {
  orgId: string
  spaceId?: string
  showModal: boolean
  initialPermissionSet?: PermissionSet
  onClose?: () => void
}

type UpdatePermissionSetInput = Pick<PermissionSet, 'id' | 'name' | 'permissions'>

const UpdateOrgPermissionSetModal: FC<Props> = ({
  orgId,
  spaceId,
  showModal,
  initialPermissionSet,
  onClose,
}) => {
  const { accessToken } = useAuth()
  const router = useRouter()

  const formRef = useRef<HTMLFormElement>(null)

  const updatePermissionSetMutation = useMutation({
    mutationKey: ['updatePermissionSetMutation', initialPermissionSet?.id],
    mutationFn: async (permissionSetInput: UpdatePermissionSetInput) => {
      if (!accessToken) return

      const { id, name, permissions } = permissionSetInput

      const { errors } = await ioteaClient(accessToken).permissions.update(
        orgId,
        id,
        name,
        permissions,
        spaceId,
      )

      if (errors && errors.length > 0) {
        addToast({
          title: 'Could not update permission set',
          body: errors[0],
          level: 'error',
        })
        return false
      }

      addToast({
        title: 'Successfully updated the permission set',
        body: `"${name}" has been successfully been updated. These changes should be reflected immediately.`,
        level: 'success',
      })

      formRef.current?.reset()
      closeModal()
      router.refresh()
      return true
    },
  })

  const handleClose = useCallback(() => {
    formRef.current?.reset()
    closeModal()
    if (onClose) onClose()
  }, [formRef, onClose])

  const handleSubmit = useCallback(() => {
    if (!formRef.current) return
    if (!initialPermissionSet) return

    const formData = new FormData(formRef.current)
    const name = formData.get('name') as string
    const permissions = structuredClone(defaultOrgPermissions)

    formData.forEach((action, namespace) => {
      // Check if the key is a namespace that we can use
      if (permissions[namespace] === undefined || permissions[namespace][action] === undefined)
        return

      // Update the permission map if the checkbox is checked
      permissions[namespace][action] = true
    })

    let permissionList: string[] = []
    for (const namespace of Object.keys(permissions)) {
      for (const action of Object.keys(permissions[namespace])) {
        if (permissions[namespace][action])
          permissionList = [...permissionList, `${namespace}:${action}`]
      }
    }

    if (!name) {
      addToast({
        title: 'Permission set name is required',
        body: 'Please add a name for your permission set and try again.',
        level: 'warning',
      })
      return
    }

    const updatePermissionSetInput: UpdatePermissionSetInput = {
      id: initialPermissionSet.id,
      name,
      permissions: permissionList,
    }

    updatePermissionSetMutation.mutate(updatePermissionSetInput)
  }, [formRef, updatePermissionSetMutation, initialPermissionSet])

  const renderInitialPermissions = useCallback(() => {
    if (!initialPermissionSet) return
    const initialPermissionsList = [...initialPermissionSet.permissions]
    const initialPermissions = structuredClone(defaultOrgPermissions)

    for (const permission of initialPermissionsList) {
      const [namespace, action] = permission.split(':')
      if (initialPermissions[namespace]) {
        initialPermissions[namespace][action] = true
      }
    }

    return Object.keys(initialPermissions).map((namespace) => {
      return (
        <div key={namespace}>
          <h3 className="mb-1 font-semibold">{namespace}</h3>
          <div className="px-6 py-4 mb-4 border border-gray-200 dark:border-gray-700 rounded-sm">
            <FormFieldSelect
              name={namespace}
              label="Select permissions"
              variant="cards"
              multiple
              value={initialPermissionsList
                .filter((permission) => permission.startsWith(namespace))
                .map((permission) => permission.split(':')[1])}
              options={Object.keys(initialPermissions[namespace]).map((action) => ({
                value: action,
                label: action,
              }))}
            />
          </div>
        </div>
      )
    })
  }, [initialPermissionSet])

  if (!showModal) return

  return (
    <Modal
      id="updateOrgPermissionSet"
      onAccept={handleSubmit}
      onClose={handleClose}
      acceptText="Update"
    >
      <h2>Update organization permission set</h2>
      <p className="text-gray-500">
        <small>Edit the details for an organization permission set.</small>
      </p>
      <div className="relative">
        <form ref={formRef} onSubmit={(e) => e.preventDefault()} className="mb-4">
          <FormFieldText
            name="name"
            label="Name"
            placeholder="Name of the permission set"
            defaultValue={initialPermissionSet?.name ?? ''}
            inputType="text"
            className="mt-4 mb-2"
            autocomplete="off"
          />
          {renderInitialPermissions()}
        </form>
      </div>
    </Modal>
  )
}

export default UpdateOrgPermissionSetModal
