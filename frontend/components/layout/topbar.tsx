"use client";

import { Button, Group, Paper, Stack, Text } from "@mantine/core";
import { IconLogout2 } from "@tabler/icons-react";
import { usePathname } from "next/navigation";

import { useAuth, UserAvatar } from "@/components/providers";

const titles: Record<string, string> = {
  "/services": "Services",
  "/domains": "Domains",
  "/nodes": "Nodes",
  "/certificates": "Certificates",
  "/system": "System",
  "/groups": "Groups",
  "/service-groups": "Service Groups",
  "/users": "Users",
  "/audit-logs": "Audit Logs",
  "/settings": "Settings"
};

export function Topbar() {
  const pathname = usePathname();
  const { user, logout } = useAuth();
  const baseTitle =
    Object.entries(titles).find(([route]) => pathname === route || pathname.startsWith(`${route}/`))?.[1] ||
    "Portlyn";
  const title = user?.role === "viewer" && baseTitle === "Services" ? "Applications" : baseTitle;

  return (
    <Paper
      radius={0}
      h="100%"
      px={32}
      style={{
        display: "flex",
        alignItems: "center",
        justifyContent: "space-between",
        borderBottom: "1px solid rgba(255,255,255,0.04)",
        background: "rgba(31,31,35,0.78)",
        backdropFilter: "blur(12px)"
      }}
    >
      <Text fw={600} fz={30} style={{ letterSpacing: "-0.03em" }}>
        {title}
      </Text>

      <Group gap="sm">
        <Paper p="xs" bg="rgba(255,255,255,0.03)">
          <Group gap="sm">
            <UserAvatar />
            <Stack gap={0}>
              <Text size="sm" fw={600}>
                {user?.email}
              </Text>
              <Text size="xs" c="#7e8795" tt="capitalize">
                {user?.role}
              </Text>
            </Stack>
          </Group>
        </Paper>
        <Button variant="subtle" color="gray" leftSection={<IconLogout2 size={16} />} onClick={logout}>
          Logout
        </Button>
      </Group>
    </Paper>
  );
}
