"use client";

import { Alert, Badge, Button, Card, Group, NumberInput, Stack, Switch, Text, TextInput } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useEffect, useState } from "react";

import { apiFetch, ApiError } from "@/lib/api";

interface NetworkSettings {
  geoip_db_path: string;
  geoip_available: boolean;
  crowdsec_enabled: boolean;
  crowdsec_api_url: string;
  crowdsec_api_key_configured: boolean;
  crowdsec_poll_interval_secs: number;
  crowdsec_active: boolean;
  crowdsec_ip_decisions: number;
  crowdsec_range_decisions: number;
}

export function NetworkSecurityCard() {
  const [settings, setSettings] = useState<NetworkSettings | null>(null);
  const [geoipPath, setGeoipPath] = useState("");
  const [crowdsecEnabled, setCrowdsecEnabled] = useState(false);
  const [crowdsecUrl, setCrowdsecUrl] = useState("");
  const [crowdsecKey, setCrowdsecKey] = useState("");
  const [pollSecs, setPollSecs] = useState<number | "">(60);
  const [saving, setSaving] = useState(false);

  const load = async () => {
    try {
      const response = await apiFetch<NetworkSettings>("/api/v1/settings/network");
      setSettings(response);
      setGeoipPath(response.geoip_db_path || "");
      setCrowdsecEnabled(response.crowdsec_enabled);
      setCrowdsecUrl(response.crowdsec_api_url || "");
      setPollSecs(response.crowdsec_poll_interval_secs || 60);
    } catch (err) {
      notifications.show({ color: "danger", message: err instanceof ApiError ? err.message : "Failed to load network settings." });
    }
  };

  useEffect(() => {
    void load();
  }, []);

  const save = async () => {
    setSaving(true);
    try {
      const body: Record<string, unknown> = {
        geoip_db_path: geoipPath,
        crowdsec_enabled: crowdsecEnabled,
        crowdsec_api_url: crowdsecUrl,
        crowdsec_poll_interval_secs: typeof pollSecs === "number" ? pollSecs : 60
      };
      if (crowdsecKey.trim()) {
        body.crowdsec_api_key = crowdsecKey.trim();
      }
      const response = await apiFetch<NetworkSettings>("/api/v1/settings/network", {
        method: "PATCH",
        body: JSON.stringify(body)
      });
      setSettings(response);
      setCrowdsecKey("");
      notifications.show({ color: "success", message: "Network security settings saved and applied." });
    } catch (err) {
      notifications.show({ color: "danger", message: err instanceof ApiError ? err.message : "Save failed." });
    } finally {
      setSaving(false);
    }
  };

  return (
    <Stack gap="md">
      <Card withBorder>
        <Stack gap="md">
          <Group justify="space-between">
            <Text fw={600}>GeoIP</Text>
            <Badge color={settings?.geoip_available ? "success" : "gray"}>
              {settings?.geoip_available ? "database loaded" : "no database"}
            </Badge>
          </Group>
          <Text size="sm" c="dimmed">
            Path to a MaxMind GeoLite2-Country .mmdb file inside the container. Per-service country allow/block lists need this loaded.
          </Text>
          <TextInput
            label="GeoIP database path"
            placeholder="/data/geoip/GeoLite2-Country.mmdb"
            value={geoipPath}
            onChange={(e) => setGeoipPath(e.currentTarget.value)}
          />
        </Stack>
      </Card>

      <Card withBorder>
        <Stack gap="md">
          <Group justify="space-between">
            <Text fw={600}>CrowdSec</Text>
            <Group gap="xs">
              {settings?.crowdsec_active ? (
                <Badge color="success">active · {settings.crowdsec_ip_decisions} IPs / {settings.crowdsec_range_decisions} ranges</Badge>
              ) : (
                <Badge color="gray">inactive</Badge>
              )}
            </Group>
          </Group>
          <Text size="sm" c="dimmed">
            Pull IP block decisions from a CrowdSec LAPI. Blocked IPs are rejected before authentication.
          </Text>
          <Switch label="Enable CrowdSec" checked={crowdsecEnabled} onChange={(e) => setCrowdsecEnabled(e.currentTarget.checked)} />
          <TextInput
            label="LAPI URL"
            placeholder="http://crowdsec:8080"
            value={crowdsecUrl}
            onChange={(e) => setCrowdsecUrl(e.currentTarget.value)}
          />
          <TextInput
            label="Bouncer API key"
            description={settings?.crowdsec_api_key_configured ? "Leave blank to keep the existing key." : "Create with: cscli bouncers add portlyn"}
            value={crowdsecKey}
            onChange={(e) => setCrowdsecKey(e.currentTarget.value)}
            type="password"
          />
          <NumberInput
            label="Poll interval (seconds)"
            value={pollSecs}
            onChange={(v) => setPollSecs(typeof v === "number" ? v : 60)}
            min={10}
            max={3600}
          />
        </Stack>
      </Card>

      <Alert color="info" variant="light">
        Changes apply immediately. GeoIP reloads the database and CrowdSec restarts its poller without a server restart.
      </Alert>

      <Group justify="flex-end">
        <Button onClick={() => void save()} loading={saving}>Save network settings</Button>
      </Group>
    </Stack>
  );
}
