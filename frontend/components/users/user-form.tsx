"use client";

import { Button, Checkbox, Select, Stack, TextInput } from "@mantine/core";
import { useEffect, useState } from "react";

import type { User, UserPayload, UserRole } from "@/lib/types";

const defaults: UserPayload = {
  email: "",
  password: "",
  role: "viewer",
  active: true
};

export function UserForm({
  initialValues,
  onSubmit,
  submitLabel,
  isLoading,
  requirePassword
}: {
  initialValues?: Partial<User>;
  onSubmit: (values: UserPayload) => Promise<void>;
  submitLabel: string;
  isLoading?: boolean;
  requirePassword?: boolean;
}) {
  const getInitialState = (): UserPayload => ({
    email: initialValues?.email || defaults.email,
    password: "",
    role: (initialValues?.role as UserRole) || defaults.role,
    active: initialValues?.active ?? defaults.active
  });

  const [values, setValues] = useState<UserPayload>(getInitialState);

  useEffect(() => {
    setValues(getInitialState());
  }, [initialValues]);

  return (
    <Stack gap="md">
      <TextInput label="Email" value={values.email} onChange={(event) => setValues({ ...values, email: event.currentTarget.value })} />
      <TextInput label={requirePassword ? "Password" : "Password (optional)"} type="password" value={values.password || ""} onChange={(event) => setValues({ ...values, password: event.currentTarget.value })} />
      <Select
        label="Role"
        data={[
          { value: "viewer", label: "viewer" },
          { value: "admin", label: "admin" }
        ]}
        value={values.role}
        onChange={(value) => setValues({ ...values, role: (value || "viewer") as UserRole })}
      />
      <Checkbox checked={values.active} onChange={(event) => setValues({ ...values, active: event.currentTarget.checked })} label="Active account" />
      <Button loading={isLoading} onClick={() => void onSubmit(values)} disabled={!values.email || (requirePassword && !values.password)}>
        {submitLabel}
      </Button>
    </Stack>
  );
}
