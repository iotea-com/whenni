const nxPreset = require('@nx/jest/preset').default

module.exports = {
  ...nxPreset,
  transform: {
    '^.+\\.ts$': ['ts-jest', { 
      tsconfig: '<rootDir>/tsconfig.test.json',
    }],
  },
}
