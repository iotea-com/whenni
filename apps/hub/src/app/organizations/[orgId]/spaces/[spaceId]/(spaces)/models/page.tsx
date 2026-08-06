import gruentClient from '@gruent/hub/lib/gruent'
import getAccessToken from '@gruent/hub/util/getAccessToken'
import CreateModelModal from '@gruent/hub/components/modals/CreateModelModal'
import ModelsTable from '@gruent/hub/components/organisms/ModelsTable'
import Container from '@gruent/libs/frontend/components/templates/Container'
import Button from '@gruent/libs/frontend/components/atoms/Button'
import { RemixIcon, riAddCircleLine, riInstanceLine } from '@mwarnerdotme/react-remixicon'
import ListTablePlaceholder from '@gruent/hub/components/molecules/ListTablePlaceholder'
import Callout from '@gruent/libs/frontend/components/molecules/Callout'
import { Tag } from '@prisma/client'

export const metadata = {
  title: 'Models | GRUENT',
}

const ModelsPage = async ({ params, searchParams }) => {
  const { orgId, spaceId } = await params
  const { page: requestedPage = 1, q: filter, tagFilter: tagFilterString } = await searchParams

  const tagFilter = tagFilterString ? tagFilterString.split(',') : undefined

  const accessToken = await getAccessToken()
  if (!accessToken) return null

  const {
    data: models,
    errors,
    page,
    totalPages,
    totalResults,
  } = await gruentClient(accessToken).models.list(spaceId, {
    page: requestedPage,
    filter,
    tagFilter,
  })

  if (errors && errors.length > 0) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not list the models for this space"
          description={`${errors[0]}`}
          variant="error"
        />
      </div>
    )
  }

  if (!models) {
    return (
      <div className="mt-4 px-8 py-2">
        <Callout
          className="max-w-xl"
          title="Could not list the models for this space"
          description={`The models could not be listed. Try logging out and logging back in.`}
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
      <CreateModelModal spaceId={spaceId} />
      <Container>
        <section
          className="relative bg-gray-50 rounded-sm py-5 px-10 border border-gray-200 dark:bg-gray-900 dark:border-gray-800"
          style={{ boxShadow: '3px 3px 10px 0 rgba(0, 0, 0, 0.03)' }}
        >
          {((Number(totalResults) > 0 && !filter) || filter || tagFilter) && (
            <ModelsTable
              orgId={orgId}
              tags={tags}
              models={models}
              spaceId={spaceId}
              page={page!}
              totalPages={totalPages!}
              totalResults={totalResults ?? 0}
              searchFilter={filter}
              tagFilter={tagFilter}
            />
          )}
          {Number(totalResults) <= 0 && !filter && !tagFilter && (
            <ListTablePlaceholder title="Models" description="Define the structure of your data.">
              <RemixIcon icon={riInstanceLine} className="text-gray-600 mb-2" size="3x" />
              <p className="text-gray-600 text-center mb-2">No models... yet!</p>
              <Button variant="primary" className="text-sm" modalId="createModel">
                <RemixIcon className="mr-1" icon={riAddCircleLine} />
                <span>Create your first model</span>
              </Button>
            </ListTablePlaceholder>
          )}
        </section>
      </Container>
    </>
  )
}

export default ModelsPage
