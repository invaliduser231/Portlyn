"use client";

import { AppShell, Box } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import type { ReactNode } from "react";

import { AuthGuard } from "@/components/layout/auth-guard";
import { Sidebar } from "@/components/layout/sidebar";
import { Topbar } from "@/components/layout/topbar";

export function DashboardShell({ children }: { children: ReactNode }) {
  const [opened, { toggle, close }] = useDisclosure(false);

  return (
    <AuthGuard>
      <AppShell
        header={{ height: 64 }}
        navbar={{ width: { base: 240, lg: 272 }, breakpoint: "sm", collapsed: { mobile: !opened } }}
        padding={{ base: "md", md: "xl" }}
      >
        <AppShell.Header style={{ background: "rgba(31,31,35,0.78)", backdropFilter: "blur(12px)" }}>
          <Topbar opened={opened} onToggle={toggle} />
        </AppShell.Header>

        <AppShell.Navbar
          p={0}
          style={{
            background: "linear-gradient(180deg, #1e1f25 0%, #16171c 100%)",
            borderRight: "1px solid rgba(255,255,255,0.04)"
          }}
        >
          <Sidebar onNavigate={close} />
        </AppShell.Navbar>

        <AppShell.Main>
          <Box maw={1180} mx="auto">
            {children}
          </Box>
        </AppShell.Main>
      </AppShell>
    </AuthGuard>
  );
}
