"use client";

import { useAuthStore } from "@/app/auth/(store)/auth.store";
import React, { useEffect } from "react";
import ClientLoading from "../Loaders/clientLoading";

export default function AuthProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const initializeAuth = useAuthStore((state) => state.initializeAuth);
  const isLoading = useAuthStore((state) => state.isLoading);



  useEffect(() => {
    initializeAuth();
  }, [initializeAuth]);

  if (isLoading) {
    return <ClientLoading />;
  }
  return <div>{children}</div>;
}
