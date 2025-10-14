'use client'

import {
  IconDefinition,
  RemixIcon,
  riBracesLine,
  riCharacterRecognitionLine,
  riCodeSSlashLine,
  riDatabase2Line,
  riLoginCircleLine,
  riMailSendLine,
  riNotification3Line,
  riQuestionMark,
  riSettings3Line,
  riSettings5Line,
  riShareCircleLine,
} from '@mwarnerdotme/react-remixicon'
import { FC, useEffect, useMemo, useState } from 'react'
import styles from './index.module.scss'

type icon =
  | 'input'
  | 'processing1'
  | 'processing2'
  | 'processing3'
  | 'code'
  | 'machineLearning'
  | 'outputDatabase'
  | 'outputApi'
  | 'outputWebhook'
  | 'actionEmail'
  | 'actionNotification'

type Props = {
  icons: icon[]
  captions: string[]
  muted?: boolean
}

const NodeCard: FC<Props> = ({
  icons = [riQuestionMark],
  captions = ['Unknown'],
  muted = false,
}) => {
  const [currentIconIndex, setCurrentIconIndex] = useState<number>(0)

  const icon: IconDefinition = useMemo(() => {
    switch (icons[currentIconIndex]) {
      case 'input':
        return riLoginCircleLine
      case 'processing1':
        return riSettings5Line
      case 'processing2':
        return riSettings3Line
      case 'processing3':
      case 'code':
        return riCodeSSlashLine
      case 'machineLearning':
        return riCharacterRecognitionLine
      case 'outputDatabase':
        return riDatabase2Line
      case 'outputApi':
        return riBracesLine
      case 'outputWebhook':
        return riShareCircleLine
      case 'actionEmail':
        return riMailSendLine
      case 'actionNotification':
        return riNotification3Line
      default:
        return riQuestionMark
    }
  }, [icons, currentIconIndex])

  const caption = useMemo(() => {
    if (captions[currentIconIndex]) return captions[currentIconIndex]
    return 'Unknown'
  }, [captions, currentIconIndex])

  useEffect(() => {
    const nextIconInterval = setInterval(() => {
      setCurrentIconIndex((current) => {
        if (current + 1 > icons.length - 1) return 0
        return current + 1
      })
    }, 3e3)

    return () => clearInterval(nextIconInterval)
  }, [icons.length])

  if (!muted) {
    return (
      <div style={{ height: 75, width: 75 }} className={styles.nodeCard}>
        <RemixIcon icon={icon} size="2x" className={styles.icon} />
        <p className="text-gray-400 text-xs">{caption}</p>
      </div>
    )
  }

  return (
    <div style={{ height: 75, width: 75 }} className={`${styles.nodeCard} ${styles.muted}`}>
      <RemixIcon icon={icon} size="2x" className={styles.icon} />
      <p className="text-gray-500 text-xs">{caption}</p>
    </div>
  )
}

export default NodeCard
