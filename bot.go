package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

type TGWeatherAPI struct {
	db         *DbClient
	log        *zerolog.Logger
	offset     int
	token      string
	client     *http.Client
	endpoint   string
	weatherApi *WeatherApi
}

type WeatherApi struct {
	Host   string
	ApiKey string
}

func NewTGWeatherAPI(log *zerolog.Logger, cfg Config, db *DbClient) *TGWeatherAPI {

	return &TGWeatherAPI{
		db:       db,
		log:      log,
		client:   &http.Client{},
		endpoint: cfg.Host,
		token:    cfg.Token,
		weatherApi: &WeatherApi{
			Host:   cfg.WeatherHost,
			ApiKey: cfg.WeatherApiKey,
		},
	}
}

//Parse Message receives a message from updates and forms new query accordingly
func (tg *TGWeatherAPI) ParseMessage(ctx context.Context, message Message) (url.Values, error) {

	var q = url.Values{}
	q.Add("chat_id", strconv.Itoa(message.Chat.Id))

	switch message.Text {
	case "/help", "/start", "/h":

		row := NewReplyKeyboardRow(KeyboardButton{
			Text: "/subscribe",
		}, KeyboardButton{
			Text: "/unsubscribe",
		})

		replymarkup := ReplyKeyboardMarkup{
			Keyboard:        NewReplyKeyboard(row),
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		}

		data, err := json.Marshal(replymarkup)
		if err != nil {
			return q, fmt.Errorf("couldnt marshal reply_markup object:%w", err)
		}

		q.Add("reply_markup", string(data))
		q.Add("text", "Please select these options or click on the flags for today's national holidays")

	case "/subscribe":
		tg.log.Debug().Msg("subscription requested")
		q.Add("text", "to subscribe send your location")

	case "/unsubscribe":
		tg.log.Debug().Msg("links menu requested")
		q.Add("text", "you will be unsubscribed")
		tg.Unsubscribe(ctx, message.Chat.Id)

	case "":
		if message.Location != nil {
			tg.Subscribe(ctx, message)
			q.Add("text", "Please select  UTC time for weather updates in millitary format[0-24hr]: HH:MM")
			q.Add("parse_mode", "HTML")
			q.Add("reply_to_message_id", strconv.Itoa(message.Id))
		}

	default:
		if num, err := parseTime(message.Text); err != nil {
			tg.log.Error().Err(err).Send()
			tg.log.Debug().Msg("invalid input requested")
			q.Add("text", `I didn't understand. For help type:/help`)
		} else {
			tg.UpsertTime(ctx, message.Chat.Id, num)
			q.Add("text", "your time for updates recorded successfully")
		}
	}

	return q, nil

}

//AddSubscription adds subscription from the Message
func (tg *TGWeatherAPI) Subscribe(ctx context.Context, msg Message) error {
	tg.log.Trace().Msg("started")
	defer tg.log.Trace().Msg("exited")

	sub := Subscription{
		ChatId: msg.Chat.Id,
		Lon:    msg.Location.Longitude,
		Lat:    msg.Location.Latitude,
	}

	_, err := tg.db.UpsertOne(ctx, sub)
	return err
}

func (tg *TGWeatherAPI) Unsubscribe(ctx context.Context, chatId int) error {
	tg.log.Trace().Msg("started")
	defer tg.log.Trace().Msg("exited")
	filter := bson.D{{"chat_id", chatId}}
	sub, err := tg.db.FindOne(ctx, filter)
	if err != nil {
		tg.log.Error().Err(err).Send()
	}
	err = tg.db.DeleteOne(ctx, sub.ID)
	if err != nil {
		tg.log.Error().Err(err).Send()
	}
	return nil

}

func (tg *TGWeatherAPI) UpsertTime(ctx context.Context, chatId int, time int) {
	tg.log.Trace().Msg("started")
	defer tg.log.Trace().Msg("exited")
	filter := bson.D{{"chat_id", chatId}}
	sub, err := tg.db.FindOne(ctx, filter)
	if err != nil {
		tg.log.Error().Err(err).Send()
		return
	}
	sub.UpdateTime = time

	_, err = tg.db.UpsertOne(ctx, sub)
	if err != nil {
		tg.log.Error().Err(err).Send()
	}

}

//send request with url parameters to TG
func (tg *TGWeatherAPI) sendMessage(method string, q url.Values) ([]byte, error) {
	tg.log.Trace().Msg("started")
	defer tg.log.Trace().Msg("exited")

	u, err := url.Parse(tg.endpoint)
	if err != nil {
		tg.log.Fatal().Err(err).Send()
	}
	u.RawQuery = q.Encode()
	u.Path = path.Join("bot"+tg.token, method)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("bot client can't form http request: %w", err)
	}
	req.Header.Set("Content-type", "application/json; charset=utf-8")

	resp, err := tg.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("bot client can't send the request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("bot client can't read from response: %w", err)

	}

	return body, nil

}

//GetUpdates receives all the updates from the Telegram API from the current offset
func (tg *TGWeatherAPI) GetUpdates() []Update {
	tg.log.Trace().Msg("GetUpdates started")
	defer tg.log.Trace().Msg("GetUpdates exited")

	q := url.Values{}
	q.Add("offset", strconv.Itoa(tg.offset))

	u, err := url.Parse(tg.endpoint)
	if err != nil {
		tg.log.Fatal().Err(err).Msg("couldn't parse the endpoint")
	}
	u.RawQuery = q.Encode()

	data, err := tg.sendMessage(getUpdate, q)
	if err != nil {
		tg.log.Error().Err(err).Send()
	}

	var resp APIResponseUpdates
	json.Unmarshal(data, &resp)
	updates := resp.Result

	tg.log.Debug().Int("length of updates", len(updates)).Int("offset", tg.offset).Send()
	if len(updates) > 0 {
		tg.offset = updates[len(updates)-1].UpdateId + 1
	}
	return updates

}

func (tg *TGWeatherAPI) PushWeatherUpdates(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)

	for t := range ticker.C {
		min := t.Minute()
		hour := t.Hour()
		updateTime := hour*100 + min
		tg.log.Debug().Int("update time", updateTime).Msg("update time for ticker")
		subs := tg.getSubscriptions(ctx, updateTime)
		tg.sendWeatherUpdates(subs)

	}

}

func (tg *TGWeatherAPI) sendWeatherUpdates(subs []Subscription) {

	for _, sub := range subs {
		upd, err := tg.doWeatherRequest(sub)
		if err != nil {
			tg.log.Error().Err(err).Send()
		}
		tg.log.Info().Str("update weather", fmt.Sprintf("%+v", upd))
		text := markdownWeatherReply(*upd)
		q := url.Values{}
		q.Add("chat_id", strconv.Itoa(sub.ChatId))
		q.Add("text", text)
		q.Add("parse_mode", "HTML")
		tg.sendMessage(sendMessage, q)

	}

}

func (tg *TGWeatherAPI) getSubscriptions(ctx context.Context, update int) []Subscription {
	tg.log.Trace().Msg("started")
	defer tg.log.Trace().Msg("exited")
	filter := bson.D{{"update_time", update}}

	subs, err := tg.db.Find(ctx, filter)
	if err != nil {
		tg.log.Error().Err(err).Send()
	}
	return subs

}

func (b *TGWeatherAPI) doWeatherRequest(sub Subscription) (*WeatherUpdate, error) {

	b.log.Trace().Msg("doWeatherRequest started")
	defer b.log.Trace().Msg("doWeatherRequest exited")

	q := url.Values{}

	q.Add("appid", b.weatherApi.ApiKey)
	q.Add("lat", fmt.Sprint(sub.Lat))
	q.Add("lon", fmt.Sprint(sub.Lon))
	q.Add("units", "metric")

	u, err := url.Parse(b.weatherApi.Host)
	if err != nil {
		return nil, fmt.Errorf("doWeatherRequest: bot client can't parse holiday api host: %w", err)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("doWeatherRequest: bot client can't form http request: %w", err)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doWeatherRequest: bot client can't send the request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b.log.Warn().Int("error", resp.StatusCode).Msg("response for weather api denied")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("doWeatherRequest bot client can't read from response: %w", err)

	}

	var upd WeatherUpdate

	err = json.Unmarshal(body, &upd)
	if err != nil {
		return nil, fmt.Errorf("doWeatherUpdate: couldn't unmarshal body to struct: %w", err)
	}

	return &upd, nil

}

func markdownWeatherReply(weather WeatherUpdate) string {
	var reply strings.Builder

	fmt.Fprintf(&reply, "<b>%s</b>: <b>%.2fdegC</b>\n",
		strings.ToUpper(weather.Name), weather.Main.Temp)
	fmt.Fprintf(&reply, "Feels like <b>%.2fdegC</b>. %s\n",
		weather.Main.Temp, strings.Title(weather.Weather[0].Description))
	fmt.Fprintf(&reply, "Wind: <b>%.2fm/s</b>, gusts of <b>%.2fm/s</b>\n",
		weather.Wind.Speed, weather.Wind.Gust)

	return reply.String()

}
