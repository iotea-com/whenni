import CreateOrganizationModal from '@iotea/hub/components/modals/CreateOrganizationModal'
import { RemixIcon, riAddCircleLine, riArrowRightLine } from '@mwarnerdotme/react-remixicon'
import Link from 'next/link'
import ioteaClient from '@iotea/hub/lib/iotea'
import getAccessToken, { getSession } from '@iotea/hub/util/getAccessToken'
import Button from '@iotea/libs/frontend/components/atoms/Button'
import Callout from '@iotea/libs/frontend/components/molecules/Callout'
export const metadata = {
  title: 'Dashboard | IOTEA',
}

const DashboardPage = async () => {
  const accessToken = await getAccessToken()
  const { userId } = await getSession(accessToken)

  if (!accessToken || !userId) return null

  const { data: user, errors } = await ioteaClient(accessToken).users.get(userId)

  if (errors && errors.length > 0) {
    return (
      <div>
        <Callout
          className="max-w-xl"
          title="Could not retrieve the user details"
          description={`${errors[0]}`}
          variant="error"
        />
      </div>
    )
  }

  if (!user) {
    return (
      <div>
        <Callout
          className="max-w-xl"
          title="User not found"
          description={`This user does not exist. Try logging out and logging back in.`}
          variant="error"
        />
      </div>
    )
  }

  return (
    <div>
      <CreateOrganizationModal />
      <h1 className="text-2xl mb-6">Welcome to IOTEA</h1>
      <section id="spaces">
        <div className="flex items-end mb-4">
          <div>
            <h2>Organizations</h2>
            <p>
              <small>Collaborate with your team across spaces.</small>
            </p>
          </div>
          <div className="grow" />
          <Button variant="transparent" className="text-sm" modalId="createOrganization">
            <RemixIcon className="mr-1" icon={riAddCircleLine} />
            <span>Create</span>
          </Button>
        </div>
        {user.organizations && user.organizations.length > 0 && (
          <div className="grid grid-cols-3 gap-8">
            {user.organizations.map((organizationMembership) => {
              return (
                <Link
                  key={organizationMembership.organization.id}
                  className="spaceCard"
                  href={`/organizations/${organizationMembership.organization.id}`}
                >
                  <div className="flex px-4 py-6 bg-gray-100 hover:bg-gray-200 dark:bg-gray-800 dark:hover:bg-gray-700 border border-gray-200 dark:border-gray-700 transition rounded-md">
                    <h2>{organizationMembership.organization.name}</h2>
                    <div className="grow" />
                    <RemixIcon icon={riArrowRightLine} className="text-gray-600" />
                  </div>
                </Link>
              )
            })}
          </div>
        )}
        {(!user.organizations || user.organizations.length <= 0) && (
          <p>You're not a member of any organizations... yet!</p>
        )}
      </section>
    </div>
  )
}

export default DashboardPage
