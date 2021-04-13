<h1>API для управления Яндекс станцией</h1>

`yapi` собран под ARM \
`GOOS=linux GOARCH=arm GOARM=6 go build -o yapi`

<h3>.env.local</h3>

LOGIN - логин от Яндекс аккаунта \
PASSWORD - пароль от Яндекс аккаунта \
STATION_ID - приложение яндекс -> устройства \
STATION_ADDR - ipadress станции

<h3>Установка</h3>

`cd /opt` \
`git clone https://github.com/ebuyan/yapi.git` \
`cp .env bin/.env.local` \
`cp yapi.service /etc/systemd/systemd` \
`systemctl daemon-reload` \
`systemctl start yapi.service` \
`systemctl enable yapi.service`

<h3>API</h3>

- Статус `GET <host>/`
- Перемотка `POST {
    "command": "rewind",
    "position" : 120
}`
- Продолжить `POST {
    "command": "play"
}`
- Пауза `POST {
    "command": "stop"
}`
- Следующий `POST {
    "command": "next"
}`
- Предыдущий `POST {
    "command": "prev"
}`
- Изменить громкость `POST {
    "command" : "setVolume",
	"volume" : 0.5
}`
- Выполнить команду `POST {
    "command" : "sendText",
	"text" : "Включи музыку"
}`
- Воспроизвести текст `POST {
    "command" : "sendText",
	"text" : "Повтори за мной 'Саша, иди собирай игрушки, а то выключу интернет в доме'"
}`
