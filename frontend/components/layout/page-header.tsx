"use client";

import { Box, Group, Stack, Text } from "@mantine/core";
import type { ReactNode } from "react";

export function PageHeader({
  description,
  action,
  children
}: {
  description?: string;
  action?: ReactNode;
  children?: ReactNode;
}) {
  const hasTopRow = Boolean(description) || Boolean(action);
  return (
    <Stack gap="md" mb="lg">
      {hasTopRow ? (
        <Group justify="space-between" align="flex-end" wrap="wrap" gap="sm">
          {description ? (
            <Text c="dimmed" size="sm" maw={640} style={{ flex: 1, minWidth: 220 }}>
              {description}
            </Text>
          ) : (
            <span />
          )}
          {action ? <Group gap="sm">{action}</Group> : null}
        </Group>
      ) : null}
      {children ? <Box>{children}</Box> : null}
    </Stack>
  );
}
