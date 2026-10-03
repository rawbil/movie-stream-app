import { useAuthStore } from "@/app/auth/(store)/auth.store";
import { useUserStore } from "@/app/auth/(store)/user.store";
import axios from "axios";

const NEXT_API_URL =
  process.env.NEXT_PUBLIC_API_URL || "https://localhost:8080/api/v1";

export const Axios = axios.create({
  baseURL: NEXT_API_URL,
  withCredentials: true,
});

export const AxiosClient = axios.create({
  baseURL: NEXT_API_URL,
  withCredentials: true,
});

type RefreshResponse = {
  access_token: string;
  user: {
    public_id: string;
    email: string;
    created_at: Date;
    updated_at: Date;
  };
};

let refreshAccessTokenPromise: Promise<string> | null = null;

async function refreshAccessToken() {
  if (!refreshAccessTokenPromise) {
    refreshAccessTokenPromise = AxiosClient.post<RefreshResponse>(
      "/auth/refresh_tokens",
      {},
    )
      .then(({ data }) => {
        useAuthStore.getState().setAccessToken(data.access_token);
        useUserStore.getState().setUser({
          public_id: data.user.public_id,
          email: data.user.email,
          created_at: data.user.created_at,
          updated_at: data.user.updated_at,
        });

        return data.access_token;
      })
      .finally(() => {
        refreshAccessTokenPromise = null;
      });
  }

  return refreshAccessTokenPromise;
}

// Add Bearer token globally to all requests
Axios.interceptors.request.use(
  (config) => {
    const accessToken = useAuthStore.getState().accessToken;
    if (accessToken) {
      config.headers.Authorization = `Bearer ${accessToken}`;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

// 401 interceptor
Axios.interceptors.response.use(
  (res) => res,
  async (error) => {
    const originalRequest = error.config;
    if (
      error.response &&
      error.response.status === 401 &&
      originalRequest &&
      !originalRequest._retry
    ) {
      originalRequest._retry = true;
      try {
        const token = useAuthStore.getState().accessToken;
        if (!token) {
          window.location.href = "/auth/login";
          return Promise.reject(error);
        }
        const accessToken = await refreshAccessToken();

        // Update the original request's Authorization header
        originalRequest.headers = originalRequest.headers ?? {};
        originalRequest.headers["Authorization"] =
          `Bearer ${accessToken}`;

        return Axios(originalRequest); //retry original request
      } catch (err) {
        useAuthStore.getState().setAccessToken(null);
        if (
          window.location.pathname !== "/auth/login" &&
          window.location.pathname !== "/auth/register"
        ) {
          window.location.href = "/auth/login";
        }

        return Promise.reject(err);
      }
    }
    return Promise.reject(error);
  },
);
