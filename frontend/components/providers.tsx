"use client";

import { Avatar, Center, Loader, MantineProvider, Stack, Text, type CSSVariablesResolver } from "@mantine/core";
import { Notifications, notifications } from "@mantine/notifications";
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode
} from "react";
import { usePathname, useRouter } from "next/navigation";

import { BootstrapWizard } from "@/components/bootstrap/wizard";
import theme from "@/theme";
import {
  getCurrentUser,
  login as loginRequest,
  logoutRequest,
  logout as clearAuthStorage
} from "@/lib/auth";
import { setApiToken, setUnauthorizedHandler } from "@/lib/api";
import type { LoginResponse, User } from "@/lib/types";

interface AuthContextValue {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<LoginResponse>;
  completeAuth: (response: LoginResponse) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

const cssVariablesResolver: CSSVariablesResolver = () => ({
  variables: {},
  light: {},
  dark: {
    "--mantine-color-dimmed": "var(--portlyn-text-dimmed)",
    "--mantine-color-text": "var(--portlyn-text)"
  }
});

export function Providers({ children }: { children: ReactNode }) {
  return (
    <MantineProvider theme={theme} defaultColorScheme="dark" cssVariablesResolver={cssVariablesResolver}>
      <Notifications position="top-right" />
      <AuthProvider>{children}</AuthProvider>
    </MantineProvider>
  );
}

function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();
  const pathname = usePathname();
  const normalizedPath = pathname.length > 1 ? pathname.replace(/\/+$/, "") : pathname;
  const isPublicAuthRoute =
    normalizedPath === "/login" ||
    normalizedPath === "/oidc/callback" ||
    normalizedPath === "/route-login" ||
    normalizedPath === "/route-forbidden";

  const clearSession = useCallback(() => {
    setUser(null);
    setToken(null);
    clearAuthStorage();
  }, []);

  const handleLogout = useCallback(() => {
    void logoutRequest();
    clearSession();
    router.push("/login");
  }, [clearSession, router]);

  const completeAuth = useCallback(
    (response: LoginResponse) => {
      setApiToken(response.token || null);
      setToken(response.token || null);
      setUser({
        ...response.user,
        bootstrap_required:
          response.user.bootstrap_required ?? response.bootstrap_required ?? false
      });
      notifications.show({
        title: "Signed in",
        message: `Role: ${response.user.role}`,
        color: "success"
      });
    },
    []
  );

  useEffect(() => {
    setUnauthorizedHandler(handleLogout);
    return () => setUnauthorizedHandler(null);
  }, [handleLogout]);

  useEffect(() => {
    getCurrentUser()
      .then((currentUser) => {
        setUser(currentUser);
      })
      .catch(() => {
        clearSession();
        if (!isPublicAuthRoute) {
          router.replace("/login");
        }
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, [clearSession, isPublicAuthRoute, router]);

  const login = useCallback(
    async (email: string, password: string) => {
      const response = await loginRequest(email, password);
      if (!response.requires_mfa) {
        completeAuth(response);
      }
      return response;
    },
    [completeAuth]
  );

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      token,
      isLoading,
      isAuthenticated: Boolean(user),
      login,
      completeAuth,
      logout: handleLogout
    }),
    [completeAuth, handleLogout, isLoading, login, token, user]
  );

  const showWizard = Boolean(
    user && (user.bootstrap_required || user.must_change_password) && !user.bootstrap_dismissed
  );

  const handleWizardComplete = useCallback(
    async (updates?: { user?: User; dismissed?: boolean }) => {
      if (updates?.user) {
        setUser((current) => (current ? { ...current, ...updates.user } : updates.user!));
      }
      try {
        const fresh = await getCurrentUser();
        setUser(fresh);
      } catch {
        if (updates?.dismissed) {
          setUser((current) => (current ? { ...current, bootstrap_dismissed: true } : current));
        }
      }
    },
    []
  );

  return (
    <AuthContext.Provider value={value}>
      {showWizard && user ? <BootstrapWizard user={user} onComplete={handleWizardComplete} /> : null}
      {isPublicAuthRoute || (!isLoading && Boolean(user)) ? (
        children
      ) : (
        <Center mih="100vh">
          <Stack gap="sm" align="center">
            <Loader color="brand" />
            <Text c="dimmed">Loading workspace</Text>
          </Stack>
        </Center>
      )}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within Providers");
  }
  return context;
}

export function UserAvatar() {
  const { user } = useAuth();
  const seed = user?.display_name || user?.username || user?.email || "PL";
  const initials = seed.slice(0, 2).toUpperCase();
  return <Avatar radius="xl" size="sm" color="brand.4">{initials}</Avatar>;
}
