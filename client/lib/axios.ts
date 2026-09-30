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
      !originalRequest._retry
    ) {
      originalRequest._retry = true;
      try {
        const token = useAuthStore.getState().accessToken;
        if (!token) {
          window.location.href = "/auth/login";
          return Promise.reject(error);
        }
        const response = await AxiosClient.post(
          "/auth/refresh_tokens",
          {},
          { withCredentials: true },
        );
        useAuthStore.getState().setAccessToken(response.data.access_token);
        useUserStore.getState().setUser({
          public_id: response.data.user.public_id,
          email: response.data.user.email,
          created_at: response.data.user.created_at,
          updated_at: response.data.user.updated_at,
        });
        //console.log("New Access Token generated", response.data.access_token);

        // Update the original request's Authorization header
        originalRequest.headers["Authorization"] =
          `Bearer ${response.data.access_token}`;

        return Axios(originalRequest); //retry original request
      } catch (err) {
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
