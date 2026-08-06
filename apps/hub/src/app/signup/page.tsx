import { verifyInvitationToken } from '@gruent/hub/actions/members'
import SignUpForm from '@gruent/hub/components/organisms/SignUpForm'

export const metadata = {
  title: 'Sign up | GRUENT',
  description:
    'The internet of things just got easier. Drag-and-drop connections to automate tasks and scale within minutes. Create an account to get started for free.',
  openGraph: {
    type: 'website',
    url: `https://app.gruent.com/signup`,
    description:
      'The internet of things just got easier. Drag-and-drop connections to automate tasks and scale within minutes. Create an account to get started for free.',
    images: ['https://gruent.com/img/logos/app-icon-primary.png'],
  },
}

const SignUpPage = async ({ searchParams }) => {
  const { inviteToken } = await searchParams

  const { data: invitation, error: inviteTokenError } = await (async () => {
    if (!inviteToken) return { data: undefined, error: undefined }

    return await verifyInvitationToken(process.env.JWT_SECRET!, inviteToken)
  })()

  return (
    <SignUpForm
      invitation={invitation}
      inviteToken={inviteToken ?? undefined}
      inviteTokenError={inviteTokenError}
    />
  )
}

export default SignUpPage
