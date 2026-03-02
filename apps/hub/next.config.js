//@ts-check

const { composePlugins, withNx } = require('@nx/next')

const nextConfig = {
  nx: {
    // Set this to true if you would like to use SVGR
    // See: https://github.com/gregberge/svgr
    svgr: false,
  },
  transpilePackages: ['@iotea/libs/frontend'],
}

const plugins = [withNx]

module.exports = composePlugins(...plugins)(nextConfig)
