package main

import "github.com/bwmarrin/discordgo"

const (
	DISCORD_ALLOW = 1
	DISCORD_DENY  = 0
	DISCORD_ROLE  = "@everyone"
)

func createDiscordTextChannel(channel_name string, discord *discordgo.Session, guildID string, role_id string) (string, error) {
	channel, err := discord.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     channel_name,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: getEnv("CATEGORY_ID"),
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:   guildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel | discordgo.PermissionReadMessageHistory,
			},
			{
				ID:    discord.State.User.ID,
				Type:  discordgo.PermissionOverwriteTypeMember,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionReadMessageHistory,
			},
			{
				ID:    role_id,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionReadMessageHistory,
			},
		},
	})

	if err != nil {
		return "", err
	}
	return channel.ID, nil
}

func createDiscordRole(role_name string, discord *discordgo.Session, guildID string) (*discordgo.Role, error) {
	role, err := discord.GuildRoleCreate(guildID, &discordgo.RoleParams{
		Name: role_name,
	})
	if err != nil {
		return nil, err
	}

	return role, nil
}

// func setDiscordPermissions(discord *discordgo.Session, channelID string, role string, allow discordgo.PermissionOverwriteType, deny discordgo.PermissionOverwriteType) error {
// 	err := discord.ChannelPermissionSet(channelID, role, allow, discordgo.PermissionViewChannel, discordgo.PermissionReadMessageHistory)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
