import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import CircularProgress from '@mui/material/CircularProgress';
import Typography from '@mui/material/Typography';

function LoadingPage({ title = "Loading Content", message = "Please wait a moment..." }) {
  return (
    <Box
      display="flex"
      justifyContent="center"
      alignItems="center"
      minHeight="100vh"
      sx={{ backgroundColor: 'background.default' }}
    >
      <Stack spacing={2} alignItems="center">
        <CircularProgress size={60} thickness={4} color="primary" />

        <Box textAlign="center">
          <Typography variant="h6" component="h1" color="text.primary">
            {title}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            {message}
          </Typography>
        </Box>
      </Stack>
    </Box>
  );
};

export default LoadingPage
