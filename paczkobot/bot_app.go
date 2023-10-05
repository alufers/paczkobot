package paczkobot

import (
	"log"
	"net/http"
	"time"

	"github.com/alufers/paczkobot/httphelpers"
	"github.com/alufers/paczkobot/inpostextra"
	"github.com/alufers/paczkobot/tghelpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type BotApp struct {
	Bot *tgbotapi.BotAPI
	DB  *gorm.DB

	CommandDispatcher     *tghelpers.CommandDispatcher
	BaseHTTPClient        *http.Client
	NotificationsService  *NotificationsService
	TrackingService       *TrackingService
	AskService            *tghelpers.AskService
	TranslationService    *TranslationService
	InpostService         inpostextra.InpostService
	FollowService         *FollowService
	InpostScannerService  *InpostScannerService
	PackagePrinterService *PackagePrinterService
	ArchiveService        *ArchiveService
	ImageScanningService  *ImageScanningService
}

func NewBotApp(b *tgbotapi.BotAPI, DB *gorm.DB) (a *BotApp) {
	a = &BotApp{
		Bot: b,
		DB:  DB,
	}
	a.BaseHTTPClient = httphelpers.NewClientWithLogger()
	a.BaseHTTPClient.Timeout = time.Second * 25
	a.AskService = tghelpers.NewAskService(a.Bot)
	a.CommandDispatcher = tghelpers.NewCommandDispatcher(
		b,
		a.AskService,
	)

	a.CommandDispatcher.RegisterCommands(
		&StartCommand{App: a, ExtraHelp: []tghelpers.Helpable{
			&AvailableProvidersExtraHelp{},
			&AuthorExtraHelp{},
		}},
		&TrackCommand{App: a},
		&FollowCommand{App: a},
		&PackagesCommand{App: a},
		&UnfollowCommand{App: a},
		&SetNameCommand{App: a},
		&UnfollowAllCommand{App: a},
		&InpostLoginCommand{App: a},
		&InpostLogoutCommand{App: a},
		&InpostScanCommand{App: a},
		&InpostAccountsCommand{App: a},
		&InpostQrCommand{App: a},
		&ArchivedCommand{App: a},
		&InpostOpenCommand{App: a},
	)

	a.NotificationsService = NewNotificationsService(a, []PackageNotificationHook{
		&InpostPackageNotificationHook{},
	})
	a.TrackingService = NewTrackingService(a)
	a.TranslationService = NewTranslationService()
	a.InpostService = inpostextra.NewInpostService(a.BaseHTTPClient)
	a.FollowService = NewFollowService(a)
	a.InpostScannerService = NewInpostScannerService(a)
	a.PackagePrinterService = NewPackagePrinterService()
	a.ArchiveService = NewArchiveService(a)
	a.ImageScanningService = NewImageScanningService(a)
	a.CommandDispatcher.RegisterUpdateHooks(
		a.ImageScanningService,
		a.AskService,
	)
	return
}

func (a *BotApp) Run() {
	if err := MigrateBadInpostAccounts(a.DB); err != nil {
		log.Fatalf("Failed to migrate bad inpost accounts: %v", err)
	}

	log.Printf("Flushing enqueued notifications...")
	if err := a.NotificationsService.FlushEnqueuedNotifications(); err != nil {
		log.Fatalf("Failed to flush enqueued notifications: %v", err)
	}
	log.Printf("Done flushing enqueued notifications!")

	go a.TrackingService.RunAutomaticTrackingLoop()

	log.Printf("Telegram bot is starting...")
	if err := a.CommandDispatcher.RequestSetMyCommands(); err != nil {
		log.Fatalf("Failed to set my commands: %v", err)
	}

	if err := a.CommandDispatcher.RunUpdateLoop(); err != nil {
		log.Fatalf("Failed to run update loop: %v", err)
	}
}
