import { createTheme } from "@mantine/core";

const theme = createTheme({
  primaryColor: "brand",
  primaryShade: 5,
  defaultRadius: "lg",
  fontFamily: "Inter, Segoe UI, Arial, sans-serif",
  headings: {
    fontFamily: "Inter, Segoe UI, Arial, sans-serif",
    fontWeight: "600",
    sizes: {
      h1: { fontSize: "2rem", lineHeight: "1.05" },
      h2: { fontSize: "1.5rem", lineHeight: "1.1" },
      h3: { fontSize: "1.125rem", lineHeight: "1.2" }
    }
  },
  colors: {
    brand: ["#f1ecf8", "#e1d5f0", "#d0bee9", "#bfa7e1", "#ae90da", "#9d79d2", "#8760bb", "#6d4c96", "#553b74", "#3c2952"],
    dark: ["#d5d9e2", "#b6bdcc", "#8d96a8", "#6a7282", "#4d5463", "#292c33", "#1f1f23", "#1a1b1e", "#121316", "#0d0e11"]
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
          background: "#121316"
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
          backgroundColor: "#1f1f23",
          boxShadow: "0 18px 40px -28px rgba(0, 0, 0, 0.9)"
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
          backgroundColor: "#1f1f23",
          boxShadow: "0 18px 40px -28px rgba(0, 0, 0, 0.9)"
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
          letterSpacing: "-0.01em"
        }
      }
    },
    NavLink: {
      defaultProps: {
        variant: "subtle"
      },
      styles: {
        root: {
          borderRadius: "12px",
          color: "#9aa3b2"
        },
        label: {
          fontSize: "0.9rem",
          fontWeight: 500
        },
        section: {
          color: "#7e8795"
        }
      }
    },
    TextInput: {
      defaultProps: {
        size: "md"
      }
    },
    Select: {
      defaultProps: {
        size: "md"
      }
    },
    Table: {
      styles: {
        table: {
          backgroundColor: "transparent"
        },
        th: {
          borderBottom: "none",
          color: "#7e8795",
          fontSize: "0.7rem",
          fontWeight: 600,
          letterSpacing: "0.08em",
          textTransform: "uppercase"
        },
        td: {
          borderTop: "1px solid rgba(255, 255, 255, 0.04)"
        },
        tr: {
          backgroundColor: "transparent"
        }
      }
    },
    Badge: {
      defaultProps: {
        radius: "xl"
      },
      styles: {
        root: {
          fontWeight: 600,
          letterSpacing: "0.01em"
        }
      }
    },
    Modal: {
      styles: {
        content: {
          backgroundColor: "#292a2d",
          boxShadow: "0 24px 60px -30px rgba(0, 0, 0, 0.85)"
        },
        header: {
          backgroundColor: "#292a2d"
        }
      }
    },
    Drawer: {
      styles: {
        content: {
          backgroundColor: "#292a2d"
        },
        header: {
          backgroundColor: "#292a2d"
        },
        body: {
          paddingTop: "0.25rem"
        }
      }
    }
  }
});

export default theme;
