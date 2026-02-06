import { type FC } from 'react';
import { Cell, List, Placeholder, Section } from '@telegram-apps/telegram-ui';
import { TournamentData, Tournaments } from './RetrieveTournamentPage';
import { Page } from '@/components/Page.tsx';
import { useLocation } from 'react-router-dom';
import { useAuth } from '@/components/AuthProvider';
import { downloadFile } from '@tma.js/sdk-react';

export const AvailableTournamentsPage: FC = () => {
  const { bearerToken } = useAuth()
  const location = useLocation();

  const onClick = (tournamentData: TournamentData) => {

    fetch(
      `/api/tournament/generate-link`,
      {
        method: "POST",
        headers: {
          'Authorization': `Bearer ${bearerToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(tournamentData),
      }
    ).then((response) => {
      response.json().then(({ user, token }) => {
        const url = `${window.location.origin}/api/tournament/download?user=${user}&token=${token}`
        if (import.meta.env.DEV) {
          window.location.href = url;
        } else {
          downloadFile(
            url,
            `${tournamentData.date}_${tournamentData.name}.pdf`
          )
        }
      }
      )
    });
  }

  const { date, tournaments } = location.state as Tournaments;
  return (
    <Page>
      <Section
        header={`Tournaments at date ${date}`}
      >

        {
          tournaments.length == 0 ? <Placeholder description="No tournaments found at selected date" /> :
            <List>
              {tournaments.map(
                (tournamentData) =>
                  <Cell
                    description={`Participants: ${tournamentData.teams.length}`}
                    onClick={() => onClick(tournamentData)}
                  >
                    {tournamentData.name}
                  </Cell>

              )}      </List>}


      </Section>
    </Page >
  );
};
