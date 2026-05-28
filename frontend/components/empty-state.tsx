"use client";

import { Paper, Stack, Text, ThemeIcon, Title } from "@mantine/core";
import { IconInbox } from "@tabler/icons-react";

export function EmptyState({ title, description }: { title: string; description: string }) {
  return (
    <Paper p="xl">
      <Stack align="center" gap={10} py={20}>
        <ThemeIcon size={44} radius="md" variant="light" color="gray">
          <IconInbox size={22} stroke={1.6} />
        </ThemeIcon>
        <Title order={4} fw={600}>
          {title}
        </Title>
        <Text c="dimmed" size="sm" ta="center" maw={360}>
          {description}
        </Text>
      </Stack>
    </Paper>
  );
}
