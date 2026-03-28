"use client";

import { ActionIcon, Group, Table, Text } from "@mantine/core";
import { IconEdit, IconPower, IconTrash } from "@tabler/icons-react";
import Link from "next/link";

import { StatusBadge } from "@/components/status-badge";
import { formatDateTime } from "@/lib/format";
import type { User } from "@/lib/types";

export function UserTable({
  users,
  currentUserId,
  onEdit,
  onToggleActive,
  onDelete
}: {
  users: User[];
  currentUserId?: number;
  onEdit: (user: User) => void;
  onToggleActive: (user: User) => void;
  onDelete: (user: User) => void;
}) {
  return (
    <Table.ScrollContainer minWidth={900}>
      <Table verticalSpacing="md" horizontalSpacing="md">
        <Table.Thead>
            <Table.Tr>
              <Table.Th>Email</Table.Th>
              <Table.Th>Identity</Table.Th>
              <Table.Th>Role</Table.Th>
              <Table.Th>Status</Table.Th>
              <Table.Th>Last Login</Table.Th>
            <Table.Th>Created</Table.Th>
            <Table.Th ta="right">Actions</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {users.map((user) => (
            <Table.Tr key={user.id}>
              <Table.Td>
                <Text component={Link} href={`/users/${user.id}`} fw={500}>{user.email}</Text>
                {user.id === currentUserId ? <Text c="#7e8795" size="xs">Current session</Text> : null}
              </Table.Td>
              <Table.Td>
                <Text tt="capitalize">{user.auth_provider}</Text>
                {user.auth_provider === "oidc" ? (
                  <Text c="#7e8795" size="xs">{user.auth_issuer || "SSO provider"}</Text>
                ) : (
                  <Text c="#7e8795" size="xs">Local password</Text>
                )}
              </Table.Td>
              <Table.Td>{user.role}</Table.Td>
              <Table.Td>{user.active ? <StatusBadge status="active" /> : <StatusBadge status="inactive" />}</Table.Td>
              <Table.Td>{formatDateTime(user.last_login_at)}</Table.Td>
              <Table.Td>{formatDateTime(user.created_at)}</Table.Td>
              <Table.Td>
                <Group justify="flex-end" gap="xs">
                  <ActionIcon variant="subtle" color="gray" onClick={() => onEdit(user)}>
                    <IconEdit size={16} />
                  </ActionIcon>
                  <ActionIcon variant="subtle" color="gray" onClick={() => onToggleActive(user)}>
                    <IconPower size={16} />
                  </ActionIcon>
                  <ActionIcon variant="subtle" color="gray" onClick={() => onDelete(user)}>
                    <IconTrash size={16} />
                  </ActionIcon>
                </Group>
              </Table.Td>
            </Table.Tr>
          ))}
        </Table.Tbody>
      </Table>
    </Table.ScrollContainer>
  );
}
