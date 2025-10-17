'use client'

import { FC, useRef } from 'react'

import Modal from '@iotea/libs/frontend/components/organisms/Modal'
import FormFieldText from '../../../atoms/FormFieldText'
import { closeModal } from '@iotea/libs/frontend/hooks/useModal'
import addNewsletterMember from '@iotea/site/actions/addNewsletterMember'
import { RemixIcon, riLoader2Fill } from '@mwarnerdotme/react-remixicon'
import { useMutation } from '@tanstack/react-query'

type Props = {}

const NewsletterModal: FC<Props> = ({}) => {
  const formRef = useRef<HTMLFormElement>(null)

  const { mutate, data, isPending } = useMutation({
    mutationFn: async (email: string) => {
      if (!email || email === '') {
        return { message: 'Please enter an email address' }
      }

      const result = await addNewsletterMember(email)
      return result
    },
  })

  const onSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const formData = new FormData(event.target as HTMLFormElement)
    const email = formData.get('email')

    await mutate(email as string)
  }

  const onAccept = () => {
    if (!formRef.current) return
    formRef.current.requestSubmit()
  }

  const onClose = () => {
    if (!formRef.current) return
    formRef.current.reset()
    closeModal()
  }

  return (
    <Modal id="newsletter" acceptText="Sign Up" onAccept={onAccept} onClose={onClose}>
      <h2>IOTEA Newsletter Signup</h2>
      <p>
        Receive emails when we have product updates, new blog posts and videos, and other
        announcements.
      </p>
      <form ref={formRef} onSubmit={onSubmit} className="mt-6 mb-4">
        <div className="flex">
          <FormFieldText
            backgroundColor="bg-gray-100"
            className="my-0! grow"
            name="email"
            label="Email Address"
            inputType="email"
            placeholder="jon.snow@winterfell.com"
          />
        </div>
        {isPending && (
          <p>
            <RemixIcon icon={riLoader2Fill} size="sm" />
          </p>
        )}
        {data?.message && (
          <p>
            <small className="text-gray-500">{data.message}</small>
          </p>
        )}
      </form>
    </Modal>
  )
}

export default NewsletterModal
