package bot

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"git.dbyte.xyz/distro/gerry/utils"
	"github.com/disgoorg/disgo"
	discord "github.com/disgoorg/disgo/bot"
	discordUtil "github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	mumble "layeh.com/gumble/gumble"
)

type Bot struct {
	Cfg             Config
	DiscordClient   discord.Client
	MumbleClient    *mumble.Client
	MumbleConfig    *mumble.Config
	CommandList     map[string]Command
	MarkovChainImpl *utils.MarkovChainImpl
	StopCh          chan struct{}
}

func (b *Bot) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if b.Cfg.Discord.Token != "" {
		log.Printf("connnecting to discord. disgo version: %s\n", disgo.Version)
		if err := b.DiscordClient.OpenGateway(ctx); err != nil {
			return err
		}
	}

	if b.Cfg.Mumble.Address != "" {
		log.Printf("connecting to mumble server %s:%v\n", b.Cfg.Mumble.Address, b.Cfg.Mumble.Port)

		var tlsConfig tls.Config
		if !b.Cfg.Mumble.TLS {
			tlsConfig.InsecureSkipVerify = true
		}

		var err error
		b.MumbleClient, err = mumble.DialWithDialer(new(net.Dialer),
			fmt.Sprintf("%s:%v", b.Cfg.Mumble.Address, b.Cfg.Mumble.Port),
			b.MumbleConfig,
			&tlsConfig,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	b.DiscordClient.Close(ctx)
}

func (b *Bot) SendMessage(msg Message, content string) error {
	// Implement sending message logic for Discord
	switch msg.GetPlatform() {
	case "discord":
		return sendDiscordMessage(b, snowflake.MustParse(msg.GetChannelID()), content)
	case "mumble":
		channelIDInt, err := strconv.ParseInt(msg.GetChannelID(), 10, 32)
		if err != nil {
			return err
		}
		return sendMumbleMessage(b, uint32(channelIDInt), content)
	default:
		return fmt.Errorf("unsupported platform: %s", msg.GetPlatform())
	}
}

func sendDiscordMessage(bot *Bot, channelID snowflake.ID, content string) error {
	// Implement sending message logic for Discord
	fmt.Printf("Sending Discord message to channel %s: %s\n", channelID, content)

	channel, err := bot.DiscordClient.Rest().GetChannel(channelID)
	if err != nil {
		return err
	}

	message := discordUtil.MessageCreate{
		Content: content,
	}

	_, err = bot.DiscordClient.Rest().CreateMessage(channel.ID(), message)
	if err != nil {
		return err
	}

	return nil
}

func sendMumbleMessage(bot *Bot, channelID uint32, content string) error {
	fmt.Printf("Sending Mumble message to channel %d: %s\n", channelID, content)

	channel := bot.MumbleClient.Channels[channelID]
	channel.Send(content, false)

	return nil
}
