"use client";

import { Box, Container } from "@mantine/core";
import type { ReactNode } from "react";

import { AuthGuard } from "@/components/layout/auth-guard";
import { Sidebar } from "@/components/layout/sidebar";
import { Topbar } from "@/components/layout/topbar";

const SIDEBAR_WIDTH = 284;
const TOPBAR_HEIGHT = 76;

export function DashboardShell({ children }: { children: ReactNode }) {
  return (
    <AuthGuard>
      <Box mih="100vh">
        <Box
          pos="fixed"
          top={0}
          left={0}
          bottom={0}
          w={SIDEBAR_WIDTH}
          style={{ zIndex: 30 }}
        >
          <Sidebar />
        </Box>

        <Box ml={SIDEBAR_WIDTH}>
          <Box
            pos="fixed"
            top={0}
            left={SIDEBAR_WIDTH}
            right={0}
            h={TOPBAR_HEIGHT}
            style={{ zIndex: 20 }}
          >
            <Topbar />
          </Box>

          <Box pt={TOPBAR_HEIGHT}>
            <Container size="xl" px={32} py={36}>
              <Box maw={1180} mx="auto">
                {children}
              </Box>
            </Container>
          </Box>
        </Box>
      </Box>
    </AuthGuard>
  );
}
