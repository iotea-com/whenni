import { RemixIcon, riHome4Line, riUser3Line } from '@mwarnerdotme/react-remixicon'
import Image from 'next/image'
import Link from 'next/link'
import styles from './layout.module.scss'
import getAccessToken, { getSession } from '@iotea/hub/util/getAccessToken'
import FeedbackButton from '@iotea/hub/components/molecules/FeedbackButton'

export const metadata = {
  title: 'Profile | IOTEA',
}

const Layout = async ({ children }) => {
  const accessToken = await getAccessToken()
  const { userId } = await getSession(accessToken)
  if (!accessToken || !userId) return null

  return (
    <>
      <div className="flex grow">
        <FeedbackButton />
        <aside className="px-5 py-4 border-r border-gray-200 dark:border-gray-800">
          <nav>
            <Link href={`/dashboard`}>
              <Image
                src="/img/logos/app-icon-primary.png"
                alt="IOTEA logo"
                width={25}
                height={25}
              />
            </Link>
            <ul className="flex flex-col gap-4 mt-4">
              <li>
                <Link href={`/dashboard`}>
                  <RemixIcon icon={riHome4Line} size={'xl'} className={styles.sidebarLink} />
                </Link>
              </li>
              <li>
                <Link href={`/profiles/${userId}`}>
                  <RemixIcon icon={riUser3Line} className={styles.sidebarLink} size="lg" />
                </Link>
              </li>
            </ul>
          </nav>
        </aside>
        <main className="overflow-y-scroll overflow-y-scroll grow">
          <div className="mt-4 px-8 py-2">{children}</div>
        </main>
      </div>
    </>
  )
}

export default Layout
