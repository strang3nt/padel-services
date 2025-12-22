package main

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"padelservices/pkg/services"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"padelservices/pkg/database"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var defaultMessage string

var ctx = context.Background()

func main() {
	conn, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	err = database.CreateDatabaseTables(ctx, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while initializing database: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	telegramBotToken := os.Getenv("BOT_TOKEN")

	state := make(map[int64]StateMachine, 0)

	opts := []bot.Option{
		bot.WithDebug(),
		bot.WithDefaultHandler(defaultHandler),
		bot.WithMessageTextHandler("/crea_torneo", bot.MatchTypeExact, createTournamentHandler(&state)),
	}

	b, err := bot.New(telegramBotToken, opts...)

	// Register this in your main function after creating the bot instance
	b.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		chatId := update.Message.Chat.ID
		if s, ok := state[chatId]; ok && s == TournamentCreated {
			return true
		}
		return false
	}, printToPdf(&state, conn))

	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func createTournamentHandler(state *map[int64]StateMachine) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		(*state)[update.Message.Chat.ID] = TournamentCreated
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ora devi messaggio nel seguente formato, tipo di torneo",
		})
	}
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      defaultMessage,
		ParseMode: models.ParseModeMarkdown,
		LinkPreviewOptions: &models.LinkPreviewOptions{
			IsDisabled: bot.True(),
		},
	})
}

func printToPdf(state *map[int64]StateMachine, conn *pgxpool.Pool) bot.HandlerFunc {

	return func(ctx context.Context, b *bot.Bot, update *models.Update) {

		tournamentData := update.Message.Text

		msgScanner := bufio.NewScanner(strings.NewReader(tournamentData))

		msgScanner.Scan()
		availableCourts, err := strconv.Atoi(msgScanner.Text())
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Errore: la prima riga deve contenere il numero di campi disponibili per il torneo.",
			})
			return
		}

		msgScanner.Scan()
		roundsNumber, err := strconv.Atoi(msgScanner.Text())
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Errore: la seconda riga deve contenere il numero di round da disputare durante il torneo.",
			})
			return
		}

		teams, _ := services.MakeTeamsFromMessage(msgScanner)
		tournament := services.CreateTournament("Rodeo", time.Now(), teams, roundsNumber, availableCourts)
		template_data := services.FromTournamentToTemplateData(tournament)
		pdfPath, err := services.CreatePdfTournament(template_data, "template/template_schedule.html")
		if err != nil {
			log.Printf("Error creating PDF: %v", err)
			return
		}
		pdfFile, err := os.Open(pdfPath)
		if err != nil {
			log.Printf("Error opening generated PDF: %v", err)
			return
		}
		err = database.CreateTournament(ctx, conn, &tournament)
		if err != nil {
			log.Printf("Error saving tournament into database: %v", err)
			return
		}
		defer os.Remove(pdfFile.Name())

		b.SendDocument(ctx, &bot.SendDocumentParams{
			ChatID:   update.Message.Chat.ID,
			Document: &models.InputFileUpload{Filename: "tournament_schedule.pdf", Data: bufio.NewReader(pdfFile)},
			Caption:  fmt.Sprintf("Document '%s'.", "tournament_schedule.pdf"),
		})
		(*state)[update.Message.Chat.ID] = Started
	}
}
