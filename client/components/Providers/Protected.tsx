"use client";

import { useAuthStore } from "@/app/auth/(store)/auth.store";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function AuthProtected({ children }: { children: React.ReactNode }) {
  const accessToken = useAuthStore((state) => state.accessToken);
  const router = useRouter();

  useEffect(() => {
    if (!accessToken) router.replace("/auth/login");
  }, [accessToken, router]);

  if (!accessToken) return null;
  return children;
}
