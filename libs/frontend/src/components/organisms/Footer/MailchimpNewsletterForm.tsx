'use client'

import { FC, useState } from 'react'
import { useFormStatus } from 'react-dom'
import FormFieldText from '../../atoms/FormFieldText'
import Button from '../../atoms/Button'
import addNewsletterMember from '@gruent/site/actions/addNewsletterMember'

const FormSubmitButton: FC = () => {
  const { pending } = useFormStatus()

  return <Button type="submit" text={pending ? 'Submitting...' : 'Sign up'} disabled={pending} />
}

const MailchimpNewsletterForm: FC = () => {
  const [message, setMessage] = useState<string>()

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const formData = new FormData(event.target as HTMLFormElement)
    const email = formData.get('email')

    if (!email || email === '') {
      setMessage('Please enter an email address')
      return
    }

    addNewsletterMember(email as string)
  }

  return (
    <>
      <form onSubmit={handleSubmit} className="mt-6">
        <div className="flex">
          <FormFieldText
            backgroundColor="bg-gray-100"
            className="my-0! grow"
            name="email"
            label="Email Address"
            inputType="email"
            placeholder="jon.snow@winterfell.com"
          />
          <FormSubmitButton />
        </div>
      </form>
      {message && (
        <p>
          <small className="text-gray-500">{message}</small>
        </p>
      )}
    </>
  )
}

export default MailchimpNewsletterForm
