import SignInForm from '@gruent/hub/components/organisms/SignInForm'

export const metadata = {
  title: 'Sign in | GRUENT',
  description:
    'Integrating devices and services just got a lot easier. Drag-and-drop connections to automate tasks and scale within minutes. Try GRUENT during our beta launch.',
  openGraph: {
    type: 'website',
    url: `https://app.gruent.com/signin`,
    description:
      'Integrating devices and services just got a lot easier. Drag-and-drop connections to automate tasks and scale within minutes. Try GRUENT during our beta launch.',
    images: ['https://gruent.com/img/logos/app-icon-primary.png'],
  },
}

const SignInPage = async () => {
  return <SignInForm allowPassword />
}

export default SignInPage
