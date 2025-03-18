module s-belichenko/ilovaiskaya2-bot

go 1.21.0

require (
	github.com/go-test/deep v1.1.1
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/sheeiavellie/go-yandexgpt v0.0.0-00010101000000-000000000000
	gopkg.in/telebot.v4 v4.0.0-beta.4
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/sheeiavellie/go-yandexgpt => github.com/s-belichenko/go-yandexgpt v1.7.1
