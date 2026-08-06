import gruentClient from '@gruent/hub/lib/gruent'
import CreateThingModal from '@gruent/hub/components/modals/CreateThingModal'
import ThingsTable from '@gruent/hub/components/organisms/ThingsTable'
import getAccessToken from '@gruent/hub/util/getAccessToken'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import { RemixIcon, riAddCircleLine, riRouterLine } from '@mwarnerdotme/react-remixicon'
import Container from '@gruent/libs/frontend/components/templates/Container'
import ListTablePlaceholder from '@gruent/hub/components/molecules/ListTablePlaceholder'
import Callout from '@gruent/libs/frontend/components/molecules/Callout'
import { Tag } from '@prisma/client'

export const metadata = {
  title: 'Things | GRUENT',
}

const ThingsPage = async ({ params, searchParams }) => {
  const { orgId, spaceId } = await params
  const { page: requestedPage = 1, q: filter, tagFilter: tagFilterString } = await searchParams

  const tagFilter = tagFilterString ? tagFilterString.split(',') : undefined

  const accessToken = await getAccessToken()
  if (!accessToken) return null

  const {
    data: things,
    errors: thingListErrors,
    page,
    totalPages,
    totalResults,
  } = await gruentClient(accessToken).things.list(spaceId, {
    page: requestedPage,
    filter,
    tagFilter,
  })

  if (thingListErrors && thingListErrors.length > 0) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not list the things for this space"
          description={`${thingListErrors[0]}`}
          variant="error"
        />
      </div>
    )
  }

  if (!things) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not list the things for this space"
          description={`The things could not be listed. Try logging out and logging back in.`}
          variant="error"
        />
      </div>
    )
  }

  const { data: tagListResponse, errors: tagListErrors } =
    await gruentClient(accessToken).tags.list(spaceId)

  if (tagListErrors && tagListErrors.length > 0) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not list the tags for this space"
          description={`${tagListErrors[0]}`}
          variant="error"
        />
      </div>
    )
  }

  if (!tagListResponse) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not list the tags for this space"
          description={`The tags could not be listed. Try logging out and logging back in.`}
          variant="error"
        />
      </div>
    )
  }

  const tags: Tag[] = []
  for (const [_, v] of Object.entries(tagListResponse)) {
    tags.push(v.tag)
  }

  return (
    <>
      <CreateThingModal orgId={orgId} spaceId={spaceId} />
      <Container>
        <section className="relative bg-gray-50 rounded-sm py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800">
          {((Number(totalResults) > 0 && !filter) || filter || tagFilter) && (
            <ThingsTable
              tags={tags ?? []}
              searchFilter={filter}
              tagFilter={tagFilter}
              orgId={orgId}
              things={things}
              spaceId={spaceId}
              page={page!}
              totalPages={totalPages!}
              totalResults={totalResults ?? 0}
            />
          )}
          {Number(totalResults) <= 0 && !filter && !tagFilter && (
            <ListTablePlaceholder
              title="Things"
              description="Devices, APIs, databases, message queues - anything you need to connect."
            >
              <RemixIcon icon={riRouterLine} className="text-gray-600 mb-2" size="3x" />
              <p className="text-gray-600 text-center mb-2">No things... yet!</p>
              <Button variant="primary" className="text-sm" modalId="createThing">
                <RemixIcon className="mr-1" icon={riAddCircleLine} />
                <span>Create your first thing</span>
              </Button>
            </ListTablePlaceholder>
          )}
        </section>
      </Container>
    </>
  )
}

export default ThingsPage
