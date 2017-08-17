package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/gebv/ff_tgbot/admin"
	"github.com/gebv/ff_tgbot/chats"
	"github.com/gebv/ff_tgbot/logger"
	"github.com/gebv/ff_tgbot/nodes"
	"github.com/gebv/ff_tgbot/tarantool"
	"github.com/gebv/ff_tgbot/telegram"
	"github.com/labstack/echo"
	tgbotapi "gopkg.in/telegram-bot-api.v4"

	tnt "github.com/tarantool/go-tarantool"
	cli "gopkg.in/urfave/cli.v2"
)

var VERSION = "dev"
var log *zap.SugaredLogger

var ChatsStore chats.Store
var NodesStore nodes.Store
var TarantoolConnection *tnt.Connection

func main() {
	log = logger.NewLogger(zap.DebugLevel).Sugar()

	app := &cli.App{}
	app.Name = "ff_tgbot"
	app.Version = VERSION
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "file", Value: "example.toml"},
	}
	app.Before = func(c *cli.Context) (err error) {
		dat, err := ioutil.ReadFile(c.String("file"))
		if err != nil {
			log.Fatalw(
				"read scripts from file",
				"file_path", c.String("file"),
				"err", err,
			)
			return err
		}

		TarantoolConnection, err = tarantool.SetupFromENV()
		if err != nil {
			log.Fatalw(
				"setup tarantool connection",
				"err", err,
			)
			return err
		}

		// TODO: testing
		// ChatsStore = chats.NewInmemory()
		ChatsStore = chats.NewTarantool(TarantoolConnection)
		NodesStore = nodes.NewInMemoryStoreNodes()

		if err := NodesStore.LoadFromToml(dat); err != nil {
			log.Errorw(
				"unmarshal toml",
				"err", err,
			)
			return err
		}
		return nil
	}
	app.Commands = []*cli.Command{
		{
			Name: "admin-run",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "port", Value: "8085"},
				&cli.StringFlag{Name: "host", Value: ""},
				&cli.StringFlag{Name: "path_views", Value: "./views/*.tpl"},
			},
			Action: func(c *cli.Context) error {
				e := echo.New()
				admin.SetLogger(log.Named("admin"))
				admin.Setup(e, c.String("path_views"))
				admin.TarantoolConnection = TarantoolConnection

				addrs := fmt.Sprintf("%s:%s", c.String("host"), c.String("port"))
				var appSignal = make(chan struct{}, 2)
				go func() {
					e.Start(addrs)
					appSignal <- struct{}{}
				}()

				osSignal := make(chan os.Signal, 2)
				close := make(chan struct{})
				signal.Notify(
					osSignal,
					os.Interrupt,
					syscall.SIGTERM,
				)

				go func() {

					defer func() {
						close <- struct{}{}
					}()

					select {
					case <-osSignal:
						log.Error("signal completion of the process: OS")
					case <-appSignal:
						log.Error("signal completion of the process: internal (http server, etc..)")
					}

					// TODO: destroy services
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()
					e.Shutdown(ctx)
				}()

				<-close
				os.Exit(0)
				return nil
			},
		},
		{
			Name: "run",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "token",
					Value: "399087349:AAFYX064SOZ9j9iAPC1V33I52n6IIQLen6U",
					Usage: "telegram bot access token",
				},
			},
			Action: func(c *cli.Context) error {
				bot, err := tgbotapi.NewBotAPI(c.String("token"))
				if err != nil {
					log.Fatalw(
						"setup telegram bot API",
						"err", err,
					)
					return err
				}
				bot.Debug = false
				log.Infof("Authorized on account %s", bot.Self.UserName)

				u := tgbotapi.NewUpdate(0)
				u.Timeout = 60
				updates, err := bot.GetUpdatesChan(u)
				if err != nil {
					log.Fatalw(
						"setup updater from telegram bot API",
						"err", err,
					)
					return err
				}

				mng := telegram.New(
					bot,
					log.Named("tg"),
					ChatsStore,
					NodesStore,
					time.Millisecond*500,
				)

				for update := range updates {
					if update.Message == nil {
						continue
					}

					go func() {
						if err := mng.Handler(update); err != nil {
							log.Errorw(
								"handler",
								"err", err,
							)
						}
					}()
				}

				return nil
			},
		},
		// {
		// 	Name: "tgbot-run",
		// 	Flags: []cli.Flag{
		// 		&cli.StringFlag{
		// 			Name:  "token",
		// 			Value: "399087349:AAFYX064SOZ9j9iAPC1V33I52n6IIQLen6U",
		// 			Usage: "telegram bot access token",
		// 		},
		// 	},
		// 	Action: func(c *cli.Context) error {
		// 		bot, err := tgbotapi.NewBotAPI(c.String("token"))
		// 		if err != nil {
		// 			log.Fatalw(
		// 				"setup telegram bot API",
		// 				"err", err,
		// 			)
		// 			return err
		// 		}
		// 		bot.Debug = false
		// 		log.Infof("Authorized on account %s", bot.Self.UserName)

		// 		u := tgbotapi.NewUpdate(0)
		// 		u.Timeout = 60
		// 		updates, err := bot.GetUpdatesChan(u)
		// 		if err != nil {
		// 			log.Fatalw(
		// 				"setup updater from telegram bot API",
		// 				"err", err,
		// 			)
		// 			return err
		// 		}

		// 		for update := range updates {
		// 			if update.Message == nil {
		// 				continue
		// 			}

		// 			chatID := fmt.Sprintf("chats:telegram:%d", update.Message.Chat.ID)
		// 			state, _ := StateStore.Find(chatID)
		// 			if state != nil {
		// 				log.Debugw(
		// 					"state for current chat",
		// 					"chat_id", update.Message.Chat.ID,
		// 					"last_question_id", state.LastQID,
		// 					"current_script_id", state.ScriptID,
		// 					"props", fmt.Sprintf("%+v", state.Props),
		// 					"last_updated_at", state.UpdatedAt.String(),
		// 				)
		// 			}

		// 			go TelegramManager.Handler(
		// 				bot,
		// 				update,
		// 			)
		// 		}

		// 		return nil
		// 	},
		// },
	}

	app.Run(os.Args)
}
