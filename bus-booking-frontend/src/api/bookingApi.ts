import { type Booking } from "../types/booking.types";
import apiClient from "./apiClient";

export interface CreateBookingPayload {
  tripId: string;
  seatNumbers: string[];
}

interface BookingApiResponse<T> {
  "thông báo": string;
  dữ_liệu?: T;
  lỗi?: string;
}

export const createBooking = async (
  payload: CreateBookingPayload
): Promise<BookingApiResponse<Booking>> => {
  try {
    const response = await apiClient.post<BookingApiResponse<Booking>>(
      "/bookings",
      payload
    );
    console.log("Response từ API createBooking:", response.data);
    return response.data;
  } catch (error: any) {
    console.error(
      "Lỗi khi tạo booking (giữ chỗ):",
      error.response?.data || error.message
    );

    if (
      error.response &&
      error.response.data &&
      (error.response.data.lỗi || error.response.data["thông báo"])
    ) {
      return error.response.data as BookingApiResponse<Booking>;
    }
    return {
      "thông báo": "Lỗi không xác định khi giữ chỗ.",
      lỗi: error.message || "Unknown error",
    };
  }
};

export const getBookingDetails = async (
  bookingId: string
): Promise<Booking | null> => {
  try {
    const response = await apiClient.get<BookingApiResponse<Booking>>(
      `/bookings/${bookingId}`
    );
    console.log("Response từ API getBookingDetails:", response.data);

    if (response.data && response.data.dữ_liệu) {
      return response.data.dữ_liệu;
    } else {
      if (response.data && response.data.lỗi) {
        console.warn(
          `getBookingDetails (API msg): ${response.data.lỗi} cho bookingId ${bookingId}.`
        );
      } else {
        console.warn(
          `getBookingDetails: Key 'dữ_liệu' không tìm thấy cho bookingId ${bookingId} hoặc dữ liệu rỗng.`
        );
      }
      return null;
    }
  } catch (error: any) {
    console.error(
      `Lỗi khi lấy chi tiết booking ${bookingId}:`,
      error.response?.data || error.message
    );
    return null;
  }
};

export const getMyBookings = async (): Promise<Booking[]> => {
  try {
    const response = await apiClient.get<BookingApiResponse<Booking[]>>(
      "/bookings/my"
    );
    console.log("Response từ API getMyBookings:", response.data);

    if (response.data && response.data.dữ_liệu) {
      return response.data.dữ_liệu;
    } else {
      if (response.data && response.data.lỗi) {
        console.warn(`getMyBookings (API msg): ${response.data.lỗi}`);
      } else {
        console.warn(
          "getMyBookings: Key 'dữ_liệu' không tìm thấy hoặc dữ liệu rỗng."
        );
      }
      return [];
    }
  } catch (error: any) {
    console.error(
      "Lỗi khi lấy lịch sử đặt vé:",
      error.response?.data || error.message
    );
    return [];
  }
};
