package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"git.dbyte.xyz/distro/gerry/bot"
	"git.dbyte.xyz/distro/gerry/commands"
	"git.dbyte.xyz/distro/gerry/handlers"
	"git.dbyte.xyz/distro/gerry/utils"
	"github.com/disgoorg/disgo"
	discord "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	mumble "layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleutil"
	_ "layeh.com/gumble/opus"
)

func main() {
	path := flag.String("config", "config.yml", "path to the config file")
	flag.Parse()

	cfg, err := bot.ReadConfig(*path)
	if err != nil {
		log.Fatalln("failed to read config", err)
	}

	log.Printf("starting gerrybot")

	b := &bot.Bot{
		Cfg: cfg,
	}

	b.MarkovChainImpl = utils.NewMarkovChain()
	b.MarkovChainImpl.LoadModel("./markovChain.json")
	b.MarkovChainImpl.LoadModel("./word_freq_map.json")

	cmds := commands.Commands{Bot: b}
	b.CommandList = make(map[string]bot.Command)

	cmds.LoadCommands()

	hdlr := &handlers.Handlers{Bot: b}

	if b.Cfg.Discord.Token != "" {
		if b.DiscordClient, err = disgo.New(cfg.Discord.Token,
			discord.WithGatewayConfigOpts(
				gateway.WithIntents(gateway.IntentsAll),
			),
			discord.WithCacheConfigOpts(
				cache.WithCaches(cache.FlagsAll),
			),
			discord.WithEventListenerFunc(hdlr.OnDiscordReady),
			discord.WithEventListenerFunc(hdlr.OnDiscordMessageCreate),
		); err != nil {
			log.Fatalln("failed to create disgo client", err)
		}
	}

	if b.Cfg.Mumble.Address != "" {
		b.MumbleConfig = mumble.NewConfig()
		b.MumbleConfig.Username = b.Cfg.Mumble.Username

		b.MumbleConfig.Attach(gumbleutil.Listener{TextMessage: hdlr.OnMumbleTextMessage})

		b.MumbleConfig.Attach(gumbleutil.Listener{Connect: hdlr.OnMumbleReady})

		b.MumbleConfig.Attach(gumbleutil.Listener{ChannelChange: hdlr.OnMumbleChannelChange})

		b.MumbleConfig.Attach(gumbleutil.Listener{UserChange: hdlr.OnMumbleUserChange})

		b.MumbleConfig.Attach(gumbleutil.Listener{Disconnect: hdlr.OnMumbleDisconnect})
	}

	if err = b.Start(); err != nil {
		log.Fatalln("failed to start bot:", err)
	}
	defer b.Stop()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
