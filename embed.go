package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func GetEmbed(record *Book) (embed *discordgo.MessageEmbed) {
	//fmt.Println("Creating embed for: ", record.Title)

	/*filesize, err := strconv.Atoi(record.Filesize)
	if err != nil {
		fmt.Println(err)
		filesize = 0
	}*/

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
		//	"**Filesize:** " + prettyByteSize(filesize) + "\n"
		"**Filesize:** " + record.Filesize + "\n"

	cover := discordgo.MessageEmbedThumbnail{URL: record.CoverURL}
	footer := discordgo.MessageEmbedFooter{
		Text:         "Search results powered by libgen",
		IconURL:      "https://i.imgur.com/C4RNeln.png",
		ProxyIconURL: "",
	}
	// Create our discordgo embed
	embed = &discordgo.MessageEmbed{
		URL:         fmt.Sprintf("https://libgen.gs/edition.php?id=%s", record.ID),
		Title:       "**" + record.Title + "**",
		Description: description,
		Thumbnail:   &cover,
		Footer:      &footer,
	}

	return embed
}
