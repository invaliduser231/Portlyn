import { describe, expect, it } from "vitest";

import {
  accessMethodLabel,
  buildServiceRequestPayload,
  defaultServicePayload,
  legacyAuthPolicyFromAccessMode,
  linesToList
} from "@/lib/access-control";
import type { ServicePayload } from "@/lib/types";

describe("access-control helpers", () => {
  it("maps access modes to legacy auth policy", () => {
    expect(
      legacyAuthPolicyFromAccessMode({
        access_mode: "public",
        allowed_roles: [],
        allowed_groups: [],
        allowed_service_groups: []
      })
    ).toBe("public");

    expect(
      legacyAuthPolicyFromAccessMode({
        access_mode: "restricted",
        allowed_roles: ["admin"],
        allowed_groups: [],
        allowed_service_groups: []
      })
    ).toBe("admin_only");

    expect(
      legacyAuthPolicyFromAccessMode({
        access_mode: "restricted",
        allowed_roles: ["viewer"],
        allowed_groups: [],
        allowed_service_groups: []
      })
    ).toBe("authenticated");
  });

  it("normalizes line/csv strings to distinct non-empty entries", () => {
    expect(linesToList("one,\ntwo\n\nthree")).toEqual(["one", "two", "three"]);
  });

  it("sanitizes service request payload for pin mode", () => {
    const base: ServicePayload = defaultServicePayload();
    const payload = buildServiceRequestPayload(
      {
        ...base,
        access_method: "pin",
        access_method_config: {
          pin: "1234",
          hint: "Front door"
        }
      },
      { omitEmptyAccessMethod: false }
    );
    expect(payload.access_method_config).toEqual({
      pin: "1234",
      hint: "Front door"
    });
  });

  it("omits access method fields when empty and option enabled", () => {
    const payload = buildServiceRequestPayload(defaultServicePayload(), { omitEmptyAccessMethod: true }) as Record<string, unknown>;
    expect(payload.access_method).toBeUndefined();
    expect(payload.access_method_config).toBeUndefined();
  });

  it("returns readable labels", () => {
    expect(accessMethodLabel("session")).toBe("Session-based");
    expect(accessMethodLabel("oidc_only")).toBe("SSO required");
    expect(accessMethodLabel("pin")).toBe("PIN protected");
    expect(accessMethodLabel("email_code")).toBe("Email code required");
  });
});
