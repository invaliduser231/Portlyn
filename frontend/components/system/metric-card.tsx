"use client";

import { Card, Group, Stack, Text } from "@mantine/core";
import type { ReactNode } from "react";

export function MetricCard({
  label,
  value,
  hint,
  accent
}: {
  label: string;
  value: ReactNode;
  hint?: string;
  accent?: string;
}) {
  return (
    <Card withBorder radius="md" p="lg" style={{ background: "rgba(255,255,255,0.02)" }}>
      <Stack gap={6}>
        <Text size="xs" tt="uppercase" fw={700} c="#7e8795">
          {label}
        </Text>
        <Group align="baseline" gap="xs">
          <Text fw={700} fz={28} c={accent || "white"}>
            {value}
          </Text>
        </Group>
        {hint ? <Text size="sm" c="dimmed">{hint}</Text> : null}
      </Stack>
    </Card>
  );
}
