import { useState, type FC } from "react";
import { Page } from "@/components/Page.tsx";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/components/AuthProvider";
import Snackbar from "@mui/material/Snackbar";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Box from "@mui/material/Box";
import Section from "@/components/Section";
import { Tournaments } from "@/api/tournament";
import retrieveTournament from "@/api/retrieveTournament";
interface FormData {
  date: string;
}

export const RetrieveTournamentPage: FC = () => {
  const { bearerToken } = useAuth();
  const [formData, setFormData] = useState<FormData>({ date: "" });
  const navigate = useNavigate();

  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const response = await retrieveTournament(
        bearerToken || "",
        new Date(formData.date),
      );

      const result = (await response.json()) as Tournaments;
      if (response.ok) {
        navigate("/available-tournaments", { state: result });
      } else {
        return (
          <Snackbar
            autoHideDuration={3000}
            onClose={() => {}}
            message="Server not reachable"
          />
        );
      }
    } catch (error) {
      console.error("Submission error:", error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Page>
      <Section title="Select tournament date">
        <form
          onSubmit={(event) => {
            void handleSubmit(event);
          }}
        >
          <Box
            sx={{
              width: "100%",
              display: "flex",
              flexDirection: "column",
              gap: 2,
            }}
          >
            <TextField
              id="outlined-basic"
              label="Insert date"
              variant="outlined"
              type="date"
              required
              slotProps={{ inputLabel: { shrink: true } }}
              onChange={(e) => setFormData({ date: e.target.value })}
            />
            <Button type="submit" disabled={isLoading}>
              {isLoading ? "Submitting..." : "Submit"}
            </Button>
          </Box>
        </form>
      </Section>
    </Page>
  );
};
