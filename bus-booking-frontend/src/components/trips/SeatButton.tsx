import React from "react";
import { Button } from "react-bootstrap";
import type { Seat } from "../../types/trip.types";

interface SeatButtonProps {
  seat: Seat;
  isSelected: boolean;
  onSeatClick: (seatNumber: string) => void;
  disabled?: boolean;
}

const SeatButton: React.FC<SeatButtonProps> = ({
  seat,
  isSelected,
  onSeatClick,
  disabled,
}) => {
  let variant: string;
  let isSeatDisabled = disabled;

  if (seat.status === "booked") {
    variant = "secondary";
    isSeatDisabled = true;
  } else if (seat.status === "held") {
    variant = "warning";
    isSeatDisabled = true;
  } else if (isSelected) {
    variant = "primary";
  } else {
    variant = "outline-success";
  }

  return (
    <Button
      variant={variant}
      onClick={() => !isSeatDisabled && onSeatClick(seat.seatNumber)}
      disabled={isSeatDisabled}
      className="m-1 p-2 text-center"
      style={{ minWidth: "50px", fontSize: "0.9rem" }}
    >
      {seat.seatNumber}
    </Button>
  );
};

export default SeatButton;
