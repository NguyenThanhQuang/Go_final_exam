export interface Seat {
  seatNumber: string;
  status: "available" | "held" | "booked";
}

export interface Route {
  from: { name: string };
  to: { name: string };
}

export interface Trip {
  id: string;
  companyName: string;
  route: Route;
  departureTime: string;
  expectedArrivalTime: string;
  price: number;
  seats: Seat[];
  availableSeats: number;
}
