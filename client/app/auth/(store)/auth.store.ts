import { useUserStore } from "./user.store";
import { AxiosClient } from "@/lib/axios";
import { create } from "zustand";

let initializeAuthPromise: Promise<void> | null = null;

type AuthStates = {
  accessToken: string | null;
  isLoading: boolean;
};

type AuthActions = {
  setAccessToken: (token: AuthStates["accessToken"]) => void;
  initializeAuth: () => Promise<void>;
  logout: () => Promise<void>;
};

export const useAuthStore = create<AuthStates & AuthActions>()((set) => ({
  accessToken: "",
  isLoading: true,
  setAccessToken: (token) => set({ accessToken: token }),
  initializeAuth: async () => {
    if (initializeAuthPromise) return initializeAuthPromise;

    initializeAuthPromise = (async () => {
      try {
        const response = await AxiosClient.post("/auth/refresh_tokens", {});
        const data = response.data;

        set({ accessToken: data.access_token });
        useUserStore.getState().setUser({
          public_id: data.user.public_id,
          username: data.user.username,
          email: data.user.email,
          created_at: data.user.created_at,
          updated_at: data.user.updated_at,
          role: data.user.role,
        });
      } catch {
        set({ accessToken: null });
        useUserStore.getState().clearUser();
      } finally {
        set({ isLoading: false });
      }
    })();

    return initializeAuthPromise;
  },
  logout: async () => {},
}));
