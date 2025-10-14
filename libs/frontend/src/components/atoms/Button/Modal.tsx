'use client'

import { FC } from 'react'
import Button, { Props as BaseProps } from '.'
import { openModal } from '@iotea/libs/frontend/hooks/useModal'

type Props = Omit<BaseProps, 'onClick' | 'href'> & {
  id: string
}

const ModalButton: FC<Props> = ({ id, text, variant, className, style }) => {
  return (
    <Button
      text={text}
      style={style}
      variant={variant}
      className={className}
      onClick={() => openModal(id)}
    />
  )
}

export default ModalButton
