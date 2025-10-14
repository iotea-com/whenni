import LoadingSkeleton from '@iotea/hub/components/organisms/LoadingSkeleton'
import Container from '@iotea/libs/frontend/components/templates/Container'

export default function Loading() {
  return (
    <Container className="px-8 pt-10 pb-20">
      <LoadingSkeleton />
    </Container>
  )
}
