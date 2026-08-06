import {
  IconDefinition,
  RemixIcon,
  riHome4Line,
  riSettings3Line,
  riShieldUserLine,
  riUser3Line,
  riUserCommunityFill,
} from '@mwarnerdotme/react-remixicon'
import Image from 'next/image'
import Link from 'next/link'
import getAccessToken, { getSession } from '@gruent/hub/util/getAccessToken'
import FeedbackButton from '@gruent/hub/components/molecules/FeedbackButton'
import gruentClient from '@gruent/hub/lib/gruent'
import OrgSelectDropdown from '@gruent/hub/components/organisms/OrgSelectDropdown'
import { FC } from 'react'

export const metadata = {
  title: 'Organization dashboard | GRUENT',
}

const Layout = async ({ children, params }) => {
  const accessToken = await getAccessToken()
  const { userId } = await getSession(accessToken)
  if (!accessToken || !userId) return null

  const { orgId } = await params

  const { data: user, errors: _getUserErrors } = await gruentClient(accessToken).users.get(userId)

  const organizations = user?.organizations.map((orgMembership) => orgMembership.organization)
  const organization = organizations?.find((org) => org.id === orgId)

  const SidebarLink: FC<{ href: string; icon: IconDefinition }> = ({ href, icon }) => {
    return (
      <Link href={href}>
        <RemixIcon
          icon={icon}
          size={'xl'}
          className="text-gray-500 hover:text-gray-900 dark:text-gray-500 dark:hover:text-gray-100 transition"
        />
      </Link>
    )
  }

  return (
    <>
      <div className="flex grow">
        <FeedbackButton />
        <aside className="px-5 py-4 border-r border-gray-200 dark:border-gray-800">
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
              <li>
                <SidebarLink href={`/organizations/${orgId}`} icon={riHome4Line} />
              </li>
              <li>
                <SidebarLink href={`/organizations/${orgId}/team`} icon={riShieldUserLine} />
              </li>
              <li>
                <SidebarLink href={`/organizations/${orgId}/settings`} icon={riSettings3Line} />
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
                <SidebarLink href={`/profiles/${userId}`} icon={riUser3Line} />
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
