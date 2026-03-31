import {
  ChooseTournamentType,
  CreateTournamentPage,
} from "@/pages/CreateTournamentPage/CreateTournamentPage";

import {
  AddPlayersPage,
  AddTeamPage,
  teamsLoader,
} from "@/pages/CreateTournamentPage/AddTeamsPage";
import { RetrieveTournamentPage } from "@/pages/RetrieveTournamentPage";
import { AvailableTournamentsPage } from "@/pages/AvailableTournamentsPage";
import { MenuPage } from "@/pages/IndexPage/MenuPage";
import { RouteObject, useRouteError } from "react-router-dom";
import { ProtectedRoute } from "@/pages/ProtectedRoute";
import GenericErrorPage from "@/pages/GenericErrorPage";

const RouteError = () => {
  const error = useRouteError();
  return <GenericErrorPage error={error} />;
};

export const routes: RouteObject[] = [
  {
    path: "/",
    element: <ProtectedRoute />,
    errorElement: <RouteError />,
    children: [
      { index: true, element: <MenuPage /> },
      {
        path: "/create-tournament",
        element: <CreateTournamentPage />,
        children: [
          {
            path: "/create-tournament/add-players",
            element: <AddPlayersPage />,
            loader: teamsLoader,
          },
          {
            path: "/create-tournament",
            element: <ChooseTournamentType />,
          },
          {
            path: "/create-tournament/add-team",
            element: <AddTeamPage />,
          },
        ],
      },
      { path: "/retrieve-tournament", element: <RetrieveTournamentPage /> },
      { path: "/available-tournaments", element: <AvailableTournamentsPage /> },
    ],
  },
];
