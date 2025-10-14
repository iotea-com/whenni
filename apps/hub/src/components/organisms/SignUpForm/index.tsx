'use client'

import Image from 'next/image'
import { FC, FormEventHandler, useEffect, useState } from 'react'
import styles from './index.module.scss'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import Link from 'next/link'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { useRouter } from 'next/navigation'
import { RemixIcon, riEyeLine } from '@mwarnerdotme/react-remixicon'
// import { signIn } from 'next-auth/react'
import { signup } from '@iotea/hub/actions/auth'

type Props = {
  invitation?: {
    email: string
    orgId: string
  }
  inviteToken?: string
  inviteTokenError?: string
}

const SignUpForm: FC<Props> = ({ invitation, inviteToken, inviteTokenError }) => {
  const [intialEmail, setIntialEmail] = useState<string>()
  const [showPassword, setShowPassword] = useState<boolean>(false)

  const router = useRouter()

  useEffect(() => {
    if (invitation) {
      setIntialEmail(invitation.email)
    }
  }, [invitation])

  useEffect(() => {
    if (inviteTokenError) {
      addToast({
        title: 'Invalid or expired invitation token',
        body: inviteTokenError,
        level: 'error',
        ttl: 5,
      })
    }
  }, [inviteTokenError])

  // const handleGithubOauthClick = async () => {
  //   await signIn('github', {
  //     redirectTo: `${window.location.origin}/dashboard`,
  //   })
  // }

  const handleSignupFormSubmit: FormEventHandler = async (e) => {
    e.preventDefault()

    const formData = new FormData(e.target as HTMLFormElement)

    const email = formData.get('email') as string
    const password = formData.get('password') as string

    // Sign up
    const { error } = await signup({
      email,
      password,
      inviteToken,
      origin: window.location.origin,
    })

    if (error) {
      addToast({
        title: 'Could not sign up',
        body: error,
        level: 'error',
        ttl: 5,
      })
      return
    }

    addToast({
      title: 'Account created',
      body: 'An email has been sent to verify your account. Once verified, sign in to continue.',
      level: 'success',
      ttl: 5,
    })
    router.push(`/verify-email?email=${encodeURIComponent(email)}`)
  }

  return (
    <div className={styles.wrapper}>
      {invitation && (
        <div className="rounded-lg bg-white px-10 py-8 w-full mb-4">
          <h2 className="text-gray-800 text-2xl">👋 You've been invited!</h2>
          <p>
            Once you sign up, you'll be able to access the organization that you were invited to.
          </p>
        </div>
      )}
      <div className={styles.signupForm}>
        <Image
          src={'/img/logos/app-icon-primary.png'}
          height={35}
          width={35}
          alt="IOTEA logo"
          className="mb-5"
        />
        <h2 className="text-gray-800 text-2xl">Sign Up</h2>
        <h3 className="text-gray-500 text-lg font-normal">to start using IOTEA</h3>
        {/* {process.env.NEXT_PUBLIC_GITHUB_ENABLED === 'true' && (
          <>
            <div className="my-6">
              <Button
                className="flex items-center w-full py-2"
                variant="transparent"
                onClick={handleGithubOauthClick}
              >
                <Image src="/img/icons/github.svg" alt="Github logo" height={18} width={18} />
                <span className="text-sm ml-2">Sign up with Github</span>
              </Button>
            </div>
            <div className="flex items-center my-6">
              <hr className="grow border-gray-400" />
              <p className="text-sm px-4">or</p>
              <hr className="grow border-gray-400" />
            </div>
          </>
        )} */}
        <form onSubmit={handleSignupFormSubmit}>
          <FormFieldText
            name="email"
            label="Email"
            inputType="email"
            backgroundColor="bg-gray-50"
            defaultValue={intialEmail}
          />
          <div className="relative flex items-center">
            <FormFieldText
              name="password"
              label="Password"
              inputType={showPassword ? 'text' : 'password'}
              backgroundColor="bg-gray-50"
              className="w-full"
            />
            <RemixIcon
              className="absolute right-3 transition text-gray-500 hover:text-gray-600 cursor-pointer"
              icon={riEyeLine}
              onClick={() => setShowPassword((current) => !current)}
            />
          </div>
          <Button id="submit" className="w-full" type="submit" text="Sign up" />
        </form>
        <p className="text-sm mt-6">
          Already have an account?{' '}
          <Link id="signin" href={'/signin'} className="text-green-700">
            Sign in
          </Link>
        </p>
      </div>
    </div>
  )
}

export default SignUpForm
