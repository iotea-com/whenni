import Image from 'next/image'
import Link from 'next/link'
import {
  RemixIcon,
  riHome4Line,
  riInstanceLine,
  riKey2Line,
  riRouterLine,
  riSettings3Line,
  riShieldUserLine,
} from '@mwarnerdotme/react-remixicon'
import styles from './layout.module.scss'
import getAccessToken from '@gruent/hub/util/getAccessToken'
import FeedbackButton from '@gruent/hub/components/molecules/FeedbackButton'

export const metadata = {
  title: 'Channels | GRUENT',
}

const Layout = async ({ children, params }) => {
  const { orgId, spaceId } = await params

  const accessToken = await getAccessToken()
  if (!accessToken) return null

  return (
    <>
      <div className="flex grow">
        <FeedbackButton />
        <aside
          id="spacesNavbar"
          className="px-5 py-4 border-r border-gray-200 dark:border-gray-800"
        >
          <nav>
            <Link href={`/dashboard`}>
              <Image
                src="/img/logos/app-icon-primary.png"
                alt="GRUENT logo"
                width={25}
                height={25}
              />
            </Link>
            <ul className="flex flex-col gap-4 mt-4">
              <li id="spaceHome">
                <Link href={`/organizations/${orgId}/spaces/${spaceId}`}>
                  <RemixIcon icon={riHome4Line} size={'xl'} className={styles.sidebarLink} />
                </Link>
              </li>
              <li id="spaceThings">
                <Link href={`/organizations/${orgId}/spaces/${spaceId}/things`}>
                  <RemixIcon icon={riRouterLine} size={'xl'} className={styles.sidebarLink} />
                </Link>
              </li>
              <li id="spaceModels">
                <Link href={`/organizations/${orgId}/spaces/${spaceId}/models`}>
                  <RemixIcon icon={riInstanceLine} size={'xl'} className={styles.sidebarLink} />
                </Link>
              </li>
              <li id="spaceSecrets">
                <Link href={`/organizations/${orgId}/spaces/${spaceId}/secrets`}>
                  <RemixIcon icon={riKey2Line} size={'xl'} className={styles.sidebarLink} />
                </Link>
              </li>
              <hr />
              <li id="spaceTeam">
                <Link href={`/organizations/${orgId}/spaces/${spaceId}/team`}>
                  <RemixIcon icon={riShieldUserLine} size={'xl'} className={styles.sidebarLink} />
                </Link>
              </li>
              <li id="spaceSettings">
                <Link href={`/organizations/${orgId}/spaces/${spaceId}/settings`}>
                  <RemixIcon icon={riSettings3Line} size={'xl'} className={styles.sidebarLink} />
                </Link>
              </li>
            </ul>
          </nav>
        </aside>
        <main className="flex flex-col grow">
          <div className="grow">{children}</div>
        </main>
      </div>
    </>
  )
}

export default Layout
