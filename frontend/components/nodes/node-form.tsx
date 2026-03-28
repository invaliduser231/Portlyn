"use client";

import { Button, Select, Stack, TextInput, Textarea } from "@mantine/core";
import { useEffect, useState } from "react";

import type { Node, NodePayload } from "@/lib/types";

const defaults: NodePayload = {
  name: "",
  description: "",
  status: "unknown",
  version: ""
};

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
    version: initialValues?.version || defaults.version
  });
  const [values, setValues] = useState<NodePayload>(getInitialState);

  useEffect(() => {
    setValues(getInitialState());
  }, [initialValues]);

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
      <Button loading={isLoading} onClick={() => void onSubmit(values)} disabled={!values.name}>
        {submitLabel}
      </Button>
    </Stack>
  );
}
