"use client";

import { Alert, Button, Group, Stack, Text } from "@mantine/core";
import { IconAlertCircle } from "@tabler/icons-react";

export function ErrorState({
  title = "Something went wrong",
  message,
  onRetry
}: {
  title?: string;
  message: string;
  onRetry?: () => void;
}) {
  return (
    <Alert
      icon={<IconAlertCircle size={16} />}
      title={title}
      color="gray"
      variant="light"
      styles={{ root: { background: "var(--portlyn-surface-raised)" }, title: { color: "var(--portlyn-text)" } }}
    >
      <Stack gap="sm">
        <Text size="sm" c="dimmed">{message}</Text>
        {onRetry ? (
          <Group>
            <Button size="xs" variant="subtle" color="gray" onClick={onRetry}>
              Retry
            </Button>
          </Group>
        ) : null}
      </Stack>
    </Alert>
  );
}
