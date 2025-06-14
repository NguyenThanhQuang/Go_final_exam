import axios from "axios";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

if (!API_BASE_URL) {
  throw new Error("Biến môi trường VITE_API_BASE_URL chưa được thiết lập!");
}

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

export default apiClient;
