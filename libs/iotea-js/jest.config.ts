/* eslint-disable */
export default {
  displayName: 'iotea-js',
  preset: '../../jest.preset.js',
  moduleFileExtensions: ['ts', 'js'],
  collectCoverageFrom: [
    // Include all .ts and .tsx files in the src directory, excluding test files
    'src/**/*.{ts,tsx}',

    // Exclude index files or anything that shouldn't be covered (optional)
    '!src/**/*.{d.ts}',
    '!src/**/index.ts',
  ],
}
