import React, { useState } from "react";
import {
  Box,
  Button,
  MenuItem,
  Select,
  TextField,
  Typography,
  FormControl,
  InputLabel, FormControlLabel, Checkbox, FormGroup,
} from "@mui/material";
import { useFormik } from "formik";
import { getGenerateRequestValidationSchema } from "../utils/validators";
import { handleVoucherSubmit } from "../services/handleSubmit";
import { AIRCRAFT_TYPES } from "../constant/aircraft";

const VoucherForm: React.FC = () => {
  const [seats, setSeats] = useState<string[] | null>(null);
  const [selectedSeats, setSelectedSeats] = useState<string[]>([]);

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
      // Add selected seats to the request body if any are selected
      const requestData = {
        ...values,
        ...(selectedSeats.length > 0 && { seats: selectedSeats })
      };

      await handleVoucherSubmit(requestData, setSeats); // enqueueSnackbar is handled internally

      setSelectedSeats([]);
    },
  });

  const handleSeatCheckboxChange = (seat: string, checked: boolean) => {
    if (checked) {
      setSelectedSeats(prev => [...prev, seat]);
    } else {
      setSelectedSeats(prev => prev.filter(s => s !== seat));
    }
  };

  const handleSelectAllSeats = (checked: boolean) => {
    if (checked && seats) {
      setSelectedSeats([...seats]);
    } else {
      setSelectedSeats([]);
    }
  };

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
          label="Flight Date (DD-MM-YY)"
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
            <MenuItem value="">
              <em>None</em>
            </MenuItem>
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

        {seats && seats.length > 0 && (
            <Box mt={3} mb={2}>
              <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
                Select seats to regenerate the voucher:
              </Typography>

              <FormControlLabel
                  control={
                    <Checkbox
                        checked={selectedSeats.length === seats.length}
                        indeterminate={selectedSeats.length > 0 && selectedSeats.length < seats.length}
                        onChange={(e) => handleSelectAllSeats(e.target.checked)}
                    />
                  }
                  label="Select All"
                  sx={{ mb: 1 }}
              />

              <FormGroup>
                {seats.map((seat) => (
                    <FormControlLabel
                        key={seat}
                        control={
                          <Checkbox
                              checked={selectedSeats.includes(seat)}
                              onChange={(e) => handleSeatCheckboxChange(seat, e.target.checked)}
                          />
                        }
                        label={seat}
                    />
                ))}
              </FormGroup>
            </Box>
        )}

        <Button
            variant="contained"
            color="primary"
            type="submit"
            fullWidth
            sx={{ mt: 2 }}
        >
          {seats ? 'Regenerate Vouchers' : 'Generate Vouchers'}
        </Button>
      </form>

      {seats && (
        <Box mt={3} textAlign="center">
          <Typography variant="subtitle1" fontWeight="bold">
            Assigned Seats:
          </Typography>
          <Typography>{seats.join(", ")}</Typography>
        </Box>
      )}
    </Box>
  );
};

export default VoucherForm;
