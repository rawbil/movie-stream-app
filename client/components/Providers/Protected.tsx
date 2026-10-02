"use client";

import { useAuthStore } from "@/app/auth/(store)/auth.store";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function AuthProtected({
  children,
}: {
  children: React.ReactNode;
}) {
  const accessToken = useAuthStore((state) => state.accessToken);
  const isAuthenticated = !!accessToken;
  const router = useRouter();

  useEffect(() => {
    if (!isAuthenticated) {
      router.push("/auth/login");
    }
  }, [isAuthenticated, router]);

  if (!isAuthenticated) return null;

  return children;
}
