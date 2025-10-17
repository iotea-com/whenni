import ioteaClient from '@iotea/hub/lib/iotea'
import Image from 'next/image'
import Link from 'next/link'
import {
  RemixIcon,
  riCircleLine,
  riHome4Line,
  riInstanceLine,
  riKey2Line,
  riRouterLine,
  riSettings3Line,
  riUserCommunityFill,
  riShieldUserLine,
} from '@mwarnerdotme/react-remixicon'
import SpaceSelectDropdown from '@iotea/hub/components/organisms/SpaceSelectDropdown'
import styles from './layout.module.scss'
import getAccessToken from '@iotea/hub/util/getAccessToken'
import FeedbackButton from '@iotea/hub/components/molecules/FeedbackButton'

export const metadata = {
  title: 'Space dashboard | IOTEA',
}

const Layout = async ({ children, params }) => {
  const { orgId, spaceId } = await params

  const accessToken = await getAccessToken()
  if (!accessToken) return null

  const { data: space, errors: _spaceErrors } = await ioteaClient(accessToken).spaces.get(
    orgId,
    spaceId,
  )

  const { data: organization, errors: _organizationErrors } = await ioteaClient(
    accessToken,
  ).organizations.get(space?.organizationId ?? '')

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
                alt="IOTEA logo"
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
          <header className="flex border-b border-gray-200 dark:border-gray-800 px-8 py-4 h-14">
            <div className="flex gap-2 items-center w-full text-gray-700 dark:text-gray-300 text-sm">
              <Link
                href={`/organizations/${orgId}`}
                className="flex items-center gap-1 transition text-gray-700 hover:text-gray-800 dark:text-gray-300 dark:hover:text-gray-200"
              >
                <div className="w-5 h-5 flex items-center justify-center bg-gray-100 dark:bg-gray-800 rounded-xs">
                  <RemixIcon icon={riUserCommunityFill} size={'sm'} />
                </div>
                {organization?.name ?? 'Error'}
              </Link>
              <span>/</span>
              <Link
                href={`/organizations/${orgId}/spaces/${spaceId}`}
                className="flex items-center gap-1 transition text-gray-700 hover:text-gray-800 dark:text-gray-300 dark:hover:text-gray-200"
              >
                <div className="w-5 h-5 flex items-center justify-center bg-gray-100 dark:bg-gray-800 rounded-xs">
                  <RemixIcon icon={riCircleLine} size={'sm'} />
                </div>
                <SpaceSelectDropdown
                  organization={organization ?? undefined}
                  currentSpaceId={spaceId}
                />
              </Link>
            </div>
          </header>
          <div className="grow px-8 pt-10 pb-20">{children}</div>
        </main>
      </div>
    </>
  )
}

export default Layout
