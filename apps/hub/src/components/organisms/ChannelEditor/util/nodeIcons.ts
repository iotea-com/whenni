'use client'

import isBrowser from '@iotea/libs/frontend/util/isBrowser'
import {
  riCodeBoxFill,
  riFilter3Fill,
  riQuestionFill,
  riShareFill,
  riPlayFill,
  riWebhookFill,
  riTimerLine,
  riFolderCloudLine,
  riDatabase2Line,
  riFileTextLine,
  riNotification4Line,
  riShapesLine,
  riFlowChart,
  riErrorWarningFill,
  riBarChartBoxLine,
} from '@mwarnerdotme/react-remixicon'

export const cloneImage = (img: HTMLImageElement) => {
  const clone = new Image()
  clone.src = img.src
  return clone
}

let questionIcon: HTMLImageElement | null = null
let sourceIcon: HTMLImageElement | null = null
let processingIcon: HTMLImageElement | null = null
let actionIcon: HTMLImageElement | null = null
let conditionalIcon: HTMLImageElement | null = null
let customIcon: HTMLImageElement | null = null
let messageQueueIcon: HTMLImageElement | null = null
let httpIcon: HTMLImageElement | null = null
let timerIcon: HTMLImageElement | null = null
let fileStorageIcon: HTMLImageElement | null = null
let databaseIcon: HTMLImageElement | null = null
let logIcon: HTMLImageElement | null = null
let notificationIcon: HTMLImageElement | null = null
let shapesIcon: HTMLImageElement | null = null
let flowChartIcon: HTMLImageElement | null = null
let errorIcon: HTMLImageElement | null = null
let metricIcon: HTMLImageElement | null = null
if (isBrowser) {
  questionIcon = new Image()
  questionIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="${riQuestionFill.pathData}" fill="currentColor"/>
    </svg>
  `)

  sourceIcon = new Image()
  sourceIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="${riPlayFill.pathData}" fill="currentColor"/>
    </svg>
  `)

  processingIcon = new Image()
  processingIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="${riFilter3Fill.pathData}" fill="currentColor"/>
    </svg>
  `)

  actionIcon = new Image()
  actionIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g transform="translate(16 16) rotate(180) translate(-8 -8)">
        <path d="${riPlayFill.pathData}" fill="currentColor"/>
      </g>
    </svg>
  `)

  conditionalIcon = new Image()
  conditionalIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="${riShareFill.pathData}" fill="currentColor"/>
    </svg>
  `)

  customIcon = new Image()
  customIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="${riCodeBoxFill.pathData}" fill="currentColor"/>
    </svg>
  `)

  messageQueueIcon = new Image()
  messageQueueIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 256 256" fill="none" xmlns="http://www.w3.org/2000/svg">
        <g fill="currentColor" transform="translate(0,-20)">
          <g transform="matrix(1,0,0,1,28,2)">
              <path d="M144.253,147.119L144.253,141.56C144.253,140.884 144.603,140.257 145.178,139.901C145.753,139.546 146.471,139.513 147.076,139.816L174.11,153.333C174.771,153.663 175.188,154.339 175.188,155.078C175.188,155.816 174.771,156.492 174.11,156.822L147.076,170.339C146.471,170.642 145.753,170.609 145.178,170.254C144.603,169.898 144.253,169.271 144.253,168.595L144.253,163.036L128,163.036C123.605,163.036 120.042,159.473 120.042,155.078C120.042,155.078 120.042,92.334 120.042,79.786L120.042,79.567L26.743,79.567C22.351,79.567 18.785,76.001 18.785,71.609C18.785,67.217 22.351,63.651 26.743,63.651L128,63.651C132.395,63.651 135.958,67.214 135.958,71.609C135.958,71.609 135.958,147.119 135.958,147.119L144.253,147.119Z"/>
          </g>
          <g transform="matrix(-1,0,0,-1,175,252)">
              <path d="M144.253,163.036L128,163.036C123.605,163.036 120.042,159.473 120.042,155.078C120.042,155.078 120.042,92.334 120.042,79.786L120.042,79.567L26.743,79.567C22.351,79.567 18.785,76.001 18.785,71.609C18.785,67.217 22.351,63.651 26.743,63.651L128,63.651C132.395,63.651 135.958,67.214 135.958,71.609C135.958,71.609 135.958,147.119 135.958,147.119L144.253,147.119L144.253,141.56C144.253,140.884 144.603,140.257 145.178,139.901C145.753,139.546 146.471,139.513 147.076,139.816L174.11,153.333C174.771,153.663 175.188,154.339 175.188,155.078C175.188,155.816 174.771,156.492 174.11,156.822L147.076,170.339C146.471,170.642 145.753,170.609 145.178,170.254C144.603,169.898 144.253,169.271 144.253,168.595L144.253,163.036Z"/>
          </g>
          <g transform="matrix(1.3,0,0,1.3,-90,-37)">
              <circle cx="124" cy="127" r="10"/>
          </g>
          <g transform="matrix(1.3,0,0,1.3,-60,-37)">
              <circle cx="124" cy="127" r="10"/>
          </g>
          <g transform="matrix(1.3,0,0,1.3,-30,-37)">
          <circle cx="124" cy="127" r="10"/>
          </g>
        </g>
    </svg>
  `)

  httpIcon = new Image()
  httpIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="${riWebhookFill.pathData}" fill="currentColor"/>
    </svg>
  `)

  timerIcon = new Image()
  timerIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="${riTimerLine.pathData}" fill="currentColor"/>
    </svg>
  `)

  fileStorageIcon = new Image()
  fileStorageIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="${riFolderCloudLine.pathData}" fill="currentColor"/>
    </svg>
  `)

  databaseIcon = new Image()
  databaseIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="${riDatabase2Line.pathData}" fill="currentColor"/>
    </svg>
  `)

  logIcon = new Image()
  logIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="${riFileTextLine.pathData}" fill="currentColor"/>
    </svg>
  `)

  metricIcon = new Image()
  metricIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="${riBarChartBoxLine.pathData}" fill="currentColor"/>
    </svg>
  `)

  notificationIcon = new Image()
  notificationIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="${riNotification4Line.pathData}" fill="currentColor"/>
    </svg>
  `)

  shapesIcon = new Image()
  shapesIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="${riShapesLine.pathData}" fill="currentColor"/>
    </svg>
  `)

  flowChartIcon = new Image()
  flowChartIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="${riFlowChart.pathData}" fill="currentColor"/>
    </svg>
  `)

  errorIcon = new Image()
  errorIcon.src =
    'data:image/svg+xml,' +
    encodeURIComponent(`
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="${riErrorWarningFill.pathData}" fill="currentColor"/>
    </svg>
  `)
}

const nodeIcons = {
  // Node types
  question: questionIcon,
  source: sourceIcon,
  processing: processingIcon,
  action: actionIcon,
  conditional: conditionalIcon,
  custom: customIcon,

  // Node labels
  messageQueue: messageQueueIcon,
  http: httpIcon,
  timer: timerIcon,
  fileStorage: fileStorageIcon,
  database: databaseIcon,
  log: logIcon,
  metric: metricIcon,
  notification: notificationIcon,
  shapes: shapesIcon,
  flowChart: flowChartIcon,

  // Other
  error: errorIcon,
}

export default nodeIcons
