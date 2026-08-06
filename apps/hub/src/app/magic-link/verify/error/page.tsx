import AuthLayout from '@gruent/hub/components/layouts/AuthLayout'
import Button from '@gruent/libs/frontend/components/atoms/Button'

export const metadata = {
  title: 'Verify magic link | GRUENT',
}

const VerifyMagicLinkPage = async ({ searchParams }) => {
  const { error } = await searchParams

  const decodedError = (() => {
    try {
      return decodeURIComponent(error)
    } catch (_e) {
      return 'An unknown error occurred when verifying the magic link.'
    }
  })()

  return (
    <AuthLayout>
      <div className="flex flex-col items-center justify-center p-20 bg-white rounded-lg shadow-md">
        <p className="mb-5">{decodedError}</p>
        <Button href="/signin" className="w-full">
          Back to sign in
        </Button>
      </div>
    </AuthLayout>
  )
}

export default VerifyMagicLinkPage
