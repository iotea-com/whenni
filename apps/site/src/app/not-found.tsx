'use client'

import '../app/styles.css'
import Container from '@iotea/libs/frontend/components/templates/Container'
import Image from 'next/image'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import { useRouter } from 'next/navigation'

const NotFound = () => {
  const router = useRouter()

  return (
    <main className="flex flex-col">
      <Container className="text-center items-center mt-24 py-20 justify-center">
        <Image
          src="/img/logos/app-icon-primary.png"
          alt="IOTEA logo"
          width={200}
          height={200}
          className="mx-auto"
        />
        <h1 className="text-4xl my-10">404</h1>
        <div className="grid grid-cols-2 gap-4 mt-4">
          <Button
            text="Return to home"
            variant="secondary"
            className="grow"
            onClick={() => router.push('/')}
          />
          <Button text="Go back" className="grow" onClick={router.back} />
        </div>
      </Container>
    </main>
  )
}

export default NotFound
