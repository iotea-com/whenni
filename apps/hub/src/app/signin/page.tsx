import SignInForm from '@iotea/hub/components/organisms/SignInForm'

export const metadata = {
  title: 'Sign in | IOTEA',
  description:
    'Integrating devices and services just got a lot easier. Drag-and-drop connections to automate tasks and scale within minutes. Try IOTEA during our beta launch.',
  openGraph: {
    type: 'website',
    url: `https://app.iotea.com/signin`,
    description:
      'Integrating devices and services just got a lot easier. Drag-and-drop connections to automate tasks and scale within minutes. Try IOTEA during our beta launch.',
    images: ['https://iotea.com/img/logos/app-icon-primary.png'],
  },
}

const SignInPage = async () => {
  return <SignInForm allowPassword />
}

export default SignInPage
