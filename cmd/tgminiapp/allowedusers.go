package main

import (
	"fmt"

	"github.com/strang3nt/padel-services/pkg/database"
)

type userData struct {
	LogoPath     string
	SportsCenter string
}

type AllowedUsers struct {
	users map[int64]userData
}

func (au AllowedUsers) GetLogoPath(userId int64) (string, error) {
	if ud, ok := au.users[userId]; ok {
		return ud.LogoPath, nil
	}

	return "", fmt.Errorf("user %d does not exist", userId)
}

func (au AllowedUsers) IsUserAllowed(userId int64) bool {
	_, ok := au.users[userId]
	return ok
}

func MakeAllowerUsersFromUserData(ud []database.UserData) AllowedUsers {
	users := map[int64]userData{}
	for _, x := range ud {
		users[x.Id] = userData{
			LogoPath:     x.Logo,
			SportsCenter: x.SportsCentre,
		}
	}
	return AllowedUsers{
		users: users,
	}
}
