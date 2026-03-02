'use client'

import { ImageBlockObjectResponse } from '@notionhq/client/build/src/api-endpoints'

const ImageModal = ({ imageUrl, caption }: { imageUrl: string; caption?: string }) => {
  return (
    <dialog
      className="backdrop:bg-gray-900 backdrop:bg-opacity-90 outline-hidden border-none rounded-sm open:flex open:items-center open:justify-center"
      onClick={(e) => {
        e.currentTarget.close()
      }}
    >
      <div className="relative max-w-[90vw] max-h-[90vh] flex flex-col items-center">
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          src={imageUrl}
          alt={caption || ''}
          className="max-w-full max-h-[80vh] object-contain"
        />
        {caption && <p className="text-white text-sm mt-2">{caption}</p>}
      </div>
    </dialog>
  )
}

const ImageBlock = ({ block }: { block: ImageBlockObjectResponse }) => {
  const imageUrl = (() => {
    if (block.image.type === 'file') return block.image.file.url
    if (block.image.type === 'external') return block.image.external.url
    return null
  })()

  if (imageUrl) {
    const caption = block.image.caption?.[0]?.plain_text
    return (
      <>
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          key={block.id}
          src={imageUrl}
          alt={caption || ''}
          className="rounded-md shadow-md cursor-zoom-in"
          onClick={(e) => {
            const dialog = e.currentTarget.nextElementSibling as HTMLDialogElement
            dialog?.showModal()
          }}
        />
        {caption && <p className="text-sm text-gray-600 my-1">{caption}</p>}
        <ImageModal imageUrl={imageUrl} caption={caption} />
      </>
    )
  }

  return null
}

export default ImageBlock
