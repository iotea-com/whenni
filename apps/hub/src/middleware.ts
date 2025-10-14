import { NextRequest } from 'next/server'

export const middleware = async (req: NextRequest) => {
  const accessToken = req.cookies.get('access_token')

  // Redirect to signin if the client is not signed in and the route is protected
  if (!accessToken) {
    const newUrl = new URL('/signin', req.nextUrl.origin)
    return Response.redirect(newUrl)
  }
}

export const config = {
  matcher: ['/dashboard/:path*', '/profiles/:path*', '/organizations/:path*'],
}
