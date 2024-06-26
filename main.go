package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Variables used for command line parameters
var (
	ConfPath string
	Conf     Config
)

func init() {
	// Read our command line options
	flag.StringVar(&ConfPath, "c", "alexandria.conf", "Path to Config File")
	flag.Parse()

	_, err := os.Stat(ConfPath)
	if err != nil {
		log.Fatal("Config file is missing: ", ConfPath)
	}

	// Verify we can actually read our config file
	err = ReadConfig(ConfPath)
	if err != nil {
		log.Fatal("error reading config file at: ", ConfPath)
		return
	}

}

func main() {
	fmt.Println("\n\n|| Starting Alexandria Bot ||")
	log.SetOutput(io.Discard)

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Conf.DiscordConfig.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	defer func(dg *discordgo.Session) {
		err := dg.Close()
		if err != nil {
			panic(err)
		}
	}(dg)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Bot is up!")
	})

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			SelectSearchResult(s, i)
		}
	})

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands))
	for i, v := range Commands {
		for _, Guild := range dg.State.Guilds {
			cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, Guild.ID, v)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}
			registeredCommands[i] = cmd
		}
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Removing commands...")
	// // We need to fetch the commands, since deleting requires the command ID.
	// // We are doing this from the returned commands on line 375, because using
	// // this will delete all the commands, which might not be desirable, so we
	// // are deleting only the commands that we added.
	// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
	// if err != nil {
	// 	log.Fatalf("Could not fetch registered commands: %v", err)
	// }

	for _, Guild := range dg.State.Guilds {
		r1, _ := dg.ApplicationCommands(dg.State.User.ID, Guild.ID)
		for _, v := range r1 {
			err := dg.ApplicationCommandDelete(dg.State.User.ID, Guild.ID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	for _, v := range registeredCommands {
		for _, Guild := range dg.State.Guilds {

			err := dg.ApplicationCommandDelete(dg.State.User.ID, Guild.ID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	// Cleanly close down the Discord session.
	err = dg.Close()
	if err != nil {
		panic(err)
	}
}
