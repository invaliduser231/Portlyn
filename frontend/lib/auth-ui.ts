import type { AuthUISettings } from "@/lib/types";

const legacyColorMap: Record<string, string> = {
  "#0a0d14": "#0d0e11",
  "#162033": "#121316",
  "#111826": "#181b20",
  "#2f6fed": "#654c96",
  "#f8fafc": "#f4f7fb",
  "#94a3b8": "#9aa3b2"
};

export const defaultAuthUI: AuthUISettings = {
  brand_name: "Portlyn",
  logo_url: "/logo.png",
  background_color: "#0d0e11",
  background_accent: "#1f232b",
  panel_color: "#181b20",
  button_color: "#654c96",
  text_color: "#f4f7fb",
  muted_text_color: "#9aa3b2",
  login_title: "Login",
  login_subtitle: "",
  route_login_title: "Login",
  route_login_subtitle: "",
  forbidden_title: "Access denied",
  forbidden_subtitle: "You do not have permission to access this route.",
  login_password_label: "Login",
  login_oidc_label: "Continue with SSO",
  login_otp_request_label: "Request code",
  login_otp_verify_label: "Verify code",
  route_continue_label: "Continue",
  route_oidc_label: "Continue with SSO",
  route_pin_label: "Unlock",
  route_email_send_label: "Send code",
  route_email_verify_label: "Verify code",
  forbidden_retry_label: "Try again"
};

export function mergeAuthUI(value?: Partial<AuthUISettings> | null): AuthUISettings {
  const merged = { ...defaultAuthUI, ...(value || {}) };

  return {
    ...merged,
    logo_url: merged.logo_url?.trim() ? merged.logo_url : defaultAuthUI.logo_url,
    background_color: normalizeLegacyColor(merged.background_color),
    background_accent: normalizeLegacyColor(merged.background_accent),
    panel_color: normalizeLegacyColor(merged.panel_color),
    button_color: normalizeLegacyColor(merged.button_color),
    text_color: normalizeLegacyColor(merged.text_color),
    muted_text_color: normalizeLegacyColor(merged.muted_text_color)
  };
}

export function authShellStyle(ui: AuthUISettings) {
  return {
    minHeight: "100vh",
    display: "grid",
    placeItems: "center",
    padding: 16,
    background: `linear-gradient(180deg, ${ui.background_color} 0%, #121316 48%, #0d0e11 100%)`
  } as const;
}

export function authCardStyle(ui: AuthUISettings) {
  return {
    backgroundColor: ui.panel_color,
    color: ui.text_color,
    borderColor: "rgba(255, 255, 255, 0.08)",
    boxShadow: "0 28px 64px -36px rgba(0, 0, 0, 0.92)"
  } as const;
}

export function buttonStyle(ui: AuthUISettings) {
  return {
    backgroundColor: ui.button_color,
    borderColor: "rgba(255, 255, 255, 0.08)",
    color: ui.text_color,
    boxShadow: "inset 0 1px 0 rgba(255,255,255,0.04)"
  } as const;
}

export function inputStyles(ui: AuthUISettings) {
  return {
    label: {
      color: ui.text_color
    },
    input: {
      backgroundColor: "rgba(255, 255, 255, 0.03)",
      color: ui.text_color,
      borderColor: "rgba(255, 255, 255, 0.1)"
    }
  } as const;
}

export function authInfoAlertStyle(ui: AuthUISettings) {
  return {
    root: {
      backgroundColor: "rgba(255, 255, 255, 0.04)",
      borderColor: "rgba(255, 255, 255, 0.08)",
      color: ui.text_color
    },
    title: {
      color: ui.text_color
    },
    body: {
      color: ui.text_color
    }
  } as const;
}

function normalizeLegacyColor(value: string) {
  return legacyColorMap[value.toLowerCase()] || value;
}
