import Container from '@iotea/libs/frontend/components/templates/Container'

import { Metadata } from 'next'

import styles from './page.module.scss'
import Button from '@iotea/libs/frontend/components/atoms/Button'

// // do not cache this page
// export const revalidate = 1 * 60 // seconds

export const metadata: Metadata = {
  title: 'Documentation, tutorials, and API reference | IOTEA',
  description:
    'Learn how to use IOTEA to seamlessly integrate devices and services. Get tutorials, API reference, and more.',
  openGraph: {
    type: 'website',
    url: `https://iotea.com/docs`,
    title: 'Documentation, tutorials, and API reference | IOTEA',
    description:
      'Learn how to use IOTEA to seamlessly integrate devices and services. Get tutorials, API reference, and more.',
    images: ['https://iotea.com/img/logos/app-icon-primary.png'],
  },
}

const DocsPage = () => {
  return (
    <>
      <header id={styles.header} className={`${styles.dotGrid} border-b border-secondary`}>
        <Container className="py-32">
          <div className={`${styles.headerContent}`}>
            <h1 className="text-2xl lg:text-4xl max-w-2xl text-secondary">Documentation</h1>
            <p className="text-xl text-gray-600">Explore our guides and examples.</p>
          </div>
        </Container>
      </header>
      <div className="overflow-hidden" style={{ marginTop: '-25rem' }}>
        <div
          style={{
            marginLeft: '55rem',
            borderRight: '2000px solid rgb(23 55 24)',
            borderTop: '400px solid transparent',
          }}
        />
      </div>
      <main>
        <section className="mt-10 mb-20">
          <Container>
            <h1 className="max-w-xl">Our documentation is coming soon!</h1>
            <p className="max-w-xl">
              Once we launch our upcoming beta platform, we will have a comprehensive documentation
              section for features, tutorials, and API reference. For now, check out our{' '}
              <a href="https://discord.gg/prAJjW426d" target="_blank">
                <Button variant="underline" className="!inline px-0 py-0">
                  community Discord
                </Button>
              </a>{' '}
              or{' '}
              <Button variant="underline" modalId="betaSignup" className="!inline px-0 py-0">
                sign up for our upcoming launch.
              </Button>
            </p>
          </Container>
        </section>
      </main>
    </>
  )
}

export default DocsPage
