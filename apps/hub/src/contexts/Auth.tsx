import { FC, PropsWithChildren, createContext, useContext } from 'react'

type Props = {
  accessToken: string | null
}

const AuthContext = createContext<Props>({
  accessToken: null,
})

const AuthProvider: FC<PropsWithChildren<Props>> = ({ children, accessToken }) => {
  return <AuthContext.Provider value={{ accessToken }}>{children}</AuthContext.Provider>
}

export const useAuthProvider = () => {
  return useContext(AuthContext)
}

export default AuthProvider
