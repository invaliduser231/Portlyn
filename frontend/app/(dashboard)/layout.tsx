import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import type { ReactNode } from "react";

import { DashboardShell } from "@/components/layout/dashboard-shell";

export default async function ProtectedLayout({ children }: { children: ReactNode }) {
  const cookieStore = await cookies();
  const sessionToken = cookieStore.get("portlyn_session")?.value;

  if (!sessionToken) {
    redirect("/login");
  }

  return <DashboardShell>{children}</DashboardShell>;
}
