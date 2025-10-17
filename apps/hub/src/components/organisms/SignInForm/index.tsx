'use client'

import Image from 'next/image'
import { FormEventHandler, useState } from 'react'
import styles from './index.module.css'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import FormFieldText from '@iotea/libs/frontend/components/atoms/FormFieldText'
import Link from 'next/link'
import { addToast } from '@iotea/libs/frontend/hooks/useToast'
import { useRouter } from 'next/navigation'
import { RemixIcon, riEyeLine, riEyeOffLine } from '@mwarnerdotme/react-remixicon'
// import { signIn } from 'next-auth/react'
import { signInWithCredentials, signInWithMagicLink } from '@iotea/hub/actions/auth'
import Modal from '@iotea/libs/frontend/components/organisms/Modal'
import { openModal } from '@iotea/libs/frontend/hooks/useModal'

type Props = {
  allowPassword?: boolean
  initialEmail?: string
}

const SignInForm = ({ allowPassword = true, initialEmail }: Props) => {
  const [passwordEnabled, setPasswordEnabled] = useState<boolean>(false)
  const [showPassword, setShowPassword] = useState<boolean>(false)

  const router = useRouter()

  // const handleGithubOauthClick = async () => {
  //   await signIn('github', {
  //     redirectTo: `${window.location.origin}/dashboard`,
  //   })
  // }

  const handleLoginFormSubmit: FormEventHandler = async (e) => {
    e.preventDefault()

    const formData = new FormData(e.target as HTMLFormElement)

    const email = formData.get('email') as string
    const password = formData.get('password') as string

    if (password) {
      const { error, redirect } = await signInWithCredentials({
        email,
        password,
        redirectTo: `${window.location.origin}/dashboard`,
      })

      if (redirect) {
        router.push(redirect)
        return
      }

      if (error === 'email is not verified') {
        router.push(`/verify-email?email=${encodeURIComponent(email)}`)
        return
      }

      if (error) {
        addToast({
          title: 'Could not sign in',
          body: error,
          level: 'error',
          ttl: 5,
        })
      } else {
        if (!error) router.push('/dashboard')
      }

      return
    }

    const { error, redirect } = await signInWithMagicLink({
      email,
      appUrl: `${window.location.origin}`,
      redirectTo: `${window.location.origin}/dashboard`,
    })

    if (redirect) {
      router.push(redirect)
      return
    }

    if (error) {
      addToast({
        title: 'Could not send magic link',
        body: error,
        level: 'error',
        ttl: 5,
      })
    } else {
      addToast({
        title: 'Magic link sent',
        body: 'Check your email to finish authenticating',
        level: 'success',
        ttl: 5,
      })
    }
  }

  return (
    <div className={styles.wrapper}>
      <div className={styles.signinForm}>
        <Modal id="forgot-password" showAccept={false}>
          <h2>Forgot password</h2>
          <p className="w-96">
            If you have forgotten your password, you can use magic link authentication or a
            third-party provider to sign in to your account. Once signed in, visit your account
            settings page to update your password.
          </p>
        </Modal>
        <Image
          src={'/img/logos/app-icon-primary.png'}
          height={35}
          width={35}
          alt="IOTEA logo"
          className="mb-5"
        />
        <h2 className="text-gray-800 text-2xl">Sign In</h2>
        <h3 className="text-gray-500 text-lg font-normal">to continue to IOTEA</h3>
        {/* {process.env.NEXT_PUBLIC_GITHUB_ENABLED === 'true' && (
          <>
            <div className="my-6">
              <Button
                className="flex items-center w-full py-2"
                variant="transparent"
                onClick={handleGithubOauthClick}
              >
                <Image src="/img/icons/github.svg" alt="Github logo" height={18} width={18} />
                <span className="text-sm ml-2">Sign in with Github</span>
              </Button>
            </div>
            <div className="flex items-center my-6">
              <hr className="grow border-gray-400" />
              <p className="text-sm px-4">or</p>
              <hr className="grow border-gray-400" />
            </div>
          </>
        )} */}
        <form onSubmit={handleLoginFormSubmit}>
          <FormFieldText
            name="email"
            label="Email"
            inputType="email"
            backgroundColor="bg-gray-50"
            defaultValue={initialEmail}
          />
          {passwordEnabled && (
            <div className="relative flex items-center">
              <FormFieldText
                name="password"
                label="Password"
                inputType={showPassword ? 'text' : 'password'}
                backgroundColor="bg-gray-50"
                className="w-full"
              />
              {showPassword ? (
                <RemixIcon
                  className="absolute right-3 transition text-gray-500 hover:text-gray-600 cursor-pointer"
                  icon={riEyeOffLine}
                  onClick={() => setShowPassword((current) => !current)}
                />
              ) : (
                <RemixIcon
                  className="absolute right-3 transition text-gray-500 hover:text-gray-600 cursor-pointer"
                  icon={riEyeLine}
                  onClick={() => setShowPassword((current) => !current)}
                />
              )}
            </div>
          )}
          <Button
            id="submit"
            className="w-full"
            type="submit"
            text={passwordEnabled ? 'Continue' : 'Receive magic link'}
          />
        </form>
        <p className="text-sm mt-6">
          No account?{' '}
          <Link id="signup" href={'/signup'} className="text-green-700">
            Sign up
          </Link>
        </p>
        {passwordEnabled ? (
          <div>
            <span
              id="magicLink"
              className="text-sm text-green-700 cursor-pointer"
              onClick={() => setPasswordEnabled(false)}
            >
              Login with magic link
            </span>
          </div>
        ) : (
          <>
            {allowPassword && (
              <div>
                <span
                  id="userpass"
                  className="text-sm text-green-700 cursor-pointer"
                  onClick={() => setPasswordEnabled(true)}
                >
                  Login with password
                </span>
              </div>
            )}
          </>
        )}
        {passwordEnabled && (
          <div>
            <span
              id="forgot-password"
              className="text-green-700 cursor-pointer text-sm font-semibold"
              onClick={() => {
                openModal('forgot-password')
              }}
            >
              Forgot password
            </span>
          </div>
        )}
      </div>
    </div>
  )
}

export default SignInForm
