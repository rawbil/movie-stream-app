import { create } from "zustand";

type AuthStates = {
  accessToken: string | null;
};

type AuthActions = {
  setAccessToken: (token: AuthStates["accessToken"]) => void;
};

export const useAuthStore = create<AuthStates & AuthActions>()((set) => ({
  accessToken: "",
  setAccessToken: (token) => set({ accessToken: token }),
}));
