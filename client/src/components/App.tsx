import { useMemo } from "react";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { useSignal, miniApp } from "@tma.js/sdk-react";
import { ThemeProvider, CssBaseline } from "@mui/material";

import { routes } from "@/navigation/routes.tsx";
import { AuthProvider } from "./AuthProvider";
import { createTheme } from "@mui/material/styles";
import { themeParams } from "@tma.js/sdk-react";

const getMuiTheme = (
  themeParams: Partial<Record<string, `#${string}`>>,
  isDark: boolean,
) => {
  return createTheme({
    palette: {
      mode: isDark ? "dark" : "light",
      primary: {
        main: themeParams.buttonColor || "#5288c1",
        contrastText: themeParams.buttonTextColor || "#ffffff",
      },
      background: {
        default: themeParams.bgColor || (isDark ? "#17212b" : "#ffffff"),
        paper: themeParams.secondaryBgColor || (isDark ? "#232e3c" : "#f4f4f5"),
      },
      text: {
        primary: themeParams.textColor || (isDark ? "#f5f5f5" : "#000000"),
        secondary: themeParams.hintColor || "#708499",
      },
      error: {
        main: themeParams.destructiveTextColor || "#ec3942",
      },
    },
    shape: {
      borderRadius: 8,
    },
    components: {
      MuiCssBaseline: {
        styleOverrides: {
          body: {
            margin: 0,
            padding: 0,
            width: "100%",
            height: "100%",
            display: "flex",
            justifyContent: "center",
          },
          "#root": {
            width: "100%",
            maxWidth: "600px",
            minHeight: "100vh",
          },
        },
      },
    },
  });
};

export function App() {
  const tp = useSignal(themeParams.state);
  const isDark = useSignal(miniApp.isDark);

  const muiTheme = useMemo(() => {
    return getMuiTheme(tp, isDark);
  }, [tp, isDark]);

  return (
    <ThemeProvider theme={muiTheme}>
      <CssBaseline />
      <AuthProvider>
        <RouterProvider router={createBrowserRouter(routes)} />
      </AuthProvider>
    </ThemeProvider>
  );
}
