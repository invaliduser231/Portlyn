"use client";

import { Alert, MultiSelect, Select, Stack, Text } from "@mantine/core";

import type { AccessPolicy, Group as UserGroup, ServiceGroup, UserRole } from "@/lib/types";

export function AccessPolicyFields({
  value,
  groups,
  serviceGroups,
  onChange,
  inheritedFrom
}: {
  value: AccessPolicy;
  groups: UserGroup[];
  serviceGroups?: ServiceGroup[];
  onChange: (next: AccessPolicy) => void;
  inheritedFrom?: string | null;
}) {
  const restricted = value.access_mode === "restricted";

  return (
    <Stack gap="sm">
      {inheritedFrom ? (
        <Alert color="brand" variant="light">
          Inherits its default policy from {inheritedFrom}.
        </Alert>
      ) : null}

      <Select
        label="Access mode"
        data={[
          { value: "public", label: "Public" },
          { value: "authenticated", label: "Authenticated" },
          { value: "restricted", label: "Restricted" }
        ]}
        value={value.access_mode}
        onChange={(next) =>
          onChange({
            ...value,
            access_mode: (next || "authenticated") as AccessPolicy["access_mode"]
          })
        }
      />

      <Text c="dimmed" size="sm">
        `Restricted` grants access when the signed-in user matches at least one selected role or group.
      </Text>

      <MultiSelect
        label="Allowed roles"
        data={[
          { value: "admin", label: "admin" },
          { value: "viewer", label: "viewer" }
        ]}
        value={value.allowed_roles}
        onChange={(next) => onChange({ ...value, allowed_roles: next as UserRole[] })}
        disabled={!restricted}
        clearable
      />

      <MultiSelect
        label="Allowed user groups"
        data={groups.map((group) => ({ value: String(group.id), label: group.name }))}
        value={value.allowed_groups.map(String)}
        onChange={(next) => onChange({ ...value, allowed_groups: next.map(Number) })}
        disabled={!restricted}
        clearable
        searchable
      />

      {serviceGroups ? (
        <MultiSelect
          label="Policy service-groups"
          description="Used to document and reuse common policies across service groups."
          data={serviceGroups.map((group) => ({ value: String(group.id), label: group.name }))}
          value={value.allowed_service_groups.map(String)}
          onChange={(next) => onChange({ ...value, allowed_service_groups: next.map(Number) })}
          clearable
          searchable
        />
      ) : null}
    </Stack>
  );
}
