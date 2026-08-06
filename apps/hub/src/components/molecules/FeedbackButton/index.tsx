'use client'

import useSettingsStore from '@gruent/hub/stores/settingsStore'
import { RemixIcon, riMessage3Fill } from '@mwarnerdotme/react-remixicon'
import { AnimatePresence, motion } from 'framer-motion'
import Image from 'next/image'
import { FC } from 'react'

const DiscordLink = 'https://discord.gg/prAJjW426d'
// const GithubLink = 'https://github.com/ongruent/gruent'

type Props = {
  hidden?: boolean
}

const FeedbackButton: FC<Props> = ({ hidden }) => {
  const showFeedbackButton = useSettingsStore((state) => state.showFeedbackButton)

  if (hidden || !showFeedbackButton) return null

  return (
    <AnimatePresence>
      <motion.div
        className="fixed bottom-0 right-10 flex items-center gap-2 z-10"
        whileHover="hover"
        initial="initial"
      >
        <div className="bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-200 rounded-t-md py-1 px-6 select-none font-bold flex items-center">
          <a
            href={DiscordLink}
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center"
          >
            Feedback & support
            <RemixIcon
              icon={riMessage3Fill}
              size={'lg'}
              className="ml-2 text-gray-600 dark:text-gray-300"
            />
          </a>
        </div>

        <motion.a
          href={DiscordLink}
          target="_blank"
          rel="noopener noreferrer"
          className="mx-2"
          variants={{
            initial: { opacity: 0, x: -20, rotate: -90 },
            hover: { opacity: 1, x: 0, rotate: 0 },
          }}
        >
          <Image src="/img/icons/discord.svg" alt="Discord" width={24} height={24} />
        </motion.a>

        {/* TODO: enable when open source */}
        {/* <motion.a
          href={GithubLink}
          target="_blank"
          rel="noopener noreferrer"
          variants={{
            initial: { opacity: 0, x: -20, rotate: -90 },
            hover: { opacity: 1, x: 0, rotate: 0 },
          }}
          transition={{ delay: 0.1 }}
        >
          <Image src="/img/icons/github.svg" alt="Github" width={24} height={24} />
        </motion.a> */}
      </motion.div>
    </AnimatePresence>
  )
}

export default FeedbackButton
