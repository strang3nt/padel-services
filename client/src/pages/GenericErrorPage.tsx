import errorImage from "../../assets/error.svg";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Stack from "@mui/material/Stack";
import { type FC } from "react";

const GenericErrorPage: FC<{ error: unknown }> = ({ error }) => {
  return (
    <Box
      display="flex"
      justifyContent="center"
      alignItems="center"
      minHeight="100vh"
      sx={{ backgroundColor: "background.default" }}
    >
      <Stack spacing={2} alignItems="center">
        <Box textAlign="center">
          <img
            alt="Error image"
            src={errorImage}
            style={{
              display: "block",
              margin: "0 auto",
              width: "144px",
              height: "144px",
            }}
          />
          <Typography variant="body2" color="text.secondary">
            {"\n\n\n"}Oops:{" "}
            {typeof error == "string" ? `${error}` : "Unexpected error"}
          </Typography>
        </Box>
      </Stack>
    </Box>
  );
};

export default GenericErrorPage;
