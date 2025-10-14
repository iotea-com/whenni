'use client'

import { FC, useEffect } from 'react'
import useSettingsStore, { Theme } from '@iotea/hub/stores/settingsStore'
import useAuth from '@iotea/hub/hooks/useAuth'

const SettingsLoader: FC = () => {
  const { accessToken } = useAuth()

  const loadSettings = useSettingsStore((state) => state.loadSettings)
  const theme = useSettingsStore((state) => state.theme)
  const setTheme = useSettingsStore((state) => state.setTheme)

  useEffect(() => {
    const settings = loadSettings()
    setTheme(settings.theme)
  }, [loadSettings, setTheme])

  useEffect(() => {
    if (!accessToken) {
      if (document.documentElement.classList.contains('dark'))
        document.documentElement.classList.remove('dark')
      return
    }

    switch (theme) {
      case Theme.dark:
        if (!document.documentElement.classList.contains('dark'))
          document.documentElement.classList.add('dark')
        break
      case Theme.light:
        if (document.documentElement.classList.contains('dark'))
          document.documentElement.classList.remove('dark')
        break
    }
  }, [theme, accessToken])

  return <></>
}

export default SettingsLoader
