"use client";

import { Badge } from "@mantine/core";

import { accessMethodLabel } from "@/lib/access-control";
import type { AccessMethod, AccessMode, AuthPolicy } from "@/lib/types";

export function StatusBadge({ status }: { status: string }) {
  const normalized = status.toLowerCase();
  const color =
    normalized === "online" || normalized === "healthy" || normalized === "active" || normalized === "ok"
      ? "teal"
      : normalized === "unhealthy" || normalized === "error"
        ? "red"
        : normalized === "offline" || normalized === "inactive"
          ? "gray"
      : normalized === "warning" || normalized === "warn" || normalized === "pending" || normalized === "degraded"
          ? "orange"
          : "gray";

  return (
    <Badge color={color} variant="light" radius="xl" styles={{ root: { opacity: 0.9 } }}>
      {status.replace("_", " ")}
    </Badge>
  );
}

export function AuthPolicyBadge({ value }: { value: AuthPolicy }) {
  const color = value === "public" ? "brand" : value === "authenticated" ? "orange" : "gray";
  return (
    <Badge color={color} variant="light" radius="xl" styles={{ root: { opacity: 0.9 } }}>
      {value}
    </Badge>
  );
}

export function AccessModeBadge({ value }: { value: AccessMode }) {
  const color = value === "public" ? "green" : value === "authenticated" ? "orange" : "grape";
  return (
    <Badge color={color} variant="light" radius="xl" styles={{ root: { opacity: 0.9 } }}>
      {value}
    </Badge>
  );
}

export function AccessMethodBadge({ value }: { value: AccessMethod | undefined }) {
  const normalized = value || "session";
  const color =
    normalized === "oidc_only" ? "brand" : normalized === "pin" ? "yellow" : normalized === "email_code" ? "cyan" : "gray";
  return (
    <Badge color={color} variant="light" radius="xl" styles={{ root: { opacity: 0.9 } }}>
      {accessMethodLabel(normalized)}
    </Badge>
  );
}

export function RiskBadge({ value }: { value: string | undefined }) {
  const normalized = (value || "low").toLowerCase();
  const color = normalized === "high" ? "red" : normalized === "medium" ? "yellow" : "teal";
  return (
    <Badge color={color} variant="light" radius="xl" styles={{ root: { opacity: 0.9 } }}>
      {normalized}
    </Badge>
  );
}
