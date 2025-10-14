'use client'

import { FC, useMemo, useState } from 'react'
import Modal from '@iotea/libs/frontend/components/organisms/Modal'
import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import { useQuery, useMutation } from '@tanstack/react-query'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import ioteaClient from '@iotea/hub/lib/iotea'
import { handleAddMember, handleInviteMember } from '@iotea/hub/actions/members'
import { closeModal } from '@iotea/libs/frontend/hooks/useModal'
import { useDebounce } from '@react-hooks-library/core'
import {
  RemixIcon,
  riLoader2Fill,
  riMailAddLine,
  riMailCheckLine,
} from '@mwarnerdotme/react-remixicon'
import useAuth from '@iotea/hub/hooks/useAuth'

type Props = {
  orgId: string
}

const CreateMemberModal: FC<Props> = ({ orgId }) => {
  const [emailInput, setEmailInput] = useState<string>('')

  const { accessToken } = useAuth()

  const debouncedEmailInput = useDebounce(emailInput, 500)

  // Search for users whenever the input changes
  const { data: userList = [], isLoading: isSearchMembersLoading } = useQuery({
    queryKey: ['searchMembers', debouncedEmailInput],
    queryFn: async () => {
      if (!accessToken) return []

      const { data: users, errors } = await ioteaClient(accessToken).users.search(
        debouncedEmailInput,
        orgId,
      )

      if (errors && errors.length > 0) {
        addToast({
          title: 'Could not search members',
          body: errors[0],
          level: 'warning',
        })
        return []
      }

      return users ?? []
    },
    enabled: debouncedEmailInput.length > 3 && !!accessToken,
  })

  const isInputValidEmail = useMemo(() => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(emailInput)
  }, [emailInput])

  const {
    mutate: sendMemberInvite,
    isPending: isMemberInvitePending,
    isSuccess: isMemberInviteSuccess,
  } = useMutation({
    mutationKey: ['inviteMember', emailInput],
    mutationFn: async (email: string) => {
      if (!accessToken) return

      const { error } = await handleInviteMember(orgId, email, location.origin)

      if (error) {
        addToast({
          title: 'Could not send invitation',
          body: error,
          level: 'error',
        })
        throw new Error(error)
      } else {
        addToast({
          title: 'Invitation sent',
          body: 'The invitation has been sent to the user.',
          level: 'success',
        })
      }

      return
    },
  })

  const handleSendInvitation = (email: string) => {
    sendMemberInvite(email)
  }

  return (
    <Modal id="addMember" showAccept={false}>
      <h2>Add a member</h2>
      <p className="text-gray-500">
        <small>Invite people to collaborate on channels, data, devices, and more.</small>
      </p>
      <div className="relative">
        <FormFieldText
          value={emailInput}
          onChange={(e) => setEmailInput(e.target.value)}
          name="search"
          label="Search"
          placeholder="Search by email address"
          inputType="email"
          className="my-4"
          autocomplete="off"
        />
        {emailInput.length > 3 && (
          <ul className="absolute top-11 bg-white dark:bg-gray-900 border border-gray-200 w-full rounded-b">
            {userList.length <= 0 && isSearchMembersLoading && (
              <li className="py-3 px-5 border-b border-gray-200 dark:border-gray-700 text-gray-600 hover:text-gray-800 hover:bg-gray-100 transition cursor-pointer">
                <RemixIcon className="animate-spin" icon={riLoader2Fill} />
              </li>
            )}
            {userList.length <= 0 && !isSearchMembersLoading && !isInputValidEmail && (
              <li className="py-3 px-5 border-b border-gray-200 dark:border-gray-700 text-gray-600 hover:text-gray-800 hover:bg-gray-100 transition cursor-pointer">
                <p className="mb-0">No users were found.</p>
                <small>Enter a full email address to invite a new user.</small>
              </li>
            )}
            {userList.length > 0 &&
              userList.map((user) => {
                const { email, id } = user
                return (
                  <li
                    key={id}
                    className="py-3 px-5 border-b border-gray-200 dark:border-gray-700 text-gray-600 hover:text-gray-800 hover:bg-gray-100 dark:hover:bg-gray-800 transition cursor-pointer"
                    onClick={async () => {
                      const { error } = await handleAddMember(orgId, id)
                      if (error) {
                        addToast({
                          title: 'Could not add the member',
                          body: `The member that you selected could not be added to the space: ${error}`,
                          level: 'error',
                        })
                      } else {
                        setEmailInput('')
                        closeModal()

                        addToast({
                          title: 'Successfully added member',
                          body: 'Successfully added the member to your space.',
                          level: 'success',
                        })
                      }
                    }}
                  >
                    <p>{email}</p>
                  </li>
                )
              })}
            {userList.length <= 0 &&
              isInputValidEmail &&
              !isMemberInvitePending &&
              !isMemberInviteSuccess && (
                <li
                  className="py-3 px-5 flex items-center text-gray-600 hover:text-gray-800 hover:bg-gray-100 dark:hover:bg-gray-800 transition cursor-pointer"
                  onClick={() => handleSendInvitation(emailInput)}
                >
                  <RemixIcon icon={riMailAddLine} className="mr-1 dark:text-gray-300" />
                  <p>Invite {emailInput}</p>
                </li>
              )}
            {isMemberInvitePending && !isMemberInviteSuccess && (
              <li className="py-3 px-5 border-b border-gray-200 dark:border-gray-700 text-gray-600 hover:text-gray-800 hover:bg-gray-100 dark:hover:bg-gray-800 transition cursor-pointer">
                <RemixIcon className="animate-spin" icon={riLoader2Fill} />
              </li>
            )}
            {isMemberInviteSuccess && (
              <li className="py-3 px-5 flex items-center border-b border-gray-200 dark:border-gray-700 text-gray-600 hover:text-gray-800 hover:bg-gray-100 dark:hover:bg-gray-800 transition">
                <RemixIcon icon={riMailCheckLine} className="mr-1 dark:text-gray-300" />
                <p>Invite sent!</p>
              </li>
            )}
          </ul>
        )}
      </div>
    </Modal>
  )
}

export default CreateMemberModal
