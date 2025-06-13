import type { Trip } from "./trip.types";

export interface Passenger {
  name?: string;
  phone?: string;
  seatNumber: string;
}

export interface Booking {
  id: string;
  userId: string;
  tripId: string;
  bookingTime: string;
  status: "pending" | "held" | "confirmed" | "cancelled" | "expired";
  heldUntil?: string;
  paymentStatus?: "pending" | "paid" | "failed";
  totalAmount: number;
  passengers: Passenger[];
  ticketCode?: string;
  createdAt: string;
  updatedAt: string;
  tripInfo?: Trip;
}
