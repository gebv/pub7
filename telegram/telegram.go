package telegram

import (
	"fmt"
	"strings"
	"time"

	"github.com/gebv/ff_tgbot/chats"
	"github.com/gebv/ff_tgbot/context"
	"github.com/gebv/ff_tgbot/errors"
	"github.com/gebv/ff_tgbot/nodes"
	"github.com/gebv/ff_tgbot/utils"
	lua "github.com/yuin/gopher-lua"

	"go.uber.org/zap"
	"gopkg.in/telegram-bot-api.v4"
)

func New(
	api *tgbotapi.BotAPI,
	log *zap.SugaredLogger,
	chats chats.Store,
	nodes nodes.Store,
	ts time.Duration, // Задержка перед отпракой сообщения
) *Manager {
	return &Manager{
		log:   log,
		api:   TelegramSenderWithTimeout(api, ts),
		nodes: nodes,
		chats: chats,
	}
}

type Manager struct {
	log   *zap.SugaredLogger
	api   TelegramAPI
	nodes nodes.Store
	chats chats.Store
}

func (m *Manager) Handler(req tgbotapi.Update) (err error) {
	if req.Message == nil {
		err = errors.New("empty message from request")
		return
	}

	chatID := req.Message.Chat.ID
	chat, err := m.FindChat(chatID)
	if err == errors.ErrNotFound {
		chat = chats.NewChat(formatTelegramChatID(chatID))
		chat.NextNodeID = "start"
	}

	m.log.Debugw(
		"new message",
		"chat_id", chatID,
		"is_command", req.Message.IsCommand(),
		"command", req.Message.Command(),
		"last_node_id", chat.NextNodeID,
		"chat_props", fmt.Sprintf("%v", chat.Props),
		"chat_updated_at", chat.UpdatedAt,
	)

	if req.Message.IsCommand() {
		chat.NextNodeID = req.Message.Command()
		// TODO: command arguments
	}

	defer func() {
		if err := m.chats.Update(chat); err != nil {
			m.log.Errorw(
				"save the state of the chat",
				"chat_id", chatID,
				"err", err,
			)
		}
	}()

	return m.handler(
		chat,
		req,
	)
}

func (m *Manager) handler(
	chat *chats.Chat,
	req tgbotapi.Update,
) (err error) {

	////////////////////////////////////////////////////////////////////////////
	// previous node
	////////////////////////////////////////////////////////////////////////////
	previousNode, err := m.FindNode(chat.PreviousNodeID)
	if err == nil && previousNode != nil {
		if len(previousNode.ParamName) > 0 {
			chat.Props[previousNode.ParamName] = req.Message.Text
		}

		if len(previousNode.Options) > 0 {
			val := strings.ToLower(req.Message.Text)
			for _, opt := range previousNode.Options {
				if strings.ToLower(opt.Text) == val {
					if len(opt.NextNodeID) > 0 {
						chat.NextNodeID = opt.NextNodeID
					}
					break
				}
			}
		}
	}

	////////////////////////////////////////////////////////////////////////////
	// next node
	////////////////////////////////////////////////////////////////////////////

	node, err := m.FindNode(chat.NextNodeID)
	if err == errors.ErrNotFound {
		node, err = m.FindNode("start")
		if err != nil {
			return fmt.Errorf("not found node %q or 'start'", chat.NextNodeID)
		}
	}
	chat.PreviousNodeID = node.ID
	chat.NextNodeID = node.NextNodeID
	m.log.Debugw(
		"handler node",
		"node_id", node.ID,
		"node_is_transit", node.IsTransit,
		"node_next_id", node.NextNodeID,
		"node_param_name", node.ParamName,
	)

	ctx := context.Telegram(
		chat.Props,
		m.api,
		req.Message.Chat.ID,
		req.Message.Text,
	)
	L := lua.NewState()
	defer L.Close()
	context.RegisterTelegramContextType(L, ctx)

	////////////////////////////////////////////////////////////////////////////
	// Before
	////////////////////////////////////////////////////////////////////////////

	if err := node.Before(L, ctx); err != nil {
		m.log.Errorw(
			"execute before script",
			"node_id", node.ID,
			"err", err,
		)
		return err
	}
	if ctx.IsAbort() {
		return
	}
	if len(ctx.RedirectTo()) > 0 {
		chat.NextNodeID = ctx.RedirectTo()
		return m.handler(chat, req)
	}

	////////////////////////////////////////////////////////////////////////////
	// handler
	////////////////////////////////////////////////////////////////////////////

	if err := node.Handler(L, ctx); err != nil {
		m.log.Errorw(
			"execute handler script",
			"node_id", node.ID,
			"err", err,
		)
		return err
	}

	for i, text := range node.Text {
		msg := tgbotapi.NewMessage(req.Message.Chat.ID, utils.ExecuteTemplate(text, ctx.Props()))
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

		// show buttons if need
		if i == len(node.Text)-1 {
			// last message
			if len(node.Options) > 0 {
				_l := len(node.Options)
				var inlineButtins = make([]tgbotapi.KeyboardButton, _l, _l)
				for i, opt := range node.Options {
					inlineButtins[i] = tgbotapi.NewKeyboardButton(opt.Text)
				}
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						inlineButtins...,
					),
				)
			}
		}

		m.api.Send(msg)
	}

	if ctx.IsAbort() {
		return
	}
	if len(ctx.RedirectTo()) > 0 {
		chat.NextNodeID = ctx.RedirectTo()
		return m.handler(chat, req)
	}

	////////////////////////////////////////////////////////////////////////////
	// after
	////////////////////////////////////////////////////////////////////////////

	if err := node.After(L, ctx); err != nil {
		m.log.Errorw(
			"execute handler script",
			"node_id", node.ID,
			"err", err,
		)
		return err
	}
	if ctx.IsAbort() {
		return
	}
	if len(ctx.RedirectTo()) > 0 {
		chat.NextNodeID = ctx.RedirectTo()
		return m.handler(chat, req)
	}

	// helpers

	if node.IsTransit {
		return m.handler(chat, req)
	}

	return nil
}

func (m *Manager) FindChat(chatID int64) (*chats.Chat, error) {
	return m.chats.Find(formatTelegramChatID(chatID))
}

func (m *Manager) FindNode(nodeID string) (*nodes.Node, error) {
	return m.nodes.Find(nodeID)
}

func formatTelegramChatID(chatID int64) string {
	return fmt.Sprintf("chats:telegram:%d", chatID)
}
