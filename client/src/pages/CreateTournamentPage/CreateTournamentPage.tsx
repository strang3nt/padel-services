import { useState, type FC } from "react";
import { Form, Outlet, useNavigate } from "react-router-dom";
import Box from "@mui/material/Box";
import InputLabel from "@mui/material/InputLabel";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Select from "@mui/material/Select";
import MenuItem from "@mui/material/MenuItem";
import Section from "@/components/Section";
import FormControl from "@mui/material/FormControl";
import { Page } from "@/components/Page";
import { getMatchesPerTeam } from "./rodeoTournament";
import { getMatchesPerPerson } from "./singlePlayerRodeoTournament";

export const CreateTournamentPage: FC = () => {
  return (
    <Page>
      <Section title="Tournament creation">
        <Outlet />
      </Section>
    </Page>
  );
};

export interface TournamentSetupData {
  availableCourts: number;
  roundsNumber: number;
  tournamentDate: string;
  numberOfTeams: number;
  selectedTournament: string;
}

interface TournamentParamsProps {
  formData: TournamentSetupData;
  setFormData: React.Dispatch<React.SetStateAction<TournamentSetupData>>;
  helperText: () => string;
  quantityDescription: string;
}

const TournamentParams: React.FC<TournamentParamsProps> = ({
  formData,
  setFormData,
  helperText,
  quantityDescription,
}) => {
  return (
    <>
      <TextField
        label="Courts available"
        name="courtsAvailable"
        type="number"
        onChange={(e) =>
          setFormData({
            ...formData,
            availableCourts: parseInt(e.target.value, 10),
          })
        }
        required
      />
      <TextField
        label="Number of rounds"
        name="roundsNumber"
        type="number"
        onChange={(e) =>
          setFormData({
            ...formData,
            roundsNumber: parseInt(e.target.value, 10),
          })
        }
        required
      />
      <TextField
        label={quantityDescription}
        name="numberOfTeams"
        type="number"
        helperText={helperText()}
        onChange={(e) =>
          setFormData({
            ...formData,
            numberOfTeams: parseInt(e.target.value, 10),
          })
        }
        required
      />
    </>
  );
};

export const ChooseTournamentType: FC = () => {
  const [formData, setFormData] = useState<TournamentSetupData>({
    availableCourts: 0,
    roundsNumber: 0,
    tournamentDate: "",
    numberOfTeams: 0,
    selectedTournament: "Rodeo",
  });
  const navigate = useNavigate();

  const handleNextStep = (url: string) => {
    return () => navigate(url, { state: formData });
  };

  const renderSwitch = (): React.ReactNode => {
    switch (formData.selectedTournament) {
      case "Rodeo": {
        return (
          <>
            <TournamentParams
              formData={formData}
              setFormData={setFormData}
              helperText={() => {
                const [totalMatches, matchesPerTurn, matchesPerTeam] =
                  getMatchesPerTeam(
                    formData.numberOfTeams,
                    formData.roundsNumber,
                    formData.availableCourts,
                  );
                if (totalMatches === 0) {
                  return "Configuration is not valid";
                } else {
                  return `Current configuration translates to ${matchesPerTeam} matches per team and at most ${Math.ceil(matchesPerTurn)} matches per turn.`;
                }
              }}
              quantityDescription="Number of teams"
            />
            <Button
              type="button"
              variant="contained"
              size="large"
              fullWidth
              onClick={handleNextStep("/create-tournament/add-teams")}
              disabled={!formData.tournamentDate || !formData.numberOfTeams}
            >
              Next: Add teams
            </Button>
          </>
        );
      }

      case "SinglePlayerRodeo": {
        return (
          <>
            <TournamentParams
              formData={formData}
              setFormData={setFormData}
              helperText={() => {
                const [totalMatches, matchesPerTurn, matchesPerPerson] =
                  getMatchesPerPerson(
                    formData.numberOfTeams,
                    formData.roundsNumber,
                    formData.availableCourts,
                  );
                if (totalMatches === 0) {
                  return "Configuration is not valid";
                } else {
                  return `Current configuration translates to ${Math.floor(matchesPerPerson)} matches per person, at most ${Math.ceil(matchesPerTurn)} matches per turn, ${totalMatches} total matches.`;
                }
              }}
              quantityDescription="Number of people"
            />
            <Button
              type="button"
              variant="contained"
              size="large"
              fullWidth
              onClick={handleNextStep("/create-tournament/add-players")}
              disabled={!formData.tournamentDate || !formData.numberOfTeams}
            >
              Next: Add Players
            </Button>
          </>
        );
      }
      default:
        return "Tournament type not supported";
    }
  };

  return (
    <Form>
      <Box
        sx={{
          width: "100%",
          display: "flex",
          flexDirection: "column",
          gap: 2,
        }}
      >
        <FormControl>
          <InputLabel id="tournament-label">Tournament type</InputLabel>
          <Select
            labelId="tournament-label"
            label="Tournament Type"
            name="tournamentType"
            value={formData.selectedTournament}
            onChange={(e) =>
              setFormData({ ...formData, selectedTournament: e.target.value })
            }
          >
            <MenuItem value={"Rodeo"}>Rodeo</MenuItem>
            <MenuItem value={"SinglePlayerRodeo"}>Single Player Rodeo</MenuItem>
          </Select>
        </FormControl>
        <TextField
          label="Tournament date"
          name="tournamentDate"
          type="date"
          slotProps={{ inputLabel: { shrink: true } }}
          onChange={(e) =>
            setFormData({
              ...formData,
              tournamentDate: e.target.value,
            })
          }
          required
        />
        {renderSwitch()}
      </Box>
    </Form>
  );
};
