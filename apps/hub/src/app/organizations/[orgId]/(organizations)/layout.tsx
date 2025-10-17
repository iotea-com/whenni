import {
  RemixIcon,
  riHome4Line,
  riSettings3Line,
  riShieldUserLine,
  riUser3Line,
  riUserCommunityFill,
} from '@mwarnerdotme/react-remixicon'
import Image from 'next/image'
import Link from 'next/link'
import styles from './layout.module.scss'
import getAccessToken, { getSession } from '@iotea/hub/util/getAccessToken'
import FeedbackButton from '@iotea/hub/components/molecules/FeedbackButton'
import ioteaClient from '@iotea/hub/lib/iotea'
import OrgSelectDropdown from '@iotea/hub/components/organisms/OrgSelectDropdown'

export const metadata = {
  title: 'Organization dashboard | IOTEA',
}

const Layout = async ({ children, params }) => {
  const accessToken = await getAccessToken()
  const { userId } = await getSession(accessToken)
  if (!accessToken || !userId) return null

  const { orgId } = await params

  const { data: user, errors: _getUserErrors } = await ioteaClient(accessToken).users.get(userId)

  const organizations = user?.organizations.map((orgMembership) => orgMembership.organization)
  const organization = organizations?.find((org) => org.id === orgId)

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
                <Link href={`/organizations/${orgId}`}>
                  <RemixIcon icon={riHome4Line} size={'xl'} className={styles.sidebarLink} />
                </Link>
              </li>
              <li>
                <Link href={`/organizations/${orgId}/team`}>
                  <RemixIcon icon={riShieldUserLine} size={'xl'} className={styles.sidebarLink} />
                </Link>
              </li>
              <li>
                <Link href={`/organizations/${orgId}/settings`}>
                  <RemixIcon icon={riSettings3Line} size={'xl'} className={styles.sidebarLink} />
                </Link>
              </li>
            </ul>
          </nav>
        </aside>
        <main className="overflow-y-scroll overflow-y-scroll grow">
          <header className="flex border-b border-gray-200 dark:border-gray-800 px-8 py-4 h-14">
            <div className="flex gap-2 items-center w-full text-gray-700 dark:text-gray-300 text-sm">
              <div className="w-5 h-5 flex items-center justify-center bg-gray-100 dark:bg-gray-800 rounded-xs">
                <RemixIcon icon={riUserCommunityFill} size={'sm'} />
              </div>
              <OrgSelectDropdown organizations={organizations ?? []} currentOrg={organization} />
              <div className="grow" />
              <div className="flex gap-4">
                <Link href={`/profiles/${userId}`}>
                  <RemixIcon icon={riUser3Line} className={styles.sidebarLink} size="lg" />
                </Link>
              </div>
            </div>
          </header>
          <div className="grow px-8 pt-10 pb-20">{children}</div>
        </main>
      </div>
    </>
  )
}

export default Layout
