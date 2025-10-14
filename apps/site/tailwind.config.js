const path = require('path')

const theme = require('../../libs/frontend/src/themes/tailwind.js')

/** @type {import('tailwindcss').Config} */
module.exports = {
  theme,
  content: [
    path.join(__dirname, './**/*.{js,ts,jsx,tsx}'),
    path.join(__dirname, '../../libs/**/*.{js,ts,jsx,tsx}'),
  ],
  darkMode: 'selector',
}
