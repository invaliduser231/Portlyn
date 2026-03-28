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
      styles={{ root: { background: "rgba(255,255,255,0.03)" }, title: { color: "#f4f7fb" } }}
    >
      <Stack gap="sm">
        <Text size="sm" c="#aab3c2">{message}</Text>
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
