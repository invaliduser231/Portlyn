export function formatDateTime(value?: string | null) {
  if (!value) {
    return "Never";
  }

  return new Intl.DateTimeFormat("de-DE", {
    dateStyle: "medium",
    timeStyle: "short"
  }).format(new Date(value));
}

export function formatNumber(value: number, digits = 2) {
  return new Intl.NumberFormat("de-DE", {
    maximumFractionDigits: digits
  }).format(value);
}
