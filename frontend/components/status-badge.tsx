"use client";

import { Badge, Tooltip } from "@mantine/core";

import { accessMethodLabel } from "@/lib/access-control";
import type { AccessMethod, AccessMode, AuthPolicy } from "@/lib/types";

export function StatusBadge({ status }: { status: string }) {
  const normalized = status.toLowerCase();
  let color: string;
  if (normalized === "online" || normalized === "healthy" || normalized === "active" || normalized === "ok" || normalized === "mfa") {
    color = "success";
  } else if (normalized === "unhealthy" || normalized === "error" || normalized === "revoked") {
    color = "danger";
  } else if (normalized === "offline" || normalized === "inactive") {
    color = "gray";
  } else if (normalized === "warning" || normalized === "warn" || normalized === "pending" || normalized === "degraded") {
    color = "warning";
  } else {
    color = "gray";
  }
  return (
    <Badge color={color}>
      {status.replace("_", " ")}
    </Badge>
  );
}

function accessColor(value: string): string {
  if (value === "public") return "info";
  if (value === "authenticated") return "brand";
  return "warning";
}

export function AuthPolicyBadge({ value }: { value: AuthPolicy }) {
  return <Badge color={accessColor(value)}>{value}</Badge>;
}

export function AccessModeBadge({ value }: { value: AccessMode }) {
  return <Badge color={accessColor(value)}>{value}</Badge>;
}

export function AccessMethodBadge({ value }: { value: AccessMethod | undefined }) {
  const normalized = value || "session";
  const color =
    normalized === "oidc_only" ? "brand" : normalized === "pin" ? "warning" : normalized === "email_code" ? "accent" : "gray";
  return <Badge color={color}>{accessMethodLabel(normalized)}</Badge>;
}

export function RiskBadge({ value }: { value: string | undefined }) {
  const normalized = (value || "low").toLowerCase();
  const color = normalized === "high" ? "danger" : normalized === "medium" ? "warning" : "success";
  return (
    <Tooltip label="Exposure risk based on access mode, method and network rules" withArrow>
      <Badge color={color}>{normalized} risk</Badge>
    </Tooltip>
  );
}
