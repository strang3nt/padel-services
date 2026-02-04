import { Text, Button, ButtonCell, InlineButtons, Input, List, Section, Select } from '@telegram-apps/telegram-ui';
import { InlineButtonsItem } from '@telegram-apps/telegram-ui/dist/components/Blocks/InlineButtons/components/InlineButtonsItem/InlineButtonsItem';
import { useState, type FC } from 'react';
import { Form, Link, Outlet, useLoaderData, LoaderFunctionArgs, replace, useNavigate } from 'react-router-dom';
import { IoIosAdd, IoIosSave, IoIosClose } from "react-icons/io";
import { Page } from '@/components/Page';
import { Team } from '@/pages/RetrieveTournamentPage';
import { useAuth } from '@/components/AuthProvider';

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
  return <Page>
    <InlineButtons mode="plain">
      <Link to='/create-tournament/add-team' replace>
        <InlineButtonsItem text="Add">
          <IoIosAdd />
        </InlineButtonsItem>
      </Link>
      <Link to='/create-tournament/tournament-type'>
        <InlineButtonsItem text="Save">
          <IoIosSave />
        </InlineButtonsItem>
      </Link>
    </InlineButtons>
    {
      teams.length == 0 ?
        'No teams added yet'
        : <List> {teams.map(({
          person1: { id: teammate1 }, person2: { id: teammate2 }, gender }) =>
          <ButtonCell
            after={< IoIosClose />}>
            {`${teammate1}, ${teammate2}, ${genderToString(gender)}`}
          </ButtonCell>)}
        </List >
    }
  </Page>
}



export const tournamentStore = {
  teams: [] as Team[],
  addTeam: (team: Team) => {
    tournamentStore.teams.push(team);
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

  // Redirect back to the list after adding
  return replace('/create-tournament');
}

export const AddTeamPage: FC = () => {

  return (
    <Page>
      <Form method="post">
        <Input
          header="First teammate"
          name="teammate1" // Important!
          type="text"
          required
        />
        <Input
          header="Second teammate"
          name="teammate2" // Important!
          type="text"
          required
        />
        <Select header="Select" name="gender">
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
              navigate("/")
            } else {
              throw Error("Could not create tournament")
            }
          })
        }
        return <>
          <Input
            header="Courts available"
            name="courtsAvailable" // Important!
            type="number"
            onChange={(e) => setFormData({
              ...formData,
              availableCourts: parseInt(e.target.value, 10)
            })}
            required
          />
          <Input
            header="Number of rounds"
            name="roundsNumber" // Important!
            type="number"
            onChange={(e) => setFormData({
              ...formData,
              roundsNumber: parseInt(e.target.value, 10)
            })}
            required
          />
          <Text>
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
          </Text>
          <Button onClick={() => sendRodeoTournament(
            selectedTournament,
            new Date(formData.tournamentDate),
            tournamentStore.teams,
            formData.roundsNumber,
            formData.availableCourts

          )}>Send</Button>
        </>

      default:
        return "Tournament type not supported";
    }
  }

  return <Page>
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
