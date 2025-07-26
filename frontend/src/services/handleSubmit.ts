import axios from "axios";
import type { GenerateRequest } from "../types/api";
import { enqueueSnackbar } from "notistack";

export const handleVoucherSubmit = async (
  values: GenerateRequest,
  setSeats: React.Dispatch<React.SetStateAction<string[] | null>>
): Promise<void> => {
  try {
    const checkRes = await axios.post("/api/check", {
      flightNumber: values.flightNumber,
      date: values.date,
    });

    if (checkRes.data.exists) {
      enqueueSnackbar("This flight already has seat assignments.", {
        variant: "info",
      });
      return;
    }

    const genRes = await axios.post("/api/generate", values);
    setSeats(genRes.data.seats);

    enqueueSnackbar(`Vouchers generated! Seats: ${genRes.data.seats.join(", ")}`, {
      variant: "success",
    });
  } catch (error: unknown) {
    let errorMessage = "An unexpected error occurred";
    if (axios.isAxiosError(error) && error.response?.data?.error) {
      errorMessage = error.response.data.error;
    } else if (error instanceof Error) {
      errorMessage = error.message;
    }

    enqueueSnackbar(errorMessage, {
      variant: "error",
    });
  }
};
