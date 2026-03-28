"use client";

import { Button, Group, Paper, Select, Skeleton, Stack, Text, TextInput } from "@mantine/core";
import { useEffect, useMemo, useState } from "react";

import { AdminOnly } from "@/components/admin-only";
import { AuditLogTable } from "@/components/audit/audit-log-table";
import { EmptyState } from "@/components/empty-state";
import { ErrorState } from "@/components/error-state";
import { apiFetch, ApiError } from "@/lib/api";
import type { AuditLogListResponse } from "@/lib/types";

export default function AuditLogsPage() {
  const [data, setData] = useState<AuditLogListResponse | null>(null);
  const [userId, setUserId] = useState("");
  const [resourceType, setResourceType] = useState("");
  const [fromDate, setFromDate] = useState("");
  const [toDate, setToDate] = useState("");
  const [query, setQuery] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const filteredItems = useMemo(() => {
    if (!data) return [];
    return data.items.filter((item) =>
      !query
        ? true
        : [item.action, item.resource_type, item.details, String(item.resource_id ?? ""), String(item.user_id ?? "")]
            .join(" ")
            .toLowerCase()
            .includes(query.toLowerCase())
    );
  }, [data, query]);

  const loadData = async () => {
    setIsLoading(true);
    setError(null);
    const params = new URLSearchParams({ limit: "50" });
    if (userId) params.set("user_id", userId);
    if (resourceType) params.set("resource_type", resourceType);
    if (fromDate) params.set("from", String(Math.floor(new Date(fromDate).getTime() / 1000)));
    if (toDate) params.set("to", String(Math.floor(new Date(toDate).getTime() / 1000)));

    try {
      const response = await apiFetch<AuditLogListResponse>(`/api/v1/audit-logs?${params.toString()}`);
      setData(response);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Unable to load audit logs.");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    void loadData();
  }, []);

  return (
    <AdminOnly>
      <Stack gap="lg">
        <Group grow>
          <TextInput placeholder="User ID" value={userId} onChange={(event) => setUserId(event.currentTarget.value)} />
          <Select
            data={[
              { value: "", label: "All resources" },
              { value: "service", label: "service" },
              { value: "domain", label: "domain" },
              { value: "node", label: "node" },
              { value: "certificate", label: "certificate" },
              { value: "user", label: "user" },
              { value: "auth", label: "auth" }
            ]}
            value={resourceType}
            onChange={(value) => setResourceType(value || "")}
          />
          <TextInput type="date" value={fromDate} onChange={(event) => setFromDate(event.currentTarget.value)} />
          <TextInput type="date" value={toDate} onChange={(event) => setToDate(event.currentTarget.value)} />
        </Group>

        <Group grow>
          <TextInput placeholder="Search" value={query} onChange={(event) => setQuery(event.currentTarget.value)} />
          <Button variant="default" onClick={() => void loadData()}>
            Apply Filters
          </Button>
        </Group>

        {error ? <ErrorState title="Failed to load audit logs" message={error} onRetry={() => void loadData()} /> : null}

        {isLoading ? (
          <Stack gap="sm"><Skeleton height={54} /><Skeleton height={54} /></Stack>
        ) : !data || filteredItems.length === 0 ? (
          <EmptyState title={!data || data.items.length === 0 ? "No audit events found" : "No matching audit events"} description={!data || data.items.length === 0 ? "Adjust the filters." : "Adjust the search."} />
        ) : (
          <Paper withBorder radius="md" p="sm">
            <Stack gap="sm">
              <Text c="dimmed" size="sm">{data.total} total events</Text>
              <AuditLogTable items={filteredItems} />
            </Stack>
          </Paper>
        )}
      </Stack>
    </AdminOnly>
  );
}
