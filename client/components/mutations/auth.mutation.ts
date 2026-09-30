import { Axios } from "@/lib/axios";
import { LoginParams, RegisterParams } from "@/lib/schema";

//! Register User
export async function RegisterUser(data: RegisterParams) {
  try {
    const response = await Axios.post("/auth/register", data);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response;
    }
    throw error;
  }
}

//! Login User
export async function LoginUser(data: LoginParams) {
  try {
    const response = await Axios.post("/auth/login", data);
    return response.data;
  } catch (error: any) {
    if (error.response) {
      throw error.response;
    }
    throw error;
  }
}
