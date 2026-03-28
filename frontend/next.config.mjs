/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  async rewrites() {
    const target = process.env.INTERNAL_API_BASE_URL?.replace(/\/$/, "") || "http://localhost:8080";
    return [
      {
        source: "/api/:path*",
        destination: `${target}/api/:path*`
      }
    ];
  }
};

export default nextConfig;
