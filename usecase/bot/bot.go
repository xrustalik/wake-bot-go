package bot

import (
	"wake-bot/usecase/callback"
	"wake-bot/usecase/translation"

	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Service represents a bot's service.
type Service struct {
	tg *tgbotapi.BotAPI
}

// NewService creates a new Service.
func NewBotService(bot *tgbotapi.BotAPI) *Service {
	return &Service{tg: bot}
}

// Button represent a config for inline markup keyboard button
type Button struct {
	text     string
	callback string
}

func NewButton(text, clb string) Button {
	return Button{text: text, callback: clb}
}

// Option sets behaviour for a message config option
type Option interface {
	apply(*tgbotapi.MessageConfig)
}

type replyToOption int

func (r replyToOption) apply(opts *tgbotapi.MessageConfig) {
	opts.ReplyToMessageID = int(r)
}

// WithReplyTo casts int to replyToOption
func WithReplyTo(r int) Option {
	return replyToOption(r)
}

type parseModeOption string

func (o parseModeOption) apply(opts *tgbotapi.MessageConfig) {
	opts.ParseMode = string(o)
}

// WithParseMode casts string to parseModeOption
func WithParseMode(p string) Option {
	return parseModeOption(p)
}

type keyboardOption struct {
	keyboard *tgbotapi.InlineKeyboardMarkup
}

func (o keyboardOption) apply(opts *tgbotapi.MessageConfig) {
	opts.ReplyMarkup = o.keyboard
}

// WithKeyboard casts a tgbotapi.InlineKeyboardMarkup to a keyboardOption
func WithKeyboard(k *tgbotapi.InlineKeyboardMarkup) Option {
	return keyboardOption{k}
}

// SendMessage sends a message to the bot with a set of options.
func (bot *Service) SendMessage(chatID int64, message string, opts ...Option) error {
	msg := tgbotapi.NewMessage(chatID, message)

	msg.ParseMode = "Markdown"
	for _, o := range opts {
		o.apply(&msg)
	}

	_, err := bot.tg.Send(msg)

	return err
}

// SendClarificationMessage creates clarification keyboard, sends a message
func (bot *Service) SendClarificationMessage(chatID int64, replyTo int, languageCode string) error {
	keyboard := bot.MakeInlineKeyboard(
		NewButton(translation.Get(translation.WakeUp, languageCode), callback.WakeUp),
		NewButton(translation.Get(translation.GoToSleep, languageCode), callback.GoToSleep),
	)

	return bot.SendMessage(
		chatID,
		translation.Get(translation.ClarificationQuestion, languageCode),
		WithReplyTo(replyTo),
		WithKeyboard(&keyboard),
	)
}

// SendTimeFormatMessage creates a time format keyboard, sends a message.
func (bot *Service) SendTimeFormatMessage(chatID int64, replyTo int, languageCode string) error {
	timeFormatButtons := bot.MakeInlineKeyboard(
		NewButton(translation.Get(translation.AMPMTimeFormat, languageCode), callback.AMPMTime),
		NewButton(translation.Get(translation.MilitaryTimeFormat, languageCode), callback.MilitaryTime),
	)

	return bot.SendMessage(
		chatID,
		translation.Get(translation.TimeFormatQuestion, languageCode),
		WithReplyTo(replyTo),
		WithKeyboard(&timeFormatButtons),
	)
}

// SendNotValidTimeFormatMessage sends a not valid time format message to the bot.
func (bot *Service) SendNotValidTimeFormatMessage(chatID int64, replyTo int, languageCode, timeFormat string) error {
	return bot.SendMessage(
		chatID,
		fmt.Sprintf(translation.Get(translation.NotValidTimeFormat, languageCode),
			time.Now().UTC().Format(timeFormat)),
		WithReplyTo(replyTo),
	)
}

// MakeInlineKeyboard creates a tgbotapi.InlineKeyboardMarkup from a set of buttons.
func (bot *Service) MakeInlineKeyboard(btns ...Button) tgbotapi.InlineKeyboardMarkup {
	buttons := make([]tgbotapi.InlineKeyboardButton, len(btns))
	for i, btn := range btns {
		buttons[i] = tgbotapi.NewInlineKeyboardButtonData(btn.text, btn.callback)
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons)
}

// AnswerOnCallback sends an AnswerOnCallback message to the user.
func (bot *Service) AnswerOnCallback(id, label string) {
	_, _ = bot.tg.AnswerCallbackQuery(tgbotapi.NewCallback(id, label))
}
