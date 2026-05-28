"use client";

import { Alert, Badge, Button, Card, Code, Group, NumberInput, Stack, Switch, Text, TextInput } from "@mantine/core";
import { notifications } from "@mantine/notifications";
import { useEffect, useState } from "react";

import { apiFetch, ApiError } from "@/lib/api";
import type { TunnelSettings } from "@/lib/types";

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
      <TunnelServerCard />

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

function TunnelServerCard() {
  const [settings, setSettings] = useState<TunnelSettings | null>(null);
  const [enabled, setEnabled] = useState(false);
  const [endpoint, setEndpoint] = useState("");
  const [listenPort, setListenPort] = useState<number | "">(51820);
  const [cidr, setCIDR] = useState("10.42.0.0/16");
  const [serverTunnelIP, setServerTunnelIP] = useState("10.42.0.1");
  const [saving, setSaving] = useState(false);

  const load = async () => {
    try {
      const response = await apiFetch<TunnelSettings>("/api/v1/tunnel/settings");
      setSettings(response);
      setEnabled(response.enabled);
      setEndpoint(response.server_endpoint || "");
      setListenPort(response.listen_port ?? 51820);
      setCIDR(response.cidr || "10.42.0.0/16");
      setServerTunnelIP(response.server_tunnel_ip || "10.42.0.1");
    } catch (err) {
      notifications.show({ color: "danger", message: err instanceof ApiError ? err.message : "Failed to load tunnel settings." });
    }
  };

  useEffect(() => {
    void load();
  }, []);

  const save = async () => {
    setSaving(true);
    try {
      const response = await apiFetch<TunnelSettings>("/api/v1/tunnel/settings", {
        method: "PATCH",
        body: JSON.stringify({
          enabled,
          server_endpoint: endpoint,
          listen_port: typeof listenPort === "number" ? listenPort : undefined,
          cidr,
          server_tunnel_ip: serverTunnelIP
        })
      });
      setSettings(response);
      notifications.show({ color: "success", message: "Tunnel settings saved." });
    } catch (err) {
      notifications.show({ color: "danger", message: err instanceof ApiError ? err.message : "Save failed." });
    } finally {
      setSaving(false);
    }
  };

  return (
    <Card withBorder>
      <Stack gap="md">
        <Group justify="space-between">
          <Text fw={600}>WireGuard tunnel server</Text>
          {settings?.configured ? (
            <Badge color="success">{settings.configured_peer_count} peers</Badge>
          ) : (
            <Badge color="gray">not configured</Badge>
          )}
        </Group>
        <Text size="sm" c="dimmed">
          Lets nodes behind NAT or CGNAT reach this server. Set a public endpoint, then run the node agent; the server keypair is generated automatically.
        </Text>

        {enabled && settings && !settings.configured ? (
          <Alert color="warning" variant="light">
            Set a public endpoint (host:port) so nodes can reach this server. The server keypair is generated on save.
          </Alert>
        ) : null}

        <Switch label="Enable tunnel" checked={enabled} onChange={(e) => setEnabled(e.currentTarget.checked)} />
        <TextInput
          label="Public endpoint"
          description="Host:port nodes will dial. Example: vpn.example.com:51820"
          value={endpoint}
          onChange={(e) => setEndpoint(e.currentTarget.value)}
        />
        <NumberInput
          label="Listen port"
          description="UDP port the tunnel listens on."
          value={listenPort}
          onChange={(v) => setListenPort(typeof v === "number" ? v : 51820)}
          min={1}
          max={65535}
        />
        <TextInput
          label="Tunnel CIDR"
          description="IP range allocated to peers (server uses .1)."
          value={cidr}
          onChange={(e) => setCIDR(e.currentTarget.value)}
        />
        <TextInput
          label="Server tunnel IP"
          value={serverTunnelIP}
          onChange={(e) => setServerTunnelIP(e.currentTarget.value)}
        />

        {settings?.configured ? (
          <Stack gap={4}>
            <Text size="sm" fw={500}>Server public key</Text>
            <Code block>{settings.server_public_key}</Code>
            <Text size="sm" c="dimmed">
              {settings.connected_peer_count} of {settings.configured_peer_count} peers reported a recent handshake.
            </Text>
          </Stack>
        ) : null}

        <Group justify="flex-end">
          <Button onClick={() => void save()} loading={saving}>Save tunnel settings</Button>
        </Group>
      </Stack>
    </Card>
  );
}
