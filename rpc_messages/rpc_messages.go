package rpc_messages

import (
	"github.com/Anti-Raid/api/types"
	"github.com/bwmarrin/discordgo"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type BaseGuildUserInfo struct {
	OwnerID   string                              `json:"owner_id"`
	Name      string                              `json:"name"`
	Icon      *string                             `json:"icon"`
	Roles     []types.SerenityRole                `json:"roles"`
	UserRoles []string                            `json:"user_roles"`
	BotRoles  []string                            `json:"bot_roles"`
	Channels  []types.GuildChannelWithPermissions `json:"channels"`
}

type CheckCommandPermission struct {
	Result *string `json:"result"`
}

type SettingsOperationRequest struct {
	Fields  orderedmap.OrderedMap[string, any] `json:"fields"`
	Op      types.CanonicalOperationType       `json:"op"`
	Setting string                             `json:"setting"`
}

type CanonicalSettingsResult struct {
	Ok *struct {
		Fields []orderedmap.OrderedMap[string, any] `json:"fields"`
	} `json:"Ok"`
	Err *struct {
		Error string `json:"error"`
	} `json:"Err"`
}

type JobserverSpawn struct {
	Name    string                 `json:"name"`
	Data    map[string]interface{} `json:"data"`
	Create  bool                   `json:"create"`
	Execute bool                   `json:"execute"`

	// If create is false, then `id`` must be set
	ID string `json:"id"`

	// The User ID who initiated the action
	UserID string `json:"user_id"`
}

type JobserverSpawnResponse struct {
	ID string `json:"id"`
}

/*
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ExecuteTemplateRequest {
    pub args: serde_json::Value,
    pub template: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum ExecuteTemplateResponse {
    Ok { result: Option<serde_json::Value> },
    ExecErr { error: String },
    PermissionError { res: PermissionResult },
}
*/

type ExecuteTemplateRequest struct {
	Args     any    `json:"args"`
	Template string `json:"template"`
}

type ExecuteTemplateResponse struct {
	Ok *struct {
		Result any `json:"result"`
	} `json:"Ok,omitempty"`
	ExecErr *struct {
		Error string `json:"error"`
	} `json:"ExecErr,omitempty"`
}

type CheckUserHasKittycatPermissionsRequest struct {
	Perm string `json:"perm"`
}

type TWState struct {
	Commands []discordgo.ApplicationCommand `json:"commands"`
	Settings []types.CanonicalConfigOption  `json:"settings"`
}
