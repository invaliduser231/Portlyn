"use client";

import { Badge, Button, Code, Group, Stack, Text } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { IconBolt, IconLinkOff } from "@tabler/icons-react";
import { useState } from "react";

import { apiFetch, ApiError } from "@/lib/api";
import { formatDateTime } from "@/lib/format";
import type { Node, TunnelStatus } from "@/lib/types";

const STATUS_COLOR: Record<string, string> = {
  inactive: "gray",
  provisioned: "warning",
  connected: "success",
  stale: "warning",
};

export function NodeTunnelPanel({ node, onChange }: { node: Node; onChange: (updated: Node) => void }) {
  const [revoking, setRevoking] = useState(false);

  const status = (node.tunnel_status as TunnelStatus) || "inactive";
  const hasTunnel = (node.wg_public_key?.length ?? 0) > 0 && (node.wg_tunnel_ip?.length ?? 0) > 0;

  const revoke = async () => {
    setRevoking(true);
    try {
      await apiFetch<void>(`/api/v1/nodes/${node.id}/wg-tunnel`, { method: "DELETE" });
      onChange({
        ...node,
        wg_public_key: "",
        wg_tunnel_ip: "",
        wg_endpoint: "",
        tunnel_status: "inactive",
      });
      notifications.show({ color: "success", message: "Tunnel revoked." });
    } catch (err) {
      notifications.show({ color: "danger", message: err instanceof ApiError ? err.message : "Revoke failed." });
    } finally {
      setRevoking(false);
    }
  };

  return (
    <Stack gap="xs">
      <Group justify="space-between">
        <Group gap="xs">
          <IconBolt size={16} />
          <Text fw={600} size="sm">Tunnel</Text>
          <Badge color={STATUS_COLOR[status] || "gray"}>{status}</Badge>
        </Group>
        {hasTunnel ? (
          <Button size="xs" variant="subtle" color="danger" leftSection={<IconLinkOff size={14} />} loading={revoking} onClick={revoke}>
            Revoke
          </Button>
        ) : null}
      </Group>

      {hasTunnel ? (
        <Stack gap={4}>
          <Text size="xs" c="dimmed">Tunnel IP</Text>
          <Code>{node.wg_tunnel_ip}</Code>
          {node.advertised_subnets ? (
            <>
              <Text size="xs" c="dimmed">Advertised subnets</Text>
              <Code>{node.advertised_subnets}</Code>
            </>
          ) : null}
          <Text size="xs" c="dimmed">Last handshake</Text>
          <Text size="sm">{formatDateTime(node.wg_last_handshake)}</Text>
        </Stack>
      ) : (
        <Text size="sm" c="dimmed">
          The tunnel is established automatically when the node agent connects with an enrollment token.
        </Text>
      )}
    </Stack>
  );
}
