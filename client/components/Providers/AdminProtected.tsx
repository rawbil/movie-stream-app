"use client";

import { useAuthStore } from "@/app/auth/(store)/auth.store";
import { useUserStore } from "@/app/auth/(store)/user.store";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function AdminProtected({ children }: { children: React.ReactNode }) {
  const accessToken = useAuthStore((state) => state.accessToken);
  const role = useUserStore((state) => state.user.role);
  const router = useRouter();
  const isAdmin = role.toUpperCase() === "ADMIN";

  useEffect(() => {
    if (!accessToken) router.replace("/auth/login");
    else if (!isAdmin) router.replace("/");
  }, [accessToken, isAdmin, router]);

  if (!accessToken || !isAdmin) return null;
  return children;
}
