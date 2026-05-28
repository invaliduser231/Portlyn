"use client";

import { Button, Card, Group, PasswordInput, Stack, Text } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { IconLock } from "@tabler/icons-react";
import { useState } from "react";

import { ApiError } from "@/lib/api";
import { changeOwnPassword } from "@/lib/auth";

export function PasswordCard() {
  const [current, setCurrent] = useState("");
  const [next, setNext] = useState("");
  const [confirm, setConfirm] = useState("");
  const [busy, setBusy] = useState(false);

  const submit = async () => {
    if (next !== confirm) {
      notifications.show({ color: "danger", message: "New passwords do not match." });
      return;
    }
    setBusy(true);
    try {
      await changeOwnPassword(current, next);
      setCurrent("");
      setNext("");
      setConfirm("");
      notifications.show({ color: "success", message: "Password changed." });
    } catch (err) {
      notifications.show({ color: "danger", message: err instanceof ApiError ? err.message : "Could not change password." });
    } finally {
      setBusy(false);
    }
  };

  return (
    <Card withBorder>
      <Stack gap="md">
        <Group gap="xs">
          <IconLock size={20} />
          <Text fw={600}>Password</Text>
        </Group>
        <PasswordInput label="Current password" value={current} onChange={(e) => setCurrent(e.currentTarget.value)} />
        <PasswordInput label="New password" description="At least 8 characters." value={next} onChange={(e) => setNext(e.currentTarget.value)} />
        <PasswordInput label="Confirm new password" value={confirm} onChange={(e) => setConfirm(e.currentTarget.value)} />
        <Group justify="flex-end">
          <Button onClick={() => void submit()} loading={busy} disabled={!current || next.length < 8 || !confirm}>
            Change password
          </Button>
        </Group>
      </Stack>
    </Card>
  );
}
