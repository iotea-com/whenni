'use client'

import { FC, useCallback, useRef } from 'react'
import Modal from '@iotea/libs/frontend/components/organisms/Modal'
import { useMutation } from '@tanstack/react-query'
import { closeModal } from '@iotea/libs/frontend/hooks/useModal'
import { useRouter } from 'next/navigation'
import FormFieldSelect from '@iotea/libs/frontend/components/atoms/FormFieldSelect'
import { PermissionSet } from '@prisma/client'

type Props = {
  orgId: string
  userId: string
  showModal: boolean
  initialPermissionSetId: string
  permissionSets: PermissionSet[]
  onClose?: () => void
}

const EditMemberPermissionsModal: FC<Props> = ({
  orgId,
  userId,
  showModal,
  initialPermissionSetId,
  permissionSets,
  onClose,
}) => {
  const router = useRouter()

  const formRef = useRef<HTMLFormElement>(null)

  const updateMemberPermissionSetMutation = useMutation({
    mutationKey: ['updateMemberPermissionSetMutation', orgId, userId, initialPermissionSetId],
    mutationFn: async (_permissionSetId: string) => {
      // const { data: _results, error } = await ioteaClient(
      //   session.jwt,
      // ).spaces.profiles.

      // if (error) {
      //   addToast({
      //     title: 'Could not update permission set',
      //     body: error,
      //     level: 'error',
      //   })
      //   return false
      // }

      // addToast({
      //   title: 'Successfully updated the member's permissions',
      //   body: `The member's permissions have been successfully been updated. These changes should be reflected immediately.`,
      //   level: 'success',
      // })

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
    const formData = new FormData(formRef.current)
    const permissionSetId = formData.get('permissionSetId') as string

    updateMemberPermissionSetMutation.mutate(permissionSetId)
  }, [formRef, updateMemberPermissionSetMutation])

  if (!showModal) return null

  return (
    <Modal
      id="editUserPermissionSet"
      onAccept={handleSubmit}
      onClose={handleClose}
      acceptText="Update"
    >
      <h2>Edit member&apos;s permission set</h2>
      <p className="text-gray-500">
        <small>Change a team member&apos;s permission set.</small>
      </p>
      <div className="relative">
        <form ref={formRef} onSubmit={(e) => e.preventDefault()} className="mb-4">
          <FormFieldSelect
            name="permissionSetId"
            label="Permission Set"
            defaultValue={initialPermissionSetId}
            options={permissionSets.map((ps) => {
              return { value: ps.id, label: ps.name }
            })}
          />
        </form>
      </div>
    </Modal>
  )
}

export default EditMemberPermissionsModal
