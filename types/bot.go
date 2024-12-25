package types

import "github.com/Anti-Raid/corelib_go/silverpelt"

type BotState struct {
	Commands           []silverpelt.CanonicalCommand      `json:"commands"`
	Settings           []silverpelt.CanonicalConfigOption `json:"settings"`
	CommandPermissions map[string][]string                `json:"command_permissions"`
}
