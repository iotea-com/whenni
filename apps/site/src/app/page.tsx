import Container from '@iotea/libs/frontend/components/templates/Container'
import Button from '@iotea/libs/frontend/components/atoms/Button'

import styles from './page.module.scss'
import {
  RemixIcon,
  riFunctionLine,
  riCodeBoxLine,
  riDragDropLine,
  riPlugLine,
  riPriceTagLine,
  riCollapseDiagonal2Line,
  riServerLine,
  riDonutChartFill,
  riBarChart2Fill,
  riLineChartFill,
  riErrorWarningLine,
} from '@mwarnerdotme/react-remixicon'
import { Metadata } from 'next'
import RotatingVideoPlayer from '../components/RotatingVideoPlayer'
import Image from 'next/image'

// // do not cache this page
// export const revalidate = 1 * 60 // seconds

export const metadata: Metadata = {
  title: 'Make your next IoT project faster and easier | IOTEA',
  description:
    'Integrating devices and services just got a lot easier. Drag-and-drop connections to automate tasks and scale within minutes. Try IOTEA during our beta launch.',
  openGraph: {
    type: 'website',
    url: `https://iotea.com`,
    title: 'Make your next IoT project faster and easier | IOTEA',
    description:
      'Integrating devices and services just got a lot easier. Drag-and-drop connections to automate tasks and scale within minutes. Try IOTEA during our beta launch.',
    images: ['https://iotea.com/img/logos/app-icon-primary.png'],
  },
}

const HomePage = () => {
  return (
    <>
      <header id={styles.header} className={styles.dotGrid}>
        <Container className="pt-32">
          <h1 className="text-5xl sm:text-6xl lg:text-8xl max-w-2xl">
            <span className="mr-4">Connect</span>
            <div className={styles.scrollingWordsBox}>
              <ul>
                <li>anything</li>
                <li>devices</li>
                <li>databases</li>
                <li>APIs</li>
                <li>webhooks</li>
                <li>MQTT</li>
                <li>Postgres</li>
                <li>InfluxDB</li>
              </ul>
            </div>{' '}
            within minutes.
          </h1>
          <p className="text-lg mt-4 max-w-md">
            Ideas with hardware devices and software services just got easier. Drag-and-drop
            connections within minutes and scale to thousands of calls per second.
          </p>
          <div className="flex gap-4 mt-4">
            <Button text="Start Now" modalId="betaSignup" />
            {/* <Button text="Learn More" href="/products" variant="transparent" /> */}
          </div>
        </Container>
      </header>
      <main>
        <section className="grid grid-cols-8 w-full bg-white">
          <div className="grow h-3 bg-green-100" />
          <div className="grow h-3 bg-green-200" />
          <div className="grow h-3 bg-green-300" />
          <div className="grow h-3 bg-green-400" />
          <div className="grow h-3 bg-green-500" />
          <div className="grow h-3 bg-green-600" />
          <div className="grow h-3 bg-green-700" />
          <div className="grow h-3 bg-green-800" />
        </section>
        <section className="pt-24 pb-0 text-gray-100 bg-secondary">
          <Container>
            <div className="max-w-3xl mx-auto">
              <div>
                <h2 className="text-5xl lg:text-6xl my-8 text-primary text-center">
                  Integrate any hardware, software, or protocol
                </h2>
                <p className="text-2xl text-white text-center">
                  With a growing node library, we support your preferred services, databases,
                  message queues, and protocols. Bring your own devices and services, or use ours to
                  build ideas quicker than ever before.
                </p>
              </div>
            </div>
          </Container>
          <div className="mt-16 rounded mx-auto w-11/12 md:w-10/12" style={{ maxWidth: 1200 }}>
            <RotatingVideoPlayer
              videos={[
                { src: '/videos/channel_management.mp4', hash: '#workflows-video' },
                { src: '/videos/thing_management.mp4', hash: '#things-video' },
                { src: '/videos/model_management.mp4', hash: '#models-video' },
                { src: '/videos/secret_management.mp4', hash: '#secrets-video' },
              ]}
            />
          </div>
        </section>
        <section className="py-8">
          <Container>
            <div className="grid grid-cols-4 text-xs md:text-base gap-4 sm:gap-8 w-fit mx-auto border border-gray-200 rounded-full py-4 px-6 sm:px-12 text-center">
              <a href="#workflows-video" className="transition text-gray-700 hover:text-gray-800">
                Workflows
              </a>
              <a href="#things-video" className="transition text-gray-700 hover:text-gray-800">
                Things
              </a>
              <a href="#models-video" className="transition text-gray-700 hover:text-gray-800">
                Models
              </a>
              <a href="#secrets-video" className="transition text-gray-700 hover:text-gray-800">
                Secrets
              </a>
            </div>
          </Container>
        </section>
        <section className="bg-gray-50 pt-12 pb-36">
          <Container>
            <div className="grid lg:grid-cols-2">
              <div>
                <h2 className="text-5xl font-bold my-8">
                  Made by developers, designed for everyone
                </h2>
                <p className="text-lg">
                  Whether you&apos;ve been writing code for years or you have no desire to learn,
                  our platform is made to be accessible to anyone with an idea in mind. We believe
                  everyone should be able to build integrations faster than ever before.
                </p>
              </div>
            </div>
            <div className="mt-16 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
              <div>
                <RemixIcon icon={riFunctionLine} size="2x" className="text-green-500" />
                <h3 className="font-bold text-gray-700 mt-3 mb-1">Developer friendly</h3>
                <p className="text-gray-500">
                  Bringing custom hardware and software is easy. Libraries and SDKs will soon be
                  available for C/C++, Go, Node.js, and Python.
                </p>
              </div>
              <div>
                <RemixIcon icon={riPlugLine} size="2x" className="text-green-500" />
                <h3 className="font-bold text-gray-700 mt-3 mb-1">
                  Plug-and-play solutions (coming soon)
                </h3>
                <p className="text-gray-500">
                  Hardware takes a lot of time to get right. We provide out-of-the-box devices and
                  modules to make connecting faster.
                </p>
              </div>
              <div>
                <RemixIcon icon={riDragDropLine} size="2x" className="text-green-500" />
                <h3 className="font-bold text-gray-700 mt-3 mb-1">No-code</h3>
                <p className="text-gray-500">
                  Deploy production ready integrations without writing a single line of code. Build
                  and scale with the same performance as a custom solution.
                </p>
              </div>
              <div>
                <RemixIcon icon={riCodeBoxLine} size="2x" className="text-green-500" />
                <h3 className="font-bold text-gray-700 mt-3 mb-1">
                  Customize with code (coming soon)
                </h3>
                <p className="text-gray-500">
                  Write code inline or upload WASM modules for full control over how your
                  integrations work. We&apos;ll take care of deployment and scaling.
                </p>
              </div>
            </div>
          </Container>
        </section>
        <hr className="my-0 border border-green-200" />
        <section className="pt-32 pb-48 bg-gray-50">
          <Container>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-10">
              <div>
                <h2 className="text-2xl">Get started in minutes</h2>
                <p className="mt-4 mb-6 text-gray-500">
                  Create an account and start testing and prototyping for free until you&apos;re
                  ready to commit.
                </p>
                <div className="flex">
                  <Button text="Start Now" modalId="betaSignup" />
                  <a href="/pricing">
                    <Button className="ml-4" variant="transparent" text="See Pricing" />
                  </a>
                </div>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-10">
                <div>
                  <RemixIcon icon={riPriceTagLine} size="2x" className="text-green-500" />
                  <h3 className="font-bold text-gray-700 mt-3 mb-1">Pay for usage</h3>
                  <p className="text-gray-500">
                    Stop spending more for team-based subscriptions. Pay only for what you use.
                  </p>
                </div>
                <div>
                  <RemixIcon icon={riCollapseDiagonal2Line} size="2x" className="text-green-500" />
                  <h3 className="font-bold text-gray-700 mt-3 mb-1">
                    Scale down to zero (coming soon)
                  </h3>
                  <p className="text-gray-500">
                    Only need to run a few times a month? We&apos;ll scale it down to keep costs at
                    zero.
                  </p>
                </div>
              </div>
            </div>
          </Container>
        </section>
        {/* <div className="overflow-hidden max-w-screen h-32 skew-y-2 bg-gray-100 -mt-8 -mb-16" /> */}
        <section className="pt-24 pb-32 text-gray-100 bg-secondary">
          <Container>
            <div className="max-w-3xl mx-auto">
              <div>
                <h2 className="text-5xl lg:text-6xl my-8 text-primary text-center">
                  Performance at scale
                </h2>
                <p className="text-2xl text-white text-center">
                  We use languages, frameworks, and deployment strategies that are proven to scale.
                  Get the same performance with our platform as you would with a custom solution.
                </p>
                <div className="mt-10">
                  <div className="flex flex-row justify-between">
                    <div
                      className="flex items-center justify-center h-20 w-20 rounded-full bg-green-900"
                      style={{ border: '3px solid rgb(0, 172, 215)' }}
                    >
                      <Image src="/img/icons/go.png" alt="Go" width={45} height={45} />
                    </div>
                    <div
                      className="flex items-center justify-center h-20 w-20 rounded-full bg-green-900"
                      style={{ border: '3px solid rgb(50,109,230)' }}
                    >
                      <Image
                        src="/img/icons/kubernetes.png"
                        alt="Kubernetes"
                        width={45}
                        height={45}
                      />
                    </div>
                    <div
                      className="flex items-center justify-center h-20 w-20 rounded-full bg-green-900"
                      style={{ border: '3px solid rgb(75,95,171)' }}
                    >
                      <Image
                        src="/img/icons/opentelemetry.png"
                        alt="OpenTelemetry"
                        width={45}
                        height={45}
                      />
                    </div>
                  </div>
                  <svg
                    version="1.1"
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="-5 5 555 175"
                    className="mx-auto px-8"
                    style={{
                      fillRule: 'evenodd',
                      clipRule: 'evenodd',
                      strokeLinejoin: 'round',
                      strokeMiterlimit: '1.5',
                    }}
                  >
                    <g transform="matrix(1,0,0,1,-677.002,-666.378)">
                      <g transform="matrix(1,0,0,2.23643,-0.577143,-730.142)">
                        <path
                          d="M949.55,624.889L949.55,705.374"
                          style={{ fill: 'none', stroke: 'rgb(50,109,230)', strokeWidth: '3px' }}
                        />
                      </g>
                      <g transform="matrix(3.02202,0,0,2.23643,-1895.92,-730.142)">
                        <path
                          d="M1031.4,624.889C1031.4,624.889 1029.34,644.434 1020.18,652.247C998.183,670.997 972.136,670.728 959.833,680.906C948.436,690.333 949.55,705.374 949.55,705.374"
                          style={{ fill: 'none', stroke: 'rgb(75,95,171)', strokeWidth: '1px' }}
                        />
                      </g>
                      <g transform="matrix(-3.02202,0,0,2.23643,3794.93,-730.142)">
                        <path
                          d="M1031.4,624.889C1031.4,624.889 1029.34,644.434 1020.18,652.247C998.183,670.997 972.136,670.728 959.833,680.906C948.436,690.333 949.55,705.374 949.55,705.374"
                          style={{ fill: 'none', stroke: 'rgb(0,172,215)', strokeWidth: '1px' }}
                        />
                      </g>
                    </g>
                  </svg>
                  <div className="rounded-xl px-8 py-4 bg-white w-fit mx-auto">
                    <Image
                      src="/img/logos/wordmark-secondary.png"
                      alt="IOTEA"
                      width={150}
                      height={150}
                    />
                  </div>
                </div>
              </div>
            </div>
            <div className={styles.performanceGrid}>
              <div className={`${styles.performanceSection}`}>
                <h3>Production Scale</h3>
                <p>
                  Have thousands of things, channels, and requests per second? IOTEA is designed for
                  your high-throughput projects.
                </p>
                <div className={styles.throughputGrid}>
                  {Array.from({ length: 96 }).map((_, i) => (
                    <div
                      key={i}
                      className={styles.throughputDot}
                      style={{ '--delay': ((i % 24) * 96) / 24 } as React.CSSProperties}
                    />
                  ))}
                </div>
              </div>
              <div className={`${styles.performanceSection}`}>
                <h3>Realtime Capable</h3>
                <p>
                  Nodes and channels can execute in milliseconds¹, even with thousands of requests
                  per second.
                </p>
                <div className="relative">
                  <div className={styles.realtimeGraphic} />
                  <div className={styles.realtimeGraphicAnimatedLine} />
                  <span className="absolute text-gray-300 text-sm left-0 top-3">Trigger</span>
                  <span className="absolute text-gray-300 text-sm right-0 top-3">10ms</span>
                  <span className="absolute text-gray-300 text-sm right-0 top-8 md:top-10 text-xs">
                    ¹Actual execution time depends on channel setup and input data
                  </span>
                </div>
              </div>
              <div className={`${styles.performanceSection} ${styles.performanceSectionWithLine}`}>
                <h3>Complete Control</h3>
                <p>
                  Hosting your own infrastructure? Use our container-based deployment strategies for
                  total control.
                </p>
                <div className="relative flex flex-row flex-wrap justify-between">
                  <div className="mt-4 w-20 h-20 bg-white rounded flex items-center justify-center">
                    <Image
                      src="/img/icons/kubernetes.png"
                      alt="Kubernetes"
                      width={45}
                      height={45}
                    />
                  </div>
                  <div className="mt-4 w-20 h-20 bg-gray-200 rounded flex items-center justify-center">
                    <RemixIcon icon={riServerLine} size="2x" className="text-green-700" />
                  </div>
                  <div className="mt-4 w-20 h-20 bg-gray-200 rounded flex items-center justify-center">
                    <RemixIcon icon={riServerLine} size="2x" className="text-green-700" />
                  </div>
                  <div className="mt-4 w-20 h-20 bg-gray-200 rounded flex items-center justify-center hidden sm:flex">
                    <RemixIcon icon={riServerLine} size="2x" className="text-green-700" />
                  </div>
                  <div className="mt-4 w-20 h-20 bg-gray-200 rounded flex items-center justify-center hidden md:flex">
                    <RemixIcon icon={riServerLine} size="2x" className="text-green-700" />
                  </div>
                </div>
              </div>
              <div className={`${styles.performanceSection}`}>
                <h3>System Observability</h3>
                <p>
                  Our platform comes packaged with observability tools to monitor organizations,
                  channels, and spaces.
                </p>
                <div className="grid grid-cols-6 gap-2 justify-between p-2 h-20 w-full mt-4 rounded">
                  <div className="transition-all col-span-2 duration-300 bg-blue-500 w-full flex items-center justify-center rounded">
                    <RemixIcon icon={riBarChart2Fill} size="2x" className="text-white" />
                  </div>
                  <div className="transition-all col-span-2 sm:col-span-3 duration-300 bg-blue-600 w-full flex items-center justify-center rounded">
                    <RemixIcon icon={riLineChartFill} size="2x" className="text-white" />
                  </div>
                  <div className="relative -translate-y-2 col-span-2 sm:col-span-1 transition-all duration-300 bg-orange-700 w-full flex items-center justify-center rounded">
                    <RemixIcon
                      icon={riErrorWarningLine}
                      size="lg"
                      className="absolute -right-2 -top-2 text-white bg-orange-900 rounded-full p-1 shadow"
                    />
                    <RemixIcon icon={riDonutChartFill} size="2x" className="text-white" />
                  </div>
                </div>
              </div>
            </div>
          </Container>
        </section>
      </main>
    </>
  )
}

export default HomePage
