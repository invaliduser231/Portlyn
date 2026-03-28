import type { Metadata } from "next";
import type { ReactNode } from "react";

import "@mantine/core/styles.css";
import "@mantine/notifications/styles.css";
import "@/app/globals.css";
import { Providers } from "@/components/providers";

export const metadata: Metadata = {
  title: "Portlyn",
  description: "Portlyn control plane",
  icons: {
    icon: "/favicon.png",
    apple: "/apple-icon.png"
  }
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="de">
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
