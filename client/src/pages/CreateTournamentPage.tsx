import { Button, Input, List, Section, Select, Cell, Placeholder, Snackbar } from '@telegram-apps/telegram-ui';
import { useState, type FC } from 'react';
import { Form, Outlet, useLoaderData, LoaderFunctionArgs, replace, useNavigate } from 'react-router-dom';
import { IoIosAdd, IoIosClose, IoIosArrowForward } from "react-icons/io";
import { Page } from '@/components/Page';
import { Team } from '@/pages/RetrieveTournamentPage';
import { useAuth } from '@/components/AuthProvider';
import { Link } from '@/components/Link/Link';
import { bem } from '@/css/bem';

const [, e] = bem('display-data');

export const CreateTournamentPage: FC = () => {
  return <Section
    header="Tournament creation"
    footer="Complete the steps to create a tournament"
  >
    <Outlet />
  </Section>

};

function genderToString(g: number): string {
  switch (g) {
    case 0:
      return "Male"
    case 1:
      return "Female"
    case 2:
      return "Mixed"
    default:
      throw Error("Team gender not recognized");
  }
}

export const AddTeamsPage: FC = () => {

  const loaderData = useLoaderData() as { teams: Team[] } | null;
  const teams = loaderData?.teams || [];
  const [_, setRefresh] = useState(false)

  const handleDeleteTeam = (index: number) => {
    tournamentStore.removeTeam(index);
    setRefresh(prev => !prev)
  };

  return <Page>
    <List>
      <Link to='/create-tournament/add-team'>
        <Cell
          after={<IoIosAdd />}
        >
          Add Team
        </Cell>
      </Link>
      <Link to='/create-tournament/tournament-type' state={teams}>
        <Cell
          after={<IoIosArrowForward />}
        >
          Next
        </Cell>
      </Link>
      <Section
        header='Teams added'
      >
        {
          teams.length == 0 ?
            <Placeholder
              description="No teams added yet"
            />
            : <List> {teams.map(({
              person1: { id: teammate1 }, person2: { id: teammate2 }, gender }, i) =>
              <Cell
              
                key={i}
                after={<IoIosClose />}
                onClick={() => handleDeleteTeam(i)}
                subtitle={`${genderToString(gender)} team`}
              >
                {`${teammate1}, ${teammate2}`}
              </Cell>)}
            </List >
        }
      </Section>
    </List>
  </Page>
}

export const tournamentStore = {
  teams: [] as Team[],
  addTeam: (team: Team) => {
    tournamentStore.teams.push(team);
  },
  removeTeam: (index: number) => {
    tournamentStore.teams.splice(index, 1)
  },
  getTeams: () => tournamentStore.teams,
};

export async function teamsLoader() {
  const teams = tournamentStore.getTeams();
  return { teams };
}

export async function addTeamAction({ request }: LoaderFunctionArgs) {
  const formData = await request.formData();

  const newTeam: Team = {
    person1: { id: formData.get("teammate1") as string },
    person2: { id: formData.get("teammate2") as string },
    gender: parseInt(formData.get("gender") as string, 10),
  };

  tournamentStore.addTeam(newTeam);

  return replace('/create-tournament');
}

export const AddTeamPage: FC = () => {

  return (
    <Page>
      <Form method="post">
        <Input
          header="First teammate"
          name="teammate1"
          type="text"
          required
        />
        <Input
          header="Second teammate"
          name="teammate2"
          type="text"
          required
        />
        <Select header="Select gender" name="gender">
          <option value={0}>Male</option>
          <option value={1}>Female</option>
          <option value={2}>Mixed</option>
        </Select>
        <Button type="submit" stretched>Add Team</Button>
      </Form>
    </Page>
  );
}

interface FormData {
  availableCourts: number;
  roundsNumber: number;
  tournamentDate: string;
}

interface NotificationContent {
  title: string;
  description: string;
  onClose?: () => void;
}

export const ChooseTournamentType: FC = () => {

  const getMatchesPerTeam = (
    teamsNumber: number,
    totalRounds: number,
    availableCourts: number): [number, number, number] => {

    let matchesPerTeam = totalRounds

    while (matchesPerTeam > 0) {

      const totalParticipations = teamsNumber * matchesPerTeam
      const totalMatchesFloat = totalParticipations / 2.0

      if (Number.isInteger(totalMatchesFloat)) {

        const matchesPerTurn = totalMatchesFloat / totalRounds

        if (matchesPerTurn <= availableCourts && teamsNumber > matchesPerTeam) {
          return [totalMatchesFloat, matchesPerTurn, matchesPerTeam]
        }
      }

      matchesPerTeam -= 1
    }

    return [0, 0.0, 0]
  }

  const [formData, setFormData] = useState<FormData>({
    availableCourts: 0,
    roundsNumber: 0,
    tournamentDate: "",
  })
  const { bearerToken } = useAuth()
  const [selectedTournament, setTournament] = useState("Rodeo");
  const navigate = useNavigate()
  const [open, setOpen] = useState<null | NotificationContent>(null);

  const renderSwitch = (): React.ReactNode => {
    switch (selectedTournament) {

      case "Rodeo":
        const sendRodeoTournament = (
          tournamentType: string,
          dateStart: Date,
          teams: Team[],
          totalRounds: number,
          availableCourts: number
        ) => {
          fetch(`/api/create-tournament?tournamentType=${tournamentType}&dateStart=${dateStart.toISOString()}&totalRounds=${totalRounds}&availableCourts=${availableCourts}`, {
            method: 'POST',
            headers: {
              'Authorization': `Bearer ${bearerToken}`,
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(teams)
          }).then(response => {
            if (response.ok) {

              setOpen(
                {
                  title: "Tournament creation success",
                  description: "Tournament created and saved successfully",
                  onClose: () => { navigate("/") }
                }
              )
            } else {
              setOpen(
                {
                  title: "Tournament creation failed",
                  description: "Try again later",
                  onClose: () => { navigate("/") }
                }
              )
            }
          })
        }
        return <List>
          <Input
            header="Courts available"
            name="courtsAvailable"
            type="number"
            onChange={(e) => setFormData({
              ...formData,
              availableCourts: parseInt(e.target.value, 10)
            })}
            required
          />
          <Input
            header="Number of rounds"
            name="roundsNumber"
            type="number"
            onChange={(e) => setFormData({
              ...formData,
              roundsNumber: parseInt(e.target.value, 10)
            })}
            required
          />
          <span className={e('line-value')}>
            {
              (() => {
                let totalMatches, matchesPerTurn, matchesPerTeam;
                [totalMatches, matchesPerTurn, matchesPerTeam] = getMatchesPerTeam(tournamentStore.teams.length, formData.roundsNumber, formData.availableCourts)
                if (totalMatches == 0) {
                  return "Configuration is not valid"
                } else {
                  return `Current configuration translates to ${matchesPerTeam} matches per team and at most ${Math.ceil(matchesPerTurn)} matches per turn.`
                }

              })()
            }
          </span>
          <Button stretched onClick={() => sendRodeoTournament(
            selectedTournament,
            new Date(formData.tournamentDate),
            tournamentStore.teams,
            formData.roundsNumber,
            formData.availableCourts

          )}>Send</Button>
        </List>

      default:
        return "Tournament type not supported";
    }
  }

  return <Page>
    {open && <Snackbar
      onClose={() => {
        setOpen(null);
        open?.onClose?.();
      }}
      children={open.title}
      description={open.description}
    />}
    <Form>
      <Select
        header="Select"
        name="tournamentType"
        value={selectedTournament}
        onChange={e => setTournament(e.target.value)}
      >
        <option>Rodeo</option>
      </Select>
      <Input
        header="Tournament date"
        name="tournamentDate"
        type="date"
        onChange={(e) => setFormData({
          ...formData,
          tournamentDate: e.target.value
        })
        }
        required
      />
      {
        renderSwitch()
      }
    </Form>
  </Page>
}
