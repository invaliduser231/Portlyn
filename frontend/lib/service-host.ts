import type { Service } from "@/lib/types";

export function buildServiceHostname(rootDomain?: string | null, subdomain?: string | null) {
  const root = (rootDomain || "").trim().toLowerCase();
  const prefix = (subdomain || "").trim().toLowerCase();
  if (!root) {
    return "";
  }
  if (!prefix) {
    return root;
  }
  return `${prefix}.${root}`;
}

export function serviceHostname(service?: Partial<Service> | null) {
  if (!service) {
    return "";
  }
  return service.hostname || buildServiceHostname(service.domain?.name, service.subdomain);
}
