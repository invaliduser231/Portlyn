"use client";

import { Paper, Stack, Text, Title } from "@mantine/core";

export function EmptyState({ title, description }: { title: string; description: string }) {
  return (
    <Paper p="xl">
      <Stack align="center" gap={6} py={12}>
        <Title order={4} fw={600}>
          {title}
        </Title>
        <Text c="#7e8795" size="sm">
          {description}
        </Text>
      </Stack>
    </Paper>
  );
}
