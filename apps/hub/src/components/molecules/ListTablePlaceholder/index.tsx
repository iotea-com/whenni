import { PropsWithChildren } from 'react'

type Props = {
  title: string
  description: string
}

const ListTablePlaceholder = ({ title, description, children }: PropsWithChildren<Props>) => {
  return (
    <>
      <h2 className="text-lg font-extrabold">{title}</h2>
      <p className="text-gray-500 text-xs mb-4">{description}</p>
      <div className="h-72 bg-gray-100 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-md flex flex-col items-center justify-center">
        {children}
      </div>
    </>
  )
}

export default ListTablePlaceholder
