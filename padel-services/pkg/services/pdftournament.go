package services

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"padelservices/pkg/tournament"
	"path/filepath"
	"strconv"
	"strings"
)

// --- Data Structures (Same as before) ---

type Match struct {
	Court       string
	TeamA       string
	ScoreA      string
	TeamB       string
	ScoreB      string
	RoundNumber int
}

type Round struct {
	RoundNumber int
	Matches     []Match
}

type TournamentData struct {
	Name      string
	StartDate string
	Rounds    []Round
}

type TemplateData struct {
	Tournament TournamentData
}

const chromeExecutable = "google-chrome"

func CreatePdfTournament(data TemplateData, templatePath string, templateFileName string) (string, error) {

	tempHTMLFile, err := executeAndSaveTemplate(templatePath, templateFileName, data)
	if err != nil {
		return "", fmt.Errorf("error executing and saving template: %v", err)
	}
	defer func() {
		err := os.Remove(tempHTMLFile)
		if err != nil {
			log.Printf("error while removing temp file: %v", err)
		}
	}()

	outputFile := "tournament_schedule.pdf"

	if err := generatePDFWithHeadlessChrome(tempHTMLFile, outputFile); err != nil {
		return "", fmt.Errorf("error generating PDF: %v", err)
	}

	return outputFile, nil
}

func executeAndSaveTemplate(tplFilePath string, tplFileName string, data TemplateData) (string, error) {

	tplContent, err := os.ReadFile(filepath.Join(tplFilePath, tplFileName))
	if err != nil {
		return "", fmt.Errorf("reading template file: %w", err)
	}

	t, err := template.New(filepath.Base(tplFileName)).Parse(string(tplContent))
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err = t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	tempFile, err := os.CreateTemp(tplFilePath, "schedule-*.html")
	if err != nil {
		return "", fmt.Errorf("creating temp file: %w", err)
	}

	defer func() {
		err := tempFile.Close()
		if err != nil {
			log.Printf("error while closing temp file: %v", err)
		}
	}()
	if _, err := tempFile.Write(buf.Bytes()); err != nil {
		return "", fmt.Errorf("writing to temp file: %w", err)
	}

	return tempFile.Name(), nil
}

func generatePDFWithHeadlessChrome(inputHTMLPath, outputPath string) error {
	fmt.Printf("Starting PDF generation using %s...\n", chromeExecutable)

	inputURL := "file://" + inputHTMLPath

	args := []string{
		"--headless=new",
		fmt.Sprintf("--print-to-pdf=%s", outputPath),
		"--print-to-pdf-no-header",
		"--landscape",
		"--no-margins",
		"--no-sandbox",
		"--disable-gpu",
		"--font-render-hinting=none",
		"--virtual-time-budget=1000",
		"--default-page-settings={\"pageSize\":\"A4\"}",
		inputURL,
	}

	cmd := exec.Command(chromeExecutable, args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("chrome executable failed: %w. Stderr: %s", err, stderr.String())
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return fmt.Errorf("PDF output file was not created. Chrome executable error log: %s", stderr.String())
	}

	return nil
}

func FromTournamentToTemplateData(tournament tournament.Tournament) TemplateData {
	return TemplateData{
		Tournament: TournamentData{
			Name:      tournament.GetName(),
			StartDate: tournament.GetDateStart().Format("2006-01-02"),
			Rounds: func() []Round {
				var rounds []Round
				for roundIndex, round := range tournament.GetRounds() {
					var matches []Match
					for _, match := range round {

						surnamePerson1TeamA := strings.Split(match.TeamA.Person_1.Id, " ")
						surnamePerson2TeamA := strings.Split(match.TeamA.Person_2.Id, " ")
						surnamePerson1TeamB := strings.Split(match.TeamB.Person_1.Id, " ")
						surnamePerson2TeamB := strings.Split(match.TeamB.Person_2.Id, " ")

						matches = append(matches, Match{
							Court:       strconv.Itoa(match.CourtId),
							TeamA:       surnamePerson1TeamA[len(surnamePerson1TeamA)-1] + ", " + surnamePerson2TeamA[len(surnamePerson2TeamA)-1],
							ScoreA:      "",
							TeamB:       surnamePerson1TeamB[len(surnamePerson1TeamB)-1] + ", " + surnamePerson2TeamB[len(surnamePerson2TeamB)-1],
							ScoreB:      "",
							RoundNumber: roundIndex + 1,
						})
					}
					rounds = append(rounds, Round{
						RoundNumber: roundIndex + 1,
						Matches:     matches,
					})
				}
				return rounds
			}(),
		},
	}
}

func FromTournamentDataToTemplateData(tournament tournament.TournamentData) TemplateData {
	return TemplateData{
		Tournament: TournamentData{
			Name:      tournament.Name,
			StartDate: tournament.Date.Format("2006-01-02"),
			Rounds: func() []Round {
				var rounds []Round
				for roundIndex, round := range tournament.Rounds {
					var matches []Match
					for _, match := range round {

						surnamePerson1TeamA := strings.Split(match.TeamA.Person_1.Id, " ")
						surnamePerson2TeamA := strings.Split(match.TeamA.Person_2.Id, " ")
						surnamePerson1TeamB := strings.Split(match.TeamB.Person_1.Id, " ")
						surnamePerson2TeamB := strings.Split(match.TeamB.Person_2.Id, " ")

						matches = append(matches, Match{
							Court:       strconv.Itoa(match.CourtId),
							TeamA:       surnamePerson1TeamA[len(surnamePerson1TeamA)-1] + ", " + surnamePerson2TeamA[len(surnamePerson2TeamA)-1],
							ScoreA:      "",
							TeamB:       surnamePerson1TeamB[len(surnamePerson1TeamB)-1] + ", " + surnamePerson2TeamB[len(surnamePerson2TeamB)-1],
							ScoreB:      "",
							RoundNumber: roundIndex + 1,
						})
					}
					rounds = append(rounds, Round{
						RoundNumber: roundIndex + 1,
						Matches:     matches,
					})
				}
				return rounds
			}(),
		},
	}
}
