module git.dbyte.xyz/distro/gerry

go 1.22.2

replace github.com/traefik/yaegi => github.com/distrobyte/yaegi v0.16.1-modules-fix

require (
	github.com/bwmarrin/discordgo v0.28.1
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/joho/godotenv v1.5.1
	github.com/traefik/yaegi v0.16.1
	gopkg.in/fsnotify.v1 v1.4.7
)

require (
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
)
