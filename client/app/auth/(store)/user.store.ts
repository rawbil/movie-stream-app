import { create } from "zustand";

type UserStates = {
  //   public_id: string;
  //   email: string;
  //   created_at: Date;
  //   updated_at: Date;
  user: {
    public_id: string;
    email: string;
    created_at: Date;
    updated_at: Date;
  };
};

type UserActions = {
  //   setPublicId: (public_id: UserStates["public_id"]) => void;
  //   setEmail: (email: UserStates["email"]) => void;
  //   setCreatedAt: (created_at: UserStates["created_at"]) => void;
  //   setUpdatedAt: (updated_at: UserStates["updated_at"]) => void;

  setUser: (user: UserStates["user"]) => void;
};

export const useUserStore = create<UserStates & UserActions>()((set) => ({
  //   public_id: "",
  //   email: "",
  //   created_at: new Date(),
  //   updated_at: new Date(),
  //   setPublicId: (public_id) => set({ public_id: public_id }),
  //   setEmail: (email) => set({ email: email }),
  //   setCreatedAt: (created_at) => set({ created_at: created_at }),
  //   setUpdatedAt: (updated_at) => set({ updated_at: updated_at }),
  user: {
    public_id: "",
    email: "",
    created_at: new Date(),
    updated_at: new Date(),
  },
  setUser: (user) => set({ user }),
}));
