import LoadingSkeleton from '@gruent/hub/components/organisms/LoadingSkeleton'
import Container from '@gruent/libs/frontend/components/templates/Container'

export default function Loading() {
  return (
    <Container className="px-8 pt-10 pb-20">
      <LoadingSkeleton />
    </Container>
  )
}
