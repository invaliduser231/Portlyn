"use client";

import { Table, Text } from "@mantine/core";

import { formatDateTime } from "@/lib/format";
import type { AuditLog } from "@/lib/types";

export function AuditLogTable({ items }: { items: AuditLog[] }) {
  return (
    <Table.ScrollContainer minWidth={980}>
      <Table striped highlightOnHover withTableBorder>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Timestamp</Table.Th>
            <Table.Th>User</Table.Th>
            <Table.Th>Action</Table.Th>
            <Table.Th>Resource Type</Table.Th>
            <Table.Th>Resource ID</Table.Th>
            <Table.Th>Details</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {items.map((item) => (
            <Table.Tr key={item.id}>
              <Table.Td>{formatDateTime(item.timestamp)}</Table.Td>
              <Table.Td>{item.user_id ? `#${item.user_id}` : "system"}</Table.Td>
              <Table.Td>{item.action.replace("_", " ")}</Table.Td>
              <Table.Td>{item.resource_type}</Table.Td>
              <Table.Td>{item.resource_id ? `#${item.resource_id}` : "-"}</Table.Td>
              <Table.Td><Text c="dimmed" size="sm">{item.details || "No details"}</Text></Table.Td>
            </Table.Tr>
          ))}
        </Table.Tbody>
      </Table>
    </Table.ScrollContainer>
  );
}
