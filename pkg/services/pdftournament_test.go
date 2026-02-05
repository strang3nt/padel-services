package services

import (
	"os"
	"testing"
)

func TestCreatePDFTournament_VerifyCreation(t *testing.T) {

	data := TemplateData{
		Tournament: TournamentData{
			Name:      "Test Tournament",
			StartDate: "2024-10-01",
			Rounds:    []Round{},
		},
	}

	tt := MakeTournamentPdfGenerator()

	pdfPath, err := tt.CreatePdfTournament(data, Rodeo, "file_name")
	if err != nil {
		t.Fatalf("CreatePdfTournament returned an error: %v", err)
	}
	t.Run("Assertion_1_PDFCreation", func(t *testing.T) {
		if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
			t.Errorf("Expected PDF file to be created at %s, but it does not exist", pdfPath)
		}
	})

	// Clean up the generated PDF file after test
	defer func() {
		err := os.Remove(pdfPath)
		if err != nil {
			t.Logf("error while removing generated PDF file: %v", err)
		}
	}()
}
