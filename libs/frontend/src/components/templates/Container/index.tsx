import { FC, PropsWithChildren } from 'react'

type Props = {
  className?: string
}

const Container: FC<PropsWithChildren<Props>> = ({ className, children }) => {
  return (
    <div
      className={`grid grid-cols-4 md:grid-cols-8 w-full px-10 xl:px-0 relative ${className ?? ''}`}
    >
      <div className="col-start-1 col-end-5 md:col-start-2 md:col-end-8">{children}</div>
    </div>
  )
}

export default Container
