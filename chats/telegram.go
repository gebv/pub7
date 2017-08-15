package chats

import (
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"gopkg.in/telegram-bot-api.v4"
)

func NewTelegram(
	store StateStore,
	logger *zap.SugaredLogger,
) *Telegram {
	return &Telegram{
		w:   &Workspace{},
		s:   store,
		log: logger.Named("TB"),
	}
}

type Telegram struct {
	w   *Workspace
	s   StateStore
	log *zap.SugaredLogger
}

func (m *Telegram) Workspace() *Workspace {
	return m.w
}

func (m *Telegram) Handler(
	api *tgbotapi.BotAPI,
	req tgbotapi.Update,
) (err error) {
	if req.Message == nil {
		err = errors.New("empty message from request")
		return
	}

	chatID := req.Message.Chat.ID
	state, err := m.loadState(chatID)
	if err == ErrNotFound {
		state = NewState(m.keyFromChatID(chatID))
	} else if err != nil {
		m.log.Errorw(
			"load state from chat",
			"chat_id", chatID,
			"err", err,
		)
		return err
	}

	defer func() {
		m.updateState(chatID, state)
	}()

	isStartCmd := false
	if state.ScriptID == "" {
		state.ScriptID = START
		isStartCmd = true
	}

	if req.Message.IsCommand() {
		state.ScriptID = req.Message.Command()
		isStartCmd = true
	}

	// current script

	script := m.w.FindScript(state.ScriptID)
	if script == nil {
		m.log.Errorw(
			"not found script by ID",
			"script_id", state.ScriptID,
			"err", "not found script",
		)
		return ErrNotFoundScript
	}

	if isStartCmd {
		// Приветственное сообщение
		msg := Format2String(script.Title, state.Props)
		api.Send(tgbotapi.NewMessage(chatID, msg))
	}

	// current question

	question := m.w.FindQuestion(state.LastQID)
	if question == nil {
		question = m.w.FindQuestion(script.StartQID)
		state.LastQID = script.StartQID
	}
	if question == nil {
		m.log.Errorw(
			"not found question by ID",
			"last_question_id", state.LastQID,
			"err", "not found question",
		)
		return ErrNotFoundQuestion
	}

	if isStartCmd {
		// Первое сообщение сценария
		if err := question.ExecuteScript(state.Props); err != nil {
			m.log.Warnw(
				"execute script",
				"question_id", question.QID,
				"err", err,
			)
		}
		for _, txt := range question.Texts {
			txt = Format2String(txt, state.Props)
			api.Send(tgbotapi.NewMessage(chatID, txt))
		}
		return nil
	}

	if len(question.ParamName) > 0 {
		state.Props[question.ParamName] = strings.TrimSpace(req.Message.Text)
	}

	// handler response

	if err := m.handler(
		state,
		script,
		question,
		api,
		req,
	); err != nil {
		if err == ErrNotFoundQuestion {
			state.LastQID = script.StartQID
			if err := question.ExecuteScript(state.Props); err != nil {
				m.log.Warnw(
					"execute script",
					"question_id", question.QID,
					"err", err,
				)
			}
			for _, txt := range question.Texts {
				txt = Format2String(txt, state.Props)
				api.Send(tgbotapi.NewMessage(chatID, txt))
			}
			return nil
		}
		m.log.Errorw(
			"handler message",
			"err", err,
			"chat_id", chatID,
		)
		return err
	}

	// after hadnler

	nextQuestion := m.w.FindQuestion(state.LastQID)
	if nextQuestion != nil && nextQuestion.NextQID == MENU {
		nextQuestion = m.FindNextQuestion(nextQuestion.NextQID)
		_, err = api.Send(m.TgMessageByQuestion(chatID, state.Props, nextQuestion))
		state.LastQID = nextQuestion.QID
	}

	return err
}

func (m *Telegram) handler(
	state *State,
	script *Script,
	question *Question,
	api *tgbotapi.BotAPI,
	req tgbotapi.Update,
) error {
	chatID := req.Message.Chat.ID
	text := strings.TrimSpace(req.Message.Text)

	// на основе ответа определяем следующий question ID

	if len(question.Options) > 0 {
		val := strings.ToLower(text)
		for _, opt := range question.Options {
			if strings.ToLower(opt.Key) == val {
				question.NextQID = opt.NextQID
				break
			}
		}
	}

	// находим следующий вопрос
	nextQuestion := m.FindNextQuestion(question.NextQID)
	if nextQuestion == nil {
		m.log.Errorw(
			"find question by ID",
			"question_id", question.NextQID,
			"err", "not found question",
		)
		state.LastQID = ""
		return ErrNotFoundQuestion
	}

	state.LastQID = nextQuestion.QID

	// текущий текст
	// if err := question.ExecuteScript(state.Props); err != nil {
	// 	m.log.Warnw(
	// 		"execute script",
	// 		"question_id", question.QID,
	// 		"err", err,
	// 	)
	// }
	// for _, txt := range question.TextsWithoutLast() {
	// 	txt = Format2String(txt, state.Props)
	// 	api.Send(tgbotapi.NewMessage(chatID, txt))
	// }

	// следующий текст
	if err := nextQuestion.ExecuteScript(state.Props); err != nil {
		m.log.Warnw(
			"execute script",
			"question_id", nextQuestion.QID,
			"err", err,
		)
	}
	for _, txt := range nextQuestion.TextsWithoutLast() {
		txt = Format2String(txt, state.Props)
		api.Send(tgbotapi.NewMessage(chatID, txt))
	}
	_, err := api.Send(m.TgMessageByQuestion(chatID, state.Props, nextQuestion))
	return err
}

func (m *Telegram) loadState(chatID int64) (*State, error) {
	return m.s.Find(m.keyFromChatID(chatID))
}

func (m *Telegram) updateState(chatID int64, obj *State) error {
	obj.ChatID = m.keyFromChatID(chatID)
	return m.s.Update(obj)
}

func (m *Telegram) keyFromChatID(chatID int64) string {
	return fmt.Sprintf("chats:telegram:%d", chatID)
}

func (m *Telegram) FindNextQuestion(
	nextQuestionID string,
) *Question {
	if len(nextQuestionID) == 0 {
		nextQuestionID = MENU
	}

	return m.w.FindQuestion(nextQuestionID)
}

func (m *Telegram) TgMessageByQuestion(
	chatID int64,
	msgArgs map[string]interface{},
	question *Question,
) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, Format2String(question.LastText(), msgArgs))
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	if len(question.Options) > 0 {
		_l := len(question.Options)
		var inlineButtins = make([]tgbotapi.KeyboardButton, _l, _l)
		for i, opt := range question.Options {
			inlineButtins[i] = tgbotapi.NewKeyboardButton(opt.Key)
		}
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				inlineButtins...,
			),
		)
	}

	return msg
}

// func (m *Telegram) FormatMessage(
// 	state *State,
// 	msg string,
// ) string {

// }
