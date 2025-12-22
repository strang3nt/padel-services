package tournament

type Gender int

const (
	Male Gender = iota
	Female
	Else
)

type Level int

const (
	Beginner Level = iota
	Intermediate
	Advanced
)

type Person struct {
	Id string
}

type Team struct {
	Person_1   Person
	Person_2   Person
	TeamGender Gender
	TeamLevel  Level
}

func NewTeam(person1, person2 Person, teamGender Gender, teamLevel Level) *Team {
	return &Team{
		Person_1:   person1,
		Person_2:   person2,
		TeamGender: teamGender,
		TeamLevel:  teamLevel,
	}
}

func MakeTeam(person1, person2 Person, teamGender Gender, teamLevel Level) Team {
	return Team{
		Person_1:   person1,
		Person_2:   person2,
		TeamGender: teamGender,
		TeamLevel:  teamLevel,
	}
}
