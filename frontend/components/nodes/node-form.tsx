"use client";

import { Button, Select, Stack, TagsInput, TextInput, Textarea } from "@mantine/core";
import { useEffect, useState } from "react";

import type { Node, NodePayload } from "@/lib/types";

const defaults: NodePayload = {
  name: "",
  description: "",
  status: "unknown",
  version: "",
  advertised_subnets: ""
};

function toSubnetList(value?: string): string[] {
  return (value || "")
    .split(",")
    .map((part) => part.trim())
    .filter(Boolean);
}

export function NodeForm({
  initialValues,
  onSubmit,
  submitLabel,
  isLoading
}: {
  initialValues?: Partial<Node>;
  onSubmit: (values: NodePayload) => Promise<void>;
  submitLabel: string;
  isLoading?: boolean;
}) {
  const getInitialState = (): NodePayload => ({
    name: initialValues?.name || defaults.name,
    description: initialValues?.description || defaults.description,
    status: initialValues?.status || defaults.status,
    version: initialValues?.version || defaults.version,
    advertised_subnets: initialValues?.advertised_subnets || defaults.advertised_subnets
  });
  const [values, setValues] = useState<NodePayload>(getInitialState);
  const [subnets, setSubnets] = useState<string[]>(toSubnetList(initialValues?.advertised_subnets));

  useEffect(() => {
    setValues(getInitialState());
    setSubnets(toSubnetList(initialValues?.advertised_subnets));
  }, [initialValues]);

  const submit = () => {
    void onSubmit({ ...values, advertised_subnets: subnets.join(",") });
  };

  return (
    <Stack gap="md">
      <TextInput label="Name" value={values.name} onChange={(event) => setValues({ ...values, name: event.currentTarget.value })} />
      <Textarea label="Description" value={values.description} onChange={(event) => setValues({ ...values, description: event.currentTarget.value })} minRows={4} />
      <Select
        label="Status"
        data={["unknown", "online", "offline"]}
        value={values.status}
        onChange={(value) => setValues({ ...values, status: value || "unknown" })}
      />
      <TextInput label="Version" value={values.version} onChange={(event) => setValues({ ...values, version: event.currentTarget.value })} />
      <TagsInput
        label="Advertised LAN subnets"
        description="CIDRs of the local networks this node exposes to the mesh, e.g. 192.168.1.0/24. Press Enter to add."
        placeholder="192.168.1.0/24"
        value={subnets}
        onChange={setSubnets}
      />
      <Button loading={isLoading} onClick={submit} disabled={!values.name}>
        {submitLabel}
      </Button>
    </Stack>
  );
}
