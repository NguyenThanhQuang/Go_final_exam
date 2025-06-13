import apiClient from "./apiClient";

export interface RegisterPayload {
  name: string;
  email: string;
  phone: string;
  password: string;
}

export interface LoginPayload {
  email: string;
  password: string;
}

interface AuthResponse {
  "thông báo": string;
  token?: string;
  dữ_liệu?: {
    id: string;
    name: string;
    email: string;
    phone: string;
  };
}

export const registerUser = async (
  payload: RegisterPayload
): Promise<AuthResponse> => {
  const response = await apiClient.post<AuthResponse>(
    "/auth/register",
    payload
  );
  return response.data;
};

export const loginUser = async (
  payload: LoginPayload
): Promise<AuthResponse> => {
  const response = await apiClient.post<AuthResponse>("/auth/login", {
    email: payload.email,
    password: payload.password,
  });
  return response.data;
};
