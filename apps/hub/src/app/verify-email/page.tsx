import SignInForm from '@gruent/hub/components/organisms/SignInForm'

const VerifyEmailPage = async ({ searchParams }) => {
  const { email } = await searchParams

  return (
    <div className="flex flex-col max-w-md">
      <div className="rounded-lg bg-white px-10 py-8 w-full mb-4">
        <h2 className="text-gray-800 text-2xl">📬 Email verification</h2>
        <p>
          To sign in for the first time, we'll send a magic link to your email address to verify
          your email. You can sign in with a password after you've verified your email.
        </p>
      </div>
      <SignInForm allowPassword={false} initialEmail={email} />
    </div>
  )
}

export default VerifyEmailPage
