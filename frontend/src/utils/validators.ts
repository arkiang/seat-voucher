import * as Yup from "yup";
import {AIRCRAFT_TYPES} from "../constant/aircraft.ts";

export const getGenerateRequestValidationSchema = () => {
  return Yup.object().shape({
    name: Yup.string().required("Crew name is required"),
    id: Yup.string().required("Crew ID is required"),
    flightNumber: Yup.string()
    .matches(/^[A-Z]{2}\d{1,4}$/, "Invalid flight number")
    .required("Required"),
    date: Yup.string()
    .matches(/^\d{2}-\d{2}-\d{2}$/, "Date must be DD-MM-YY")
    .required("Required"),
    aircraft: Yup.string()
    .required("Aircraft type is required")
    .oneOf(AIRCRAFT_TYPES, "Invalid aircraft")
    .required("Required"),
  });
}