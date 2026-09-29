import axios from "axios";

const NEXT_API_URL =
  process.env.NEXT_PUBLIC_API_URL || "https://localhost:8080/api/v1";

export const Axios = axios.create({
  baseURL: NEXT_API_URL,
  withCredentials: true,
});
