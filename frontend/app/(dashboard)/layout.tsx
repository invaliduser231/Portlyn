"use client";

import { useRouter } from "next/navigation";
import { useEffect, type ReactNode } from "react";

import { DashboardShell } from "@/components/layout/dashboard-shell";
import { useAuth } from "@/components/providers";

export default function ProtectedLayout({ children }: { children: ReactNode }) {
  const router = useRouter();
  const { isAuthenticated, isLoading } = useAuth();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.replace("/login");
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading || !isAuthenticated) {
    return null;
  }

  return <DashboardShell>{children}</DashboardShell>;
}
