const { withNx } = require('@nx/next/plugins/with-nx')

/**
 * @type {import('@nx/next/plugins/with-nx').WithNxOptions}
 **/
const nextConfig = {
  nx: {
    // Set this to true if you would like to use SVGR
    // See: https://github.com/gregberge/svgr
    svgr: false,
  },
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'prod-files-secure.s3.us-west-2.amazonaws.com',
        pathname: '/**',
      },
      {
        protocol: 'https',
        hostname: 's3-us-west-2.amazonaws.com',
        pathname: '/**',
      },
      {
        protocol: 'https',
        hostname: 'www.notion.so',
        pathname: '/**',
      },
    ],
  },
}

module.exports = withNx(nextConfig)
