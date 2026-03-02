'use client'

import { FC, useCallback, useRef } from 'react'
import Modal from '@iotea/libs/frontend/components/organisms/Modal'
import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { closeModal } from '@iotea/libs/frontend/hooks/useModal'
import { handleAddPermissionSet } from '@iotea/hub/actions/permissionSets'
import { defaultOrgPermissions, defaultSpacePermissions } from '@iotea/libs/http/permissions'
import FormFieldSelect from '@iotea/libs/frontend/components/atoms/FormFieldSelect'

type Props = {
  orgId: string
  spaceId?: string
}

const CreatePermissionSetModal: FC<Props> = ({ orgId, spaceId }) => {
  const formRef = useRef<HTMLFormElement>(null)

  const handleFormClose = useCallback(() => {
    formRef.current?.reset()
    closeModal()
  }, [formRef])

  const handleFormSubmit = () => {
    formRef.current?.requestSubmit()
  }

  const defaultPermissions = (() => {
    if (spaceId) return structuredClone(defaultSpacePermissions)
    else return structuredClone(defaultOrgPermissions)
  })()

  return (
    <Modal
      id="addPermissionSet"
      acceptText="Create"
      onAccept={handleFormSubmit}
      onClose={handleFormClose}
    >
      <h2>Add a permission set</h2>
      <p className="text-gray-500">
        <small>Create a set of permissions to assign to your team members.</small>
      </p>
      <form
        ref={formRef}
        className="relative mb-4"
        onSubmit={async (e) => {
          e.preventDefault()

          const formData = new FormData(e.target as HTMLFormElement)
          const name = formData.get('name') as string
          const permissions = structuredClone(defaultPermissions)

          formData.forEach((action, namespace) => {
            // Check if the key is a namespace that we can use
            if (
              permissions[namespace] === undefined ||
              permissions[namespace][action] === undefined
            )
              return

            // Update the permission map if the checkbox is checked
            permissions[namespace][action] = true
          })

          let permissionList: string[] = []
          for (const category of Object.keys(permissions)) {
            for (const selector of Object.keys(permissions[category])) {
              if (permissions[category][selector])
                permissionList = [...permissionList, `${category}:${selector}`]
            }
          }

          const { error } = await handleAddPermissionSet(orgId, name, permissionList, spaceId)
          if (error) {
            addToast({
              title: 'Could not create the permission set',
              body: `The permission set that you provided could not be created: ${error}`,
              level: 'error',
            })
          } else {
            formRef.current?.reset()
            closeModal()

            addToast({
              title: 'Successfully created permission set',
              body: `Successfully added the new permission set to your ${spaceId ? 'space' : 'organization'}.`,
              level: 'success',
            })
          }
        }}
      >
        <FormFieldText
          name="name"
          label="Name"
          placeholder="Name of the permission set"
          className="my-4"
          autocomplete="off"
        />
        {Object.keys(defaultPermissions).map((namespace) => {
          return (
            <div key={namespace}>
              <h3 className="mb-1 font-semibold">{namespace}</h3>
              <div
                key={namespace}
                className="px-6 py-4 mb-4 border border-gray-200 dark:border-gray-700 rounded-xs"
              >
                <FormFieldSelect
                  name={namespace}
                  label="Select permissions"
                  variant="cards"
                  multiple
                  options={Object.keys(defaultPermissions[namespace]).map((action) => ({
                    value: action,
                    label: action,
                  }))}
                />
              </div>
            </div>
          )
        })}
      </form>
    </Modal>
  )
}

export default CreatePermissionSetModal
