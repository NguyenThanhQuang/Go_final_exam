import { type Trip } from "../types/trip.types";
import apiClient from "./apiClient";

interface SearchParams {
  from: string;
  to: string;
  date: string;
}

interface ApiResponse {
  thong_bao: string;
  du_lieu: Trip[];
}

export const searchTrips = async (params: SearchParams): Promise<Trip[]> => {
  try {
    const response = await apiClient.get("/trips", { params });
    return response.data.dữ_liệu || [];
  } catch (error) {
    console.error("Lỗi khi tìm kiếm chuyến đi:", error);
    throw error;
  }
};
