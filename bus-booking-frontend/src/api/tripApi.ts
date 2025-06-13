import type { Trip } from "../types/trip.types";
import apiClient from "./apiClient";

interface TripDetailsApiResponse {
  "thông báo": string;
  dữ_liệu: Trip;
}

export const getAllTrips = async (): Promise<Trip[]> => {
  try {
    const response = await apiClient.get("/trips");
    if (response.data && response.data.dữ_liệu) {
      return response.data.dữ_liệu;
    } else {
      return [];
    }
  } catch (error) {
    console.error("Lỗi trong getAllTrips function (tripApi.ts):", error);
    return [];
  }
};

interface SearchParams {
  from: string;
  to: string;
  date: string;
}
export const searchTrips = async (params: SearchParams): Promise<Trip[]> => {
  try {
    const response = await apiClient.get("/trips", { params });
    if (response.data && response.data.dữ_liệu) {
      return response.data.dữ_liệu;
    } else {
      return [];
    }
  } catch (error) {
    console.error("Lỗi khi tìm kiếm chuyến đi:", error);
    throw error;
  }
};

export const getTripDetails = async (tripId: string): Promise<Trip | null> => {
  try {
    const response = await apiClient.get<TripDetailsApiResponse>(
      `/trips/${tripId}`
    );
    if (response.data && response.data.dữ_liệu) {
      return response.data.dữ_liệu;
    } else {
      return null;
    }
  } catch (error) {
    console.error(`Lỗi khi lấy chi tiết chuyến đi ${tripId}:`, error);
    return null;
  }
};
