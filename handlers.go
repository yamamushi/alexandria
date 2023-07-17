package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

var (
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){

		"library": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			output := SearchLibgen(searchInput)

			if len(output) == 0 {
				err := s.InteractionResponseDelete(i.Interaction)
				if err != nil {
					fmt.Println("could not delete interaction" + err.Error())
				}
				_, _ = s.ChannelMessageSend(i.ChannelID, "Sorry "+i.Member.Mention()+" no results were found for **"+searchInput+"**")
				return
			}

			embed := GetEmbed(output[0])

			var searchResultsDropdown []discordgo.SelectMenuOption
			adb := AlexandriaDB{}
			err := adb.OpenDB()
			if err != nil {
				fmt.Println(err)
				return // Exit if we can't open the database
			}
			for _, result := range output {
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

				filesize, err := strconv.Atoi(result.Filesize)
				if err != nil {
					fmt.Println(err)
					filesize = 0
				}

				searchResultsDropdown = append(searchResultsDropdown, discordgo.SelectMenuOption{

					Label: result.Title,
					// As with components, this things must have their own unique "id" to identify which is which.
					// In this case such id is Value field.
					Value: result.Md5,
					// You can also make it a default option, but in this case we won't.
					Default:     false,
					Description: fmt.Sprintf("%s - %s - %s - %s", result.Author, result.Year, result.Extension, prettyByteSize(filesize)),
				})
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
				fmt.Println("could not delete interaction" + err.Error())
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

	for _, value := range i.Interaction.MessageComponentData().Values {
		// if i.Message.Interaction.Member.User.ID != i.Interaction.Member.User.ID // Not working as intended
		//fmt.Println("Value: ", value)
		adb := AlexandriaDB{}
		adb.OpenDB()
		defer adb.CloseDB()
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
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Ignore type for now, they will be discussed in "responses"
			Type: discordgo.InteractionResponseUpdateMessage,
		})
	}
}
