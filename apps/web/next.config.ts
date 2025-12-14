import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  async rewrites() {
    return [
      {
        source: '/docs',
        destination: 'http://localhost:3001/docs',
      },
      {
        source: '/docs/:path*',
        destination: 'http://localhost:3001/docs/:path*',
      },
    ];
  },
};

export default nextConfig;
