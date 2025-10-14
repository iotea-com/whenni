'use client'

// import { Session } from 'next-auth'
// // import { useSession } from 'next-auth/react'
// import { useRouter, useSearchParams } from 'next/navigation'
// import { FC, useEffect } from 'react'

// type Props = {
//   session: Session | null
// }

// // Redirects to the callbackUrl if it exists, otherwise redirects to the dashboard.
// // This client-side component is required since layouts cannot receive search params.
// const CallbackUrlRedirect: FC<Props> = ({ session: serverSession }) => {
//   const router = useRouter()
//   const searchParams = useSearchParams()
//   const callbackUrl = searchParams.get('callbackUrl')

//   // const { data: clientSession } = useSession()

//   useEffect(() => {
//     if (serverSession && clientSession) router.push(callbackUrl || '/dashboard')
//   }, [router, callbackUrl, serverSession, clientSession])

//   return null
// }

// export default CallbackUrlRedirect

export default null
