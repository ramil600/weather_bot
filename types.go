package main

//bot will get updates or send message to chat
const (
	getUpdate   = "getUpdates"
	sendMessage = "sendMessage"
)

type APIResponse struct {
	Ok     bool          `json:"ok"`
	Result []interface{} `json:"result,omitempty"`
}

type APIResponseUpdates struct {
	Ok     bool
	Result []Update
}

//Update is a Telegram object that we receive every time an user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

//Message is a Telegram object that can be found in an update.
type Message struct {
	Id       int       `json:"message_id"`
	Text     string    `json:"text"`
	Chat     Chat      `json:"chat"`
	Location *Location `json:"location"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// A Telegram Chat indicates the conversation to which the message belongs.
type Chat struct {
	Id int `json:"id"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool               `json:"one_time_keyboard,omitempty"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}
type ReplyMarkup struct {
	ReplyMarkup ReplyKeyboardMarkup `json:"reply_markup"`
}

func NewReplyKeyboardRow(btns ...KeyboardButton) []KeyboardButton {
	keyrow := []KeyboardButton{}
	keyrow = append(keyrow, btns...)
	return keyrow
}

func NewReplyKeyboard(rows ...[]KeyboardButton) [][]KeyboardButton {
	var keyboard [][]KeyboardButton
	keyboard = append(keyboard, rows...)
	return keyboard
}

type HolidayUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type WeatherUpdate struct {
	Name    string `json:"name"`
	Main    Main   `json:"main"`
	Wind    Wind   `json:"wind"`
	Weather []Weather
}

type Weather struct {
	Id          int    `json:"id,omitempty"` // 800,
	Main        string `json:"main"`         // "Clear",
	Description string `json:"description"`  // "clear sky",
	Icon        string `json:"icon"`         // "01d"

}

type Main struct {
	Temp      float64 `json:"temp"`       // 36.32,
	FeelsLike float64 `json:"feels_like"` // 33.72,
	TempMin   float64 `json:"temp_min"`   // 36.32,
	TempMax   float64 `json:"temp_max"`   // 36.32,
	Pressure  int     `json:"pressure"`   // 1008,
	Humidity  int     `json:"humidity"`   // 13,
	SeaLevel  int     `json:"sea_level"`  // 1008,
	GrndLevel int     `json:"grnd_level"` // 991
}

type Wind struct {
	Speed float64 `json:"speed"` // 7.33,
	Deg   int     `json:"deg"`   // 4,
	Gust  float64 `json:"gust"`  // 9.47
}
