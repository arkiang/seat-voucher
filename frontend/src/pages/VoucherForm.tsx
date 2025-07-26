import React, { useState } from "react";
import {
  Box,
  Button,
  MenuItem,
  Select,
  TextField,
  Typography,
  FormControl,
  InputLabel,
} from "@mui/material";
import { useFormik } from "formik";
import { getGenerateRequestValidationSchema } from "../utils/validators";
import { handleVoucherSubmit } from "../services/handleSubmit";
import { AIRCRAFT_TYPES } from "../constant/aircraft";

const VoucherForm: React.FC = () => {
  const [seats, setSeats] = useState<string[] | null>(null);

  const formik = useFormik({
    initialValues: {
      name: "",
      id: "",
      flightNumber: "",
      date: "",
      aircraft: "",
    },
    validationSchema: getGenerateRequestValidationSchema(),
    onSubmit: async (values) => {
      await handleVoucherSubmit(values, setSeats); // enqueueSnackbar is handled internally
    },
  });

  return (
    <Box maxWidth={500} mx="auto" mt={5} p={3} boxShadow={3} borderRadius={2} bgcolor="#fff">
      <Typography variant="h5" align="center" gutterBottom>
        Crew Seat Voucher Generator
      </Typography>

      <form onSubmit={formik.handleSubmit}>
        <TextField
          label="Crew Name"
          fullWidth
          margin="normal"
          {...formik.getFieldProps("name")}
          error={formik.touched.name && Boolean(formik.errors.name)}
          helperText={formik.touched.name && formik.errors.name}
        />

        <TextField
          label="Crew ID"
          fullWidth
          margin="normal"
          {...formik.getFieldProps("id")}
          error={formik.touched.id && Boolean(formik.errors.id)}
          helperText={formik.touched.id && formik.errors.id}
        />

        <TextField
          label="Flight Number"
          fullWidth
          margin="normal"
          {...formik.getFieldProps("flightNumber")}
          error={formik.touched.flightNumber && Boolean(formik.errors.flightNumber)}
          helperText={formik.touched.flightNumber && formik.errors.flightNumber}
        />

        <TextField
          label="Flight Date (DD-MM-YYYY)"
          fullWidth
          margin="normal"
          {...formik.getFieldProps("date")}
          error={formik.touched.date && Boolean(formik.errors.date)}
          helperText={formik.touched.date && formik.errors.date}
        />

        <FormControl
          fullWidth
          margin="normal"
          error={formik.touched.aircraft && Boolean(formik.errors.aircraft)}
        >
          <InputLabel>Aircraft Type</InputLabel>
          <Select
            {...formik.getFieldProps("aircraft")}
            value={formik.values.aircraft}
            onChange={formik.handleChange}
            onBlur={formik.handleBlur}
            label="Aircraft Type"
          >
            <MenuItem value=""><em>None</em></MenuItem>
            {AIRCRAFT_TYPES.map((type) => (
              <MenuItem key={type} value={type}>
                {type}
              </MenuItem>
            ))}
          </Select>
          {formik.touched.aircraft && formik.errors.aircraft && (
            <Typography variant="caption" color="error">
              {formik.errors.aircraft}
            </Typography>
          )}
        </FormControl>

        <Button
          variant="contained"
          color="primary"
          type="submit"
          fullWidth
          sx={{ mt: 2 }}
        >
          Generate Vouchers
        </Button>
      </form>

      {seats && (
        <Box mt={3} textAlign="center">
          <Typography variant="subtitle1" fontWeight="bold">Assigned Seats:</Typography>
          <Typography>{seats.join(", ")}</Typography>
        </Box>
      )}
    </Box>
  );
};

export default VoucherForm;