package main

type Config struct {
	//Telegram
	Host  string `env:"HOST"`
	Token string `env:"TOKEN"`

	DbHost     string `env:"MONGODB_URI"`
	Collection string `env:"COLLECTION"`
	Db         string `env:"DB"`
	DbTimeout  int    `env:"DB_TIMEOUT"`

	//Weather API
	WeatherHost   string `env:"WEATHER_HOST"`
	WeatherApiKey string `env:"WEATHER_API_KEY"`

	//Logger parameters
	Infolevel int  `env:"INFOLEVEL"`
	Timestamp bool `env:"TIMESTAMP"`
	Caller    bool `env:"CALLER"`
	Pretty    bool `env:"PRETTY"`
}
