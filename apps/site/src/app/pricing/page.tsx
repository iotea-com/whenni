import Container from '@gruent/libs/frontend/components/templates/Container'
import ChannelsCostCalculator from '@gruent/libs/frontend/components/organisms/GridExamplePriceCalculator'
import {
  RemixIcon,
  ri24HoursFill,
  riApps2AiFill,
  riArrowLeftRightFill,
  riArrowRightLine,
  riArticleFill,
  riCloudFill,
  riDragDropFill,
  riGitForkFill,
  riInstanceFill,
  riKey2Fill,
  riLineChartFill,
  riLock2Fill,
  riRouterLine,
  riServerFill,
  riShieldUserFill,
  riSpeedMiniFill,
  riStackFill,
  riSwap3Fill,
  riTeamFill,
} from '@mwarnerdotme/react-remixicon'

import styles from './page.module.scss'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import { Metadata } from 'next'
import RuntimeSelector from './RuntimeSelector'

export const metadata: Metadata = {
  title: 'Usage-based pricing for your IoT & automation needs | GRUENT',
  description:
    'Pay only for what you use. Predictable pricing without hidden fees. Get started in minutes.',
  openGraph: {
    type: 'website',
    url: `https://gruent.com/pricing`,
    title: 'Usage-based pricing for your IoT & automation needs | GRUENT',
    description:
      'Pay only for what you use. Predictable pricing without hidden fees. Get started in minutes.',
    images: ['https://gruent.com/img/logos/app-icon-primary.png'],
  },
}

const Pricing = async ({ searchParams }) => {
  const { runtime = 'sm' } = await searchParams

  const prices = {
    xs: {
      base: 2,
      executions: 1,
      connections: '∞',
      executionsPerSecond: 3,
    },
    sm: {
      base: 5,
      executions: 0.5,
      connections: '∞',
      executionsPerSecond: 50,
    },
    md: {
      base: 10,
      executions: 0.25,
      connections: '∞',
      executionsPerSecond: 200,
    },
    lg: {
      base: 15,
      executions: 0.1,
      connections: '∞',
      executionsPerSecond: 500,
    },
  }

  return (
    <div className="mt-[75px]">
      <header className="relative">
        <div className="overflow-hidden w-screen z-0 absolute right-0">
          <div
            style={{
              borderLeft: '3000px solid #0d3e25',
              borderBottom: '700px solid transparent',
            }}
          />
        </div>
        <div className="z-10 relative pt-24 pb-14">
          <Container>
            <h1 className="text-center text-primary text-4xl">
              Usage-based pricing, for everyone.
            </h1>
            <h2 className="text-center text-white text-2xl font-light mt-2 mb-12">
              Pay only for what you use. Predictable pricing without hidden fees.
            </h2>
            <div className="grid lg:grid-cols-2 gap-0">
              <div
                className="p-[2px] w-fit h-fit mx-auto mt-10 shadow-lg"
                style={{
                  background: 'linear-gradient(135deg, rgb(77 255 104), #195E13)',
                }}
              >
                <div className="bg-white" style={{ maxWidth: 600 }}>
                  <div className="pt-6 pb-4">
                    <h3 className="text-center text-green-600 text-lg uppercase">Fully managed</h3>
                  </div>
                  <hr />
                  <div className="px-14 py-8">
                    <h4 className="text-center mb-8 font-normal text-xl text-gray-600">
                      Access a complete IoT and automation platform with usage-based pricing
                      starting at
                    </h4>
                    <div className="grid grid-cols-3 gap-4 mb-10">
                      <div>
                        <p className="text-4xl md:text-5xl text-gray-800 text-center font-title mb-1">
                          ${prices[runtime].executions}
                        </p>
                        <p className="text-center text-gray-500 text-sm">
                          per 1,000 channel executions
                        </p>
                        <p className="text-center text-gray-500 text-xs font-bold">
                          Minimum ${prices[runtime].base}/mo
                        </p>
                        <p className="text-center text-gray-500 text-xs font-bold">
                          Bulk discounts available
                        </p>
                      </div>
                      <div>
                        <p className="text-4xl md:text-5xl text-gray-800 text-center font-title mb-1">
                          1¢
                        </p>
                        <p className="text-center text-gray-500 text-sm">per GB ingress data</p>
                        <p className="text-center text-gray-500 text-xs font-bold">
                          Bulk discounts available
                        </p>
                      </div>
                      <div>
                        <p className="text-4xl md:text-5xl text-gray-800 text-center font-title mb-1">
                          15¢
                        </p>
                        <p className="text-center text-gray-500 text-sm">per GB egress data</p>
                        <p className="text-center text-gray-500 text-xs font-bold">
                          Bulk discounts available
                        </p>
                      </div>
                    </div>
                    <RuntimeSelector runtime={runtime} />
                    <ul className="text-gray-500 text-lg flex flex-col gap-2 mx-auto w-fit">
                      <li className="flex items-center">
                        <RemixIcon
                          icon={riDragDropFill}
                          size="xl"
                          className="text-green-700 mr-2"
                        />{' '}
                        Everything you need to connect{' '}
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
                        </div>
                      </li>
                      <li className="flex items-center">
                        <RemixIcon
                          icon={riSpeedMiniFill}
                          size="xl"
                          className="text-green-700 mr-2"
                        />{' '}
                        Up to {prices[runtime].executionsPerSecond} executions per second
                      </li>
                      <li className="flex items-center">
                        <RemixIcon icon={ri24HoursFill} size="xl" className="text-green-700 mr-2" />{' '}
                        99.99% uptime SLA
                      </li>
                    </ul>
                  </div>
                  <p className="text-gray-400 mb-2 mx-14 text-center">
                    <small>Prices subject to change throughout beta</small>
                  </p>
                  <Button
                    variant="blank"
                    modalId="betaSignup"
                    className="py-4 px-14 mx-auto w-full bg-gray-100 border-0"
                  >
                    Get started in minutes <RemixIcon icon={riArrowRightLine} />
                  </Button>
                </div>
              </div>
              <div className="p-[2px] w-fit h-fit mx-auto mt-10" style={{ background: '#0d3e25' }}>
                <div className="bg-white" style={{ maxWidth: 400 }}>
                  <div className="pt-6 pb-4">
                    <h3 className="text-center text-green-600 text-lg uppercase">Self hosted</h3>
                  </div>
                  <hr className="border-gray-300" />
                  <div className="px-14 py-8">
                    <h4 className="text-center mb-8 font-normal text-xl text-gray-700">
                      Host your own setup for the most control
                    </h4>
                    <ul className="text-gray-500 text-lg flex flex-col gap-2 mx-auto w-fit">
                      <li className="flex items-center">
                        <RemixIcon icon={riServerFill} size="xl" className="text-green-700 mr-2" />{' '}
                        Container-based deployments
                      </li>
                      <li className="flex items-center">
                        <RemixIcon icon={riStackFill} size="xl" className="text-green-700 mr-2" />{' '}
                        Maximize customizability
                      </li>
                      <li className="flex items-center">
                        <RemixIcon icon={riCloudFill} size="xl" className="text-green-700 mr-2" />{' '}
                        Optional data forwarding
                      </li>
                    </ul>
                  </div>
                  <div className="bg-gray-100 py-4 px-14 text-center">
                    <p className="text-lg">Coming soon</p>
                  </div>
                </div>
              </div>
            </div>
          </Container>
        </div>
      </header>
      <main>
        <section className="mt-28">
          <Container>
            <h2 className="text-center text-4xl w-full">Base Features</h2>
            <p className="text-center w-full text-2xl mt-1 text-gray-600">
              Everything needed to get started with the Internet of Things.
            </p>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-10 mt-10">
              <div>
                <RemixIcon icon={riDragDropFill} size="2x" className="text-green-700" />
                <h3 className="text-lg my-4">Drag and drop configurations</h3>
                <ul className="text-gray-500 ml-4" style={{ listStyleType: 'disclosure-closed' }}>
                  <li>Supports any major protocol</li>
                  <li>Up to 100 nodes in a channel</li>
                  <li>Deploys within seconds</li>
                </ul>
              </div>
              <div>
                <RemixIcon icon={riRouterLine} size="2x" className="text-green-700" />
                <h3 className="text-lg my-4">Simple service and device management</h3>
                <ul className="text-gray-500 ml-4" style={{ listStyleType: 'disclosure-closed' }}>
                  <li>Create and manage up to 10,000 things</li>
                  <li>Thing categories and tags</li>
                  <li>Use across multiple channels</li>
                </ul>
              </div>
              <div>
                <RemixIcon icon={riArrowLeftRightFill} size="2x" className="text-green-700" />
                <h3 className="text-lg my-4">Managed MQTT broker</h3>
                <ul className="text-gray-500 ml-4" style={{ listStyleType: 'disclosure-closed' }}>
                  <li>Send messages to any connected device</li>
                  <li>Certificate-based authentication</li>
                  <li>Realtime speeds with stable connections</li>
                </ul>
              </div>
              <div>
                <RemixIcon icon={riInstanceFill} size="2x" className="text-green-700" />
                <h3 className="text-lg my-4">Define your own data models</h3>
                <ul className="text-gray-500 ml-4" style={{ listStyleType: 'disclosure-closed' }}>
                  <li>Ensure data integrity</li>
                  <li>Use between any node</li>
                  <li>Import common models from the community</li>
                </ul>
              </div>
              <div>
                <RemixIcon icon={riApps2AiFill} size="2x" className="text-green-700" />
                <h3 className="text-lg my-4">AI assistant</h3>
                <ul className="text-gray-500 ml-4" style={{ listStyleType: 'disclosure-closed' }}>
                  <li>Build solutions faster</li>
                  <li>Debug and fix channel errors</li>
                  <li>Hardware and integration recommendations</li>
                </ul>
              </div>
              <div>
                <RemixIcon icon={riLock2Fill} size="2x" className="text-green-700" />
                <h3 className="text-lg my-4">Security by default</h3>
                <ul className="text-gray-500 ml-4" style={{ listStyleType: 'disclosure-closed' }}>
                  <li>Authentication enforcement and policies</li>
                  <li>Secrets for sensitive information</li>
                  <li>PII redaction in logs and traces</li>
                </ul>
              </div>
              <div>
                <RemixIcon icon={riArticleFill} size="2x" className="text-green-700" />
                <h3 className="text-lg my-4">Comprehensive developer resources</h3>
                <ul className="text-gray-500 ml-4" style={{ listStyleType: 'disclosure-closed' }}>
                  <li>Documentation and guides</li>
                  <li>SDKs & libraries in several languages</li>
                  <li>Plug-and-play hardware</li>
                </ul>
              </div>
              <div>
                <RemixIcon icon={riTeamFill} size="2x" className="text-green-700" />
                <h3 className="text-lg my-4">Bring your whole team</h3>
                <ul className="text-gray-500 ml-4" style={{ listStyleType: 'disclosure-closed' }}>
                  <li>Up to 10,000 team members</li>
                  <li>Granular role-based authorization</li>
                  <li>Optional multi-factor authentication</li>
                </ul>
              </div>
              <div>
                <RemixIcon icon={riLineChartFill} size="2x" className="text-green-700" />
                <h3 className="text-lg my-4">Observability and analytics</h3>
                <ul className="text-gray-500 ml-4" style={{ listStyleType: 'disclosure-closed' }}>
                  <li>View uptime and requests per second</li>
                  <li>Logs and traces for every channel</li>
                  <li>Set alerts for important events</li>
                </ul>
              </div>
              {/* <div>
                <RemixIcon icon={riShieldCheckFill} size="2x" className="text-green-700" />
                <h3 className="text-lg my-4">Advanced security & compliance</h3>
                <ul className="text-gray-500 ml-4" style={{ listStyleType: 'disclosure-closed' }}>
                  <li>SOC 2 compliant</li>
                  <li>Messages encrypted in transit and at rest</li>
                  <li>Monitoring and logs</li>
                </ul>
              </div> */}
              {/* <div>
                <RemixIcon icon={riRouterFill} size="2x" className="text-green-700" />
                <h3 className="text-lg my-4">24x7 service availability</h3>
                <ul className="text-gray-500 ml-4" style={{ listStyleType: 'disclosure-closed' }}>
                  <li>99.999% uptime</li>
                  <li>Monitoring and logs</li>
                  <li>Advanced troubleshooting tools</li>
                </ul>
              </div> */}
            </div>
          </Container>
        </section>
        <section className="mt-28 mb-20">
          <Container>
            <h2 className="text-4xl text-center">Integrated Services</h2>
            <p className="text-2xl text-center mt-1 text-gray-600">
              Go further with data storage, device location, machine learning, and more
            </p>
            <div className="gap-6">
              <div className="px-10 py-6 shadow-md border border-green-500 mt-10">
                <h3 className="text-2xl">Hub</h3>
                <h4 className="text-gray-600 font-normal">
                  A complete Internet of Things platform
                </h4>
                <div className="flex flex-col md:flex-row gap-4 mt-4">
                  <div className="grow">
                    <h4 className="flex items-center gap-1 text-lg mt-4 font-normal">
                      <RemixIcon icon={riGitForkFill} size="lg" className="text-green-700" />{' '}
                      Channels
                    </h4>
                    <p className="text-gray-500 max-w-lg">
                      Connect anything with channels. Define inputs from a variety of sources,
                      process them, and send the results or perform actions.
                    </p>
                  </div>
                  <div
                    className="border-l border-gray-200 flex flex-col justify-center pl-6 w-64"
                    style={{ minWidth: '16rem' }}
                  >
                    <h5 className="text-xl">$1 to $0.0005</h5>
                    <p className="text-gray-400">per 1,000 executions</p>
                    <h5 className="text-xl mt-2">Up to 10,000</h5>
                    <p className="text-gray-400">channels per space</p>
                    <h5 className="text-xl mt-2">Up to 100</h5>
                    <p className="text-gray-400">nodes per channel</p>
                  </div>
                </div>
                <div className="flex flex-col md:flex-row gap-4 mt-4">
                  <div className="grow">
                    <h4 className="flex items-center gap-1 text-lg mt-4 font-normal">
                      <RemixIcon icon={riRouterLine} size="lg" className="text-green-700" /> Things
                    </h4>
                    <p className="text-gray-500 max-w-lg">
                      Integrate devices, APIs, third-party services, databases, messages queues, and
                      more.
                    </p>
                  </div>
                  <div
                    className="border-l border-gray-200 flex flex-col justify-center pl-6 w-64"
                    style={{ minWidth: '16rem' }}
                  >
                    <h5 className="text-xl">Included</h5>
                    <p className="text-gray-400">up to 10,000 things</p>
                  </div>
                </div>
                <div className="flex flex-col md:flex-row gap-4 mt-4">
                  <div className="grow">
                    <h4 className="flex items-center gap-1 text-lg mt-4 font-normal">
                      <RemixIcon icon={riInstanceFill} size="lg" className="text-green-700" />{' '}
                      Models
                    </h4>
                    <p className="text-gray-500 max-w-lg">
                      Define your own data models to ensure data integrity and consistency.
                    </p>
                  </div>
                  <div
                    className="border-l border-gray-200 flex flex-col justify-center pl-6 w-64"
                    style={{ minWidth: '16rem' }}
                  >
                    <h5 className="text-xl">Included</h5>
                    <p className="text-gray-400">up to 10,000 models</p>
                  </div>
                </div>
                <div className="flex flex-col md:flex-row gap-4 mt-4">
                  <div className="grow">
                    <h4 className="flex items-center gap-1 text-lg mt-4 font-normal">
                      <RemixIcon icon={riShieldUserFill} size="lg" className="text-green-700" />{' '}
                      Team
                    </h4>
                    <p className="text-gray-500 max-w-lg">
                      Team members to collaborate in your organization or space.
                    </p>
                  </div>
                  <div
                    className="border-l border-gray-200 flex flex-col justify-center pl-6 w-64"
                    style={{ minWidth: '16rem' }}
                  >
                    <h5 className="text-xl">Included</h5>
                    <p className="text-gray-400">up to 10,000 members</p>
                  </div>
                </div>
                <div className="flex flex-col md:flex-row gap-4 mt-4">
                  <div className="grow">
                    <h4 className="flex items-center gap-1 text-lg mt-4 font-normal">
                      <RemixIcon icon={riKey2Fill} size="lg" className="text-green-700" /> Secrets
                    </h4>
                    <p className="text-gray-500 max-w-lg">
                      Securely store and manage secrets for your channels and things.
                    </p>
                  </div>
                  <div
                    className="border-l border-gray-200 flex flex-col justify-center pl-6 w-64"
                    style={{ minWidth: '16rem' }}
                  >
                    <h5 className="text-xl">Included</h5>
                    <p className="text-gray-400">up to 10,000 secrets</p>
                  </div>
                </div>
                <div className="flex flex-col md:flex-row gap-4 mt-4">
                  <div className="grow">
                    <h4 className="flex items-center gap-1 text-lg mt-4 font-normal">
                      <RemixIcon icon={riSwap3Fill} size="lg" className="text-green-700" /> API
                    </h4>
                    <p className="text-gray-500 max-w-lg">
                      Programmatic access to all GRUENT API routes.
                    </p>
                  </div>
                  <div
                    className="border-l border-gray-200 flex flex-col justify-center pl-6 w-64"
                    style={{ minWidth: '16rem' }}
                  >
                    <h5 className="text-xl">Included</h5>
                    <p className="text-gray-400">up to 1 request/second</p>
                  </div>
                </div>
                <hr className="my-8" />
                <div>
                  <h4 className="font-normal text-base mb-1 font-normal">Price calculator</h4>
                  <ChannelsCostCalculator runtime={runtime} baseCost={prices[runtime].base} />
                </div>
              </div>
              <div className="px-10 py-6 shadow-md border border-gray-100 mt-6">
                <h3 className="text-2xl">Managed things</h3>
                <h4 className="text-gray-600 font-normal">
                  Built-in databases, file storage, and message brokers
                </h4>
                <p className="text-gray-600">Coming soon!</p>
              </div>
              <div className="px-10 py-6 shadow-md border border-gray-100 mt-6">
                <h3 className="text-2xl">Advanced devices</h3>
                <h4 className="text-gray-600 font-normal">
                  Level up your devices with device location, over-the-air updates, and digital
                  twins
                </h4>
                <p className="text-gray-600">Coming soon!</p>
              </div>
              <div className="px-10 py-6 shadow-md border border-gray-100 mt-6">
                <h3 className="text-2xl">Machine Learning</h3>
                <h4 className="text-gray-600 font-normal">
                  Train and query machine learning models with data
                </h4>
                <p>Coming soon!</p>
              </div>
            </div>
          </Container>
        </section>
      </main>
    </div>
  )
}

export default Pricing
