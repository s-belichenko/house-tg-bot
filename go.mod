module s-belichenko/ilovaiskaya2-bot

go 1.21.0

require (
	github.com/go-test/deep v1.1.1
	github.com/joho/godotenv v1.5.1
	github.com/rs/zerolog v1.33.0
	github.com/sheeiavellie/go-yandexgpt v0.0.0-00010101000000-000000000000
	gopkg.in/telebot.v4 v4.0.0-beta.4
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
)

replace github.com/sheeiavellie/go-yandexgpt => github.com/s-belichenko/go-yandexgpt v1.7.1
