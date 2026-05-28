import { createTheme, defaultVariantColorsResolver, parseThemeColor, rgba, type VariantColorsResolver } from "@mantine/core";

const variantColorResolver: VariantColorsResolver = (input) => {
  const resolved = defaultVariantColorsResolver(input);
  const parsed = parseThemeColor({ color: input.color || input.theme.primaryColor, theme: input.theme });

  if (input.variant === "light") {
    return {
      background: rgba(parsed.value, 0.18),
      hover: rgba(parsed.value, 0.26),
      color: `var(--mantine-color-${parsed.color}-2)`,
      border: `1px solid ${rgba(parsed.value, 0.3)}`
    };
  }

  if (input.variant === "subtle") {
    return {
      ...resolved,
      color: `var(--mantine-color-${parsed.color}-2)`,
      hover: rgba(parsed.value, 0.14)
    };
  }

  return resolved;
};

const fontStack = "var(--font-inter), Segoe UI, Arial, sans-serif";

const theme = createTheme({
  primaryColor: "brand",
  primaryShade: 5,
  autoContrast: true,
  focusRing: "auto",
  variantColorResolver,
  defaultRadius: "lg",
  fontFamily: fontStack,
  headings: {
    fontFamily: fontStack,
    fontWeight: "600",
    sizes: {
      h1: { fontSize: "2rem", lineHeight: "1.05", fontWeight: "650" },
      h2: { fontSize: "1.5rem", lineHeight: "1.15" },
      h3: { fontSize: "1.125rem", lineHeight: "1.2" }
    }
  },
  colors: {
    brand: [
      "#f5f0fc",
      "#e8def5",
      "#d4c1ec",
      "#b89ade",
      "#9c79d0",
      "#6a4a99",
      "#553a7e",
      "#422c63",
      "#2f1f48",
      "#1d1330"
    ],
    success: [
      "#e9faf3",
      "#ccf3e0",
      "#9ce8c2",
      "#69dca3",
      "#3fd28a",
      "#22c97a",
      "#179c5e",
      "#0e7a48",
      "#075833",
      "#03361f"
    ],
    warning: [
      "#fff6e5",
      "#ffe8bf",
      "#fcd58c",
      "#fac15a",
      "#f8b13a",
      "#f6a522",
      "#cd8413",
      "#a3660c",
      "#794907",
      "#502d03"
    ],
    danger: [
      "#fdecec",
      "#f9cccc",
      "#f3a1a1",
      "#ed7575",
      "#ea5050",
      "#e63d3d",
      "#bd2a2a",
      "#931e1e",
      "#691414",
      "#400a0a"
    ],
    info: [
      "#e6f5ff",
      "#bee2ff",
      "#8fccff",
      "#5cb5ff",
      "#34a3ff",
      "#1493ee",
      "#0a74c0",
      "#055691",
      "#023b65",
      "#01233e"
    ],
    accent: [
      "#e8fbfa",
      "#c2f4f0",
      "#8fe8e2",
      "#5edcd3",
      "#37d2c7",
      "#1ac9bc",
      "#0d9d92",
      "#077871",
      "#035351",
      "#012f30"
    ],
    dark: [
      "#d5d9e2",
      "#b6bdcc",
      "#8d96a8",
      "#6a7282",
      "#4d5463",
      "#292c33",
      "#1f1f23",
      "#1a1b1e",
      "#121316",
      "#0d0e11"
    ]
  },
  black: "#0d0e11",
  white: "#f4f7fb",
  defaultGradient: {
    from: "brand.4",
    to: "brand.6",
    deg: 135
  },
  components: {
    AppShell: {
      styles: {
        main: {
          background: "var(--portlyn-app-bg)"
        }
      }
    },
    Paper: {
      defaultProps: {
        radius: "lg",
        p: "lg"
      },
      styles: {
        root: {
          backgroundColor: "var(--portlyn-surface)",
          border: "1px solid var(--portlyn-border)",
          boxShadow: "0 1px 2px rgba(0, 0, 0, 0.3)"
        }
      }
    },
    Card: {
      defaultProps: {
        radius: "lg",
        p: "lg"
      },
      styles: {
        root: {
          backgroundColor: "var(--portlyn-surface)",
          border: "1px solid var(--portlyn-border)",
          boxShadow: "0 1px 2px rgba(0, 0, 0, 0.3)"
        }
      }
    },
    Button: {
      defaultProps: {
        radius: "md"
      },
      styles: {
        root: {
          fontWeight: 600,
          letterSpacing: "-0.005em"
        }
      }
    },
    NavLink: {
      defaultProps: {
        variant: "subtle"
      },
      styles: {
        root: {
          borderRadius: "10px",
          color: "var(--portlyn-text-muted)"
        },
        label: {
          fontSize: "0.9rem",
          fontWeight: 500
        },
        section: {
          color: "var(--portlyn-text-dimmed)"
        }
      }
    },
    TextInput: {
      defaultProps: {
        size: "md"
      }
    },
    PasswordInput: {
      defaultProps: {
        size: "md"
      }
    },
    Select: {
      defaultProps: {
        size: "md"
      }
    },
    NumberInput: {
      defaultProps: {
        size: "md"
      }
    },
    Avatar: {
      defaultProps: {
        color: "brand.5",
        variant: "filled"
      },
      styles: {
        placeholder: {
          fontWeight: 600,
          color: "#ffffff"
        }
      }
    },
    Table: {
      styles: {
        table: {
          backgroundColor: "transparent"
        },
        th: {
          borderBottom: "none",
          color: "var(--portlyn-text-dimmed)",
          fontSize: "0.7rem",
          fontWeight: 600,
          letterSpacing: "0.08em",
          textTransform: "uppercase"
        },
        td: {
          borderTop: "1px solid var(--portlyn-border)",
          fontVariantNumeric: "tabular-nums"
        },
        tr: {
          backgroundColor: "transparent"
        }
      }
    },
    Badge: {
      defaultProps: {
        radius: "sm",
        variant: "filled"
      },
      styles: {
        root: {
          fontWeight: 600,
          letterSpacing: "0.02em",
          textTransform: "none"
        }
      }
    },
    Modal: {
      styles: {
        content: {
          backgroundColor: "var(--portlyn-surface)",
          border: "1px solid var(--portlyn-border)",
          boxShadow: "0 16px 48px -24px rgba(0, 0, 0, 0.7)"
        },
        header: {
          backgroundColor: "var(--portlyn-surface)"
        }
      }
    },
    Drawer: {
      styles: {
        content: {
          backgroundColor: "var(--portlyn-surface)"
        },
        header: {
          backgroundColor: "var(--portlyn-surface)"
        },
        body: {
          paddingTop: "0.25rem"
        }
      }
    }
  }
});

export default theme;
