package rpc_messages

import (
	"github.com/Anti-Raid/corelib_go/ext_types"
	"github.com/Anti-Raid/corelib_go/silverpelt"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type BaseGuildUserInfo struct {
	OwnerID   string                                  `json:"owner_id"`
	Name      string                                  `json:"name"`
	Icon      *string                                 `json:"icon"`
	Roles     []ext_types.SerenityRole                `json:"roles"`
	UserRoles []string                                `json:"user_roles"`
	BotRoles  []string                                `json:"bot_roles"`
	Channels  []ext_types.GuildChannelWithPermissions `json:"channels"`
}

type CheckCommandPermission struct {
	Result *string `json:"result"`
}

type CheckCommandPermissionRequest struct {
	Command string `json:"command"`
}

type SettingsOperationRequest struct {
	Fields  orderedmap.OrderedMap[string, any] `json:"fields"`
	Op      silverpelt.CanonicalOperationType  `json:"op"`
	Setting string                             `json:"setting"`
}

type CanonicalSettingsResult struct {
	Ok *struct {
		Fields []orderedmap.OrderedMap[string, any] `json:"fields"`
	} `json:"Ok"`
	Err *struct {
		Error silverpelt.CanonicalSettingsError `json:"error"`
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

/**
#[derive(Serialize, Deserialize)]
pub struct BotState {
    pub commands: Vec<crate::botlib::canonical::CanonicalCommand>,
    pub settings: Vec<ar_settings::types::Setting>,
    pub command_permissions: crate::botlib::CommandPermissionMetadata,
}
*/

type BotState struct {
	Commands           []silverpelt.CanonicalCommand      `json:"commands"`
	Settings           []silverpelt.CanonicalConfigOption `json:"settings"`
	CommandPermissions map[string][]string                `json:"command_permissions"`
}
