import Container from '@gruent/libs/frontend/components/templates/Container'

const BlogDetailsLayout = ({ children }) => {
  return (
    <Container>
      <div className="mt-24 mb-20">{children}</div>
    </Container>
  )
}

export default BlogDetailsLayout
