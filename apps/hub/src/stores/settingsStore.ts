import { create } from 'zustand'

export enum Theme {
  light = 'light',
  dark = 'dark',
}

type State = {
  theme: Theme
  showChannelEditorGrid: boolean
  showFeedbackButton: boolean
}

type Action = {
  loadSettings: () => Settings
  setTheme: (theme: Theme) => void
  setShowChannelEditorGrid: (showChannelEditorGrid: boolean) => void
  setShowFeedbackButton: (showFeedbackButton: boolean) => void
}

type Settings = State

const DefaultSettings: Settings = {
  theme: Theme.light,
  showChannelEditorGrid: true,
  showFeedbackButton: true,
}

const updateSettingsLocalStorage = (settings: Settings) => {
  const settingsLocalStorage = {
    theme: settings.theme,
    showChannelEditorGrid: settings.showChannelEditorGrid,
    showFeedbackButton: settings.showFeedbackButton,
  }

  localStorage.setItem('settings', JSON.stringify(settingsLocalStorage))
}

const loadSettingsLocalStorage = () => {
  const settingsLocalStorage = localStorage.getItem('settings')

  // Initialize local storage settings if none are found
  if (!settingsLocalStorage) {
    updateSettingsLocalStorage(DefaultSettings)
    return DefaultSettings
  }

  const settings = JSON.parse(settingsLocalStorage) as Settings
  return settings
}

const useSettingsStore = create<State & Action>((set, get) => ({
  ...DefaultSettings,
  loadSettings: loadSettingsLocalStorage,
  setTheme: (theme) => {
    set({ theme })
    if (document) document.documentElement.setAttribute('data-theme', theme)
    const currentState = get()
    updateSettingsLocalStorage(currentState)
  },
  setShowChannelEditorGrid: (showChannelEditorGrid) => {
    set({ showChannelEditorGrid })
    const currentState = get()
    updateSettingsLocalStorage(currentState)
  },
  setShowFeedbackButton: (showFeedbackButton) => {
    set({ showFeedbackButton })
    const currentState = get()
    updateSettingsLocalStorage(currentState)
  },
}))

export default useSettingsStore
