package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/ciehanski/libgen-cli/libgen"
	"strconv"
)

func GetEmbed(record *libgen.Book) (embed *discordgo.MessageEmbed) {
	filesize, err := strconv.Atoi(record.Filesize)
	if err != nil {
		fmt.Println(err)
		filesize = 0
	}

	description := "" +
		"**Title:** " + record.Title + "\n" +
		"**Author:** " + record.Author + "\n" +
		"**Pages:** " + record.Pages + "\n" +
		"**Year:** " + record.Year + "\n" +
		//"**MD5:** " + output[0].Md5 + "\n" +
		"**Publisher:** " + record.Publisher + "\n" +
		//"**Edition:** " + output[0].Edition + "\n" +
		"**Language:** " + record.Language + "\n" +
		"**Extension:** " + record.Extension + "\n" +
		"**Filesize:** " + prettyByteSize(filesize) + "\n"

	cover := discordgo.MessageEmbedThumbnail{URL: fmt.Sprintf("https://library.lol/covers/%s", record.CoverURL)}
	footer := discordgo.MessageEmbedFooter{
		Text:         "Search results powered by libgen",
		IconURL:      "https://i.imgur.com/C4RNeln.png",
		ProxyIconURL: "",
	}
	// Create our discordgo embed
	embed = &discordgo.MessageEmbed{
		URL:         fmt.Sprintf("https://library.lol/main/%s", record.Md5),
		Title:       "**" + record.Title + "**",
		Description: description,
		Thumbnail:   &cover,
		Footer:      &footer,
	}

	return embed
}
