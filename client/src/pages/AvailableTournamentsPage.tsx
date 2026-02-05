import { type FC } from 'react';
import { List } from '@telegram-apps/telegram-ui';
import { TournamentData, Tournaments } from './RetrieveTournamentPage';
import { Page } from '@/components/Page.tsx';
import { useLocation } from 'react-router-dom';
import { useAuth } from '@/components/AuthProvider';
import { postEvent } from '@tma.js/sdk-react';
import { DisplayData } from '@/components/DisplayData/DisplayData';

export const AvailableTournamentsPage: FC = () => {
  const { bearerToken } = useAuth()
  const location = useLocation();

  const onClick = (tournamentData: TournamentData) => {

    const apiRequest = fetch(
      `/api/tournament/generate-link`,
      {
        method: "POST",
        headers: {
          'Authorization': `Bearer ${bearerToken}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(tournamentData),
      }
    )
    apiRequest.then((response) => {
      response.json().then(({ user, token }) => {
        const url = `${window.location.origin}/api/tournament/download?user=${user}&token=${token}`
        if (import.meta.env.DEV) {
          window.location.href = url;
        } else {
          postEvent('web_app_request_file_download', {
            url: url,
            file_name: `${tournamentData.date}_${tournamentData.name}.pdf`
          })
        }
      }
      )
    });
  }

  const { date, tournaments } = location.state as Tournaments;
  return (
    <Page>
      <List>
        <DisplayData
          header={`Tournaments at ${date}`}
          rows={
            tournaments
              .map((tournamentData) => (
                {
                  title: tournamentData.name,
                  value: <span onClick={() => onClick(tournamentData)}>
                    Participants: {tournamentData.teams.length}
                  </span>
                }
              ))
          }
        />
      </List>
    </Page>
  );
};
