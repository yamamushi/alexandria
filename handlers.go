package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var (
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){

		"library": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			//fmt.Println("Received command: ", i.ApplicationCommandData().Name)
			//fmt.Println("Starting Search")
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			searchInput := ""

			// Get the value from the option map.
			// When the option exists, ok = true
			if option, ok := optionMap["input"]; ok {
				searchInput = option.StringValue()
			}

			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, they will be discussed in "responses"
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Searching for **" + searchInput + "**... (limited to 10 results)",
				},
			})

			// Take our message and pass it into parser
			output, err := gsSearch(searchInput)

			if len(output) == 0 || err != nil {
				if err != nil {
					fmt.Println("error: " + err.Error())
				}
				fmt.Println("No results found for: ", searchInput)
				err := s.InteractionResponseDelete(i.Interaction)
				if err != nil {
					fmt.Println("could not delete interaction" + err.Error())
				}
				_, _ = s.ChannelMessageSend(i.ChannelID, "Sorry "+i.Member.Mention()+" no results were found for **"+searchInput+"**")
				return
			}

			// Create an embed for the first result
			embed := GetEmbed(output[0])

			var searchResultsDropdown []discordgo.SelectMenuOption
			adb := AlexandriaDB{}
			err = adb.OpenDB()
			if err != nil {
				fmt.Println(err)
				return // Exit if we can't open the database
			}

			var ids []string = make([]string, len(output))
			for count, result := range output {
				err = adb.StoreRecord(result)
				if err != nil {
					fmt.Println("could not store record: " + err.Error())
					return // Exit if we can't store the record
				}

				if len(result.Title) > 90 {
					result.Title = result.Title[:90]
				}
				if len(result.Author) > 70 {
					result.Author = result.Author[:70]
				}

				/*filesize, err := strconv.Atoi(result.Filesize)
				if err != nil {
					fmt.Println(err)
					filesize = 0
				}*/

				//fmt.Println("Creating dropdown option for: ", result.Title)

				if StringInSlice(result.ID, ids) {
					continue
				}
				searchResultsDropdown = append(searchResultsDropdown, discordgo.SelectMenuOption{

					Label: result.Title,
					// As with components, this things must have their own unique "id" to identify which is which.
					// In this case such id is Value field.
					Value: result.ID,
					// You can also make it a default option, but in this case we won't.
					Default: false,
					Description: truncate(fmt.Sprintf("%s - %s - %s - %s",
						result.Author, result.Year, result.Extension, result.Filesize), 24),
					Emoji: &discordgo.ComponentEmoji{Name: "📚"},
				})
				ids = append(ids, result.ID)
				if count == 10 {
					break
				}
			}
			err = adb.CloseDB()
			if err != nil {
				fmt.Println(err)
				return // Exit if we can't close the database
			}

			components := []discordgo.MessageComponent{discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						// Select menu, as other components, must have a customID, so we set it to this value.
						CustomID:    "select",
						Placeholder: "Choose a search result:",
						Options:     searchResultsDropdown,
					}},
			}}

			err = s.InteractionResponseDelete(i.Interaction)
			if err != nil {
				fmt.Println("could not delete response interaction: " + err.Error())
			}

			resultOutput := discordgo.MessageSend{
				Content:    i.Member.Mention() + "'s search results for **" + searchInput + "** (limited to 10 results):",
				Components: components,
				Embed:      embed,
			}

			_, err = s.ChannelMessageSendComplex(i.ChannelID, &resultOutput)
			if err != nil {
				fmt.Println(err)
			}
		},
	}
)

func SelectSearchResult(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//fmt.Println(i.Message.ID, i.Interaction.MessageComponentData().CustomID)
	adb := AlexandriaDB{}
	err := adb.OpenDB()
	if err != nil {
		return
	}
	defer func(adb *AlexandriaDB) {
		err := adb.CloseDB()
		if err != nil {
			fmt.Println(err)
		}
	}(&adb)

	for _, value := range i.Interaction.MessageComponentData().Values {
		// if i.Message.Interaction.Member.User.ID != i.Interaction.Member.User.ID // Not working as intended
		//fmt.Println("Value: ", value)
		record, err := adb.GetRecord(value)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = s.ChannelMessageEditEmbed(i.ChannelID, i.Message.ID, GetEmbed(record))
		if err != nil {
			fmt.Println("error updating embed: " + err.Error())
			return
		}
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Ignore type for now, they will be discussed in "responses"
			Type: discordgo.InteractionResponseUpdateMessage,
		})
		if err != nil {
			fmt.Println("error updating interaction: " + err.Error())
			return
		}
	}
}
