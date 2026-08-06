'use client'

import isBrowser from '@gruent/libs/frontend/util/isBrowser'
import { useMemo, useEffect, useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'

type Props = {
  videos: {
    src: string
    hash: string
  }[]
}

const RotatingVideoPlayer = ({ videos }: Props) => {
  const [hash, setHash] = useState(isBrowser ? location.hash : '')

  useEffect(() => {
    const handleHashChange = () => {
      setHash(location.hash)
    }

    window.addEventListener('hashchange', handleHashChange)
    return () => window.removeEventListener('hashchange', handleHashChange)
  }, [])

  const currentVideoSrc = useMemo(() => {
    if (videos.length === 0) return null

    for (const video of videos) {
      if (hash.includes(video.hash)) return video.src
    }

    return videos[0].src
  }, [videos, hash])

  return (
    <div className="overflow-x-hidden w-full">
      <AnimatePresence mode="wait">
        <motion.video
          key={currentVideoSrc}
          initial={{ opacity: 0, translateX: -1000 }}
          animate={{ opacity: 1, translateX: 0 }}
          exit={{ opacity: 0, translateX: 1000 }}
          transition={{ duration: 0.5, ease: 'easeOut' }}
          src={currentVideoSrc ?? undefined}
          className="mx-auto rounded-t mt-10"
          autoPlay
          loop
          muted
          playsInline
          controls
          width={'100%'}
        />
      </AnimatePresence>
    </div>
  )
}

export default RotatingVideoPlayer
