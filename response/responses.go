package response

import "github.com/rbns/gomatrix/event"

// Error is the standard JSON error response from Homeservers. It also implements the Golang "error" interface.
// See http://matrix.org/docs/spec/client_server/r0.2.0.html#api-standards
type Error struct {
	ErrCode string `json:"errcode"`
	Err     string `json:"error"`
}

// Error returns the errcode and error message.
func (e Error) Error() string {
	return e.ErrCode + ": " + e.Err
}

// CreateFilter is the JSON response for http://matrix.org/docs/spec/client_server/r0.2.0.html#post-matrix-client-r0-user-userid-filter
type CreateFilter struct {
	FilterID string `json:"filter_id"`
}

// Versions is the JSON response for http://matrix.org/docs/spec/client_server/r0.2.0.html#get-matrix-client-versions
type Versions struct {
	Versions []string `json:"versions"`
}

// JoinRoom is the JSON response for http://matrix.org/docs/spec/client_server/r0.2.0.html#post-matrix-client-r0-rooms-roomid-join
type JoinRoom struct {
	RoomID string `json:"room_id"`
}

// LeaveRoom is the JSON response for http://matrix.org/docs/spec/client_server/r0.2.0.html#post-matrix-client-r0-rooms-roomid-leave
type LeaveRoom struct{}

// ForgetRoom is the JSON response for http://matrix.org/docs/spec/client_server/r0.2.0.html#post-matrix-client-r0-rooms-roomid-forget
type ForgetRoom struct{}

// InviteUser is the JSON response for http://matrix.org/docs/spec/client_server/r0.2.0.html#post-matrix-client-r0-rooms-roomid-invite
type InviteUser struct{}

// KickUser is the JSON response for http://matrix.org/docs/spec/client_server/r0.2.0.html#post-matrix-client-r0-rooms-roomid-kick
type KickUser struct{}

// BanUser is the JSON response for http://matrix.org/docs/spec/client_server/r0.2.0.html#post-matrix-client-r0-rooms-roomid-ban
type BanUser struct{}

// UnbanUser is the JSON response for http://matrix.org/docs/spec/client_server/r0.2.0.html#post-matrix-client-r0-rooms-roomid-unban
type UnbanUser struct{}

// Typing is the JSON response for https://matrix.org/docs/spec/client_server/r0.2.0.html#put-matrix-client-r0-rooms-roomid-typing-userid
type Typing struct{}

// JoinedRooms is the JSON response for TODO-SPEC https://github.com/matrix-org/synapse/pull/1680
type JoinedRooms struct {
	JoinedRooms []string `json:"joined_rooms"`
}

// JoinedMembers is the JSON response for TODO-SPEC https://github.com/matrix-org/synapse/pull/1680
type JoinedMembers struct {
	Joined map[string]struct {
		DisplayName *string `json:"display_name"`
		AvatarURL   *string `json:"avatar_url"`
	} `json:"joined"`
}

// Messages is the JSON response for https://matrix.org/docs/spec/client_server/r0.2.0.html#get-matrix-client-r0-rooms-roomid-messages
type Messages struct {
	Start string        `json:"start"`
	Chunk []event.Event `json:"chunk"`
	End   string        `json:"end"`
}

// SendEvent is the JSON response for http://matrix.org/docs/spec/client_server/r0.2.0.html#put-matrix-client-r0-rooms-roomid-send-eventtype-txnid
type SendEvent struct {
	EventID string `json:"event_id"`
}

// MediaUpload is the JSON response for http://matrix.org/docs/spec/client_server/r0.2.0.html#post-matrix-media-r0-upload
type MediaUpload struct {
	ContentURI string `json:"content_uri"`
}

// UserInteractive is the JSON response for https://matrix.org/docs/spec/client_server/r0.2.0.html#user-interactive-authentication-api
type UserInteractive struct {
	Flows []struct {
		Stages []string `json:"stages"`
	} `json:"flows"`
	Params    map[string]interface{} `json:"params"`
	Session   string                 `json:"string"`
	Completed []string               `json:"completed"`
	ErrCode   string                 `json:"errcode"`
	Error     string                 `json:"error"`
}

// HasSingleStageFlow returns true if there exists at least 1 Flow with a single stage of stageName.
func (r UserInteractive) HasSingleStageFlow(stageName string) bool {
	for _, f := range r.Flows {
		if len(f.Stages) == 1 && f.Stages[0] == stageName {
			return true
		}
	}
	return false
}

// UserDisplayName is the JSON response for https://matrix.org/docs/spec/client_server/r0.2.0.html#get-matrix-client-r0-profile-userid-displayname
type UserDisplayName struct {
	DisplayName string `json:"displayname"`
}

// Register is the JSON response for http://matrix.org/docs/spec/client_server/r0.2.0.html#post-matrix-client-r0-register
type Register struct {
	AccessToken  string `json:"access_token"`
	DeviceID     string `json:"device_id"`
	HomeServer   string `json:"home_server"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
}

// Login is the JSON response for http://matrix.org/docs/spec/client_server/r0.2.0.html#post-matrix-client-r0-login
type Login struct {
	AccessToken string `json:"access_token"`
	DeviceID    string `json:"device_id"`
	HomeServer  string `json:"home_server"`
	UserID      string `json:"user_id"`
}

// Logout is the JSON response for http://matrix.org/docs/spec/client_server/r0.2.0.html#post-matrix-client-r0-logout
type Logout struct{}

// CreateRoom is the JSON response for https://matrix.org/docs/spec/client_server/r0.2.0.html#post-matrix-client-r0-createroom
type CreateRoom struct {
	RoomID string `json:"room_id"`
}

// Sync is the JSON response for http://matrix.org/docs/spec/client_server/r0.2.0.html#get-matrix-client-r0-sync
type Sync struct {
	NextBatch   string `json:"next_batch"`
	AccountData struct {
		Events []event.Event `json:"events"`
	} `json:"account_data"`
	Presence struct {
		Events []event.Event `json:"events"`
	} `json:"presence"`
	Rooms struct {
		Leave map[string]struct {
			State struct {
				Events []event.Event `json:"events"`
			} `json:"state"`
			Timeline struct {
				Events    []event.Event `json:"events"`
				Limited   bool          `json:"limited"`
				PrevBatch string        `json:"prev_batch"`
			} `json:"timeline"`
		} `json:"leave"`
		Join map[string]struct {
			State struct {
				Events []event.Event `json:"events"`
			} `json:"state"`
			Timeline struct {
				Events    []event.Event `json:"events"`
				Limited   bool          `json:"limited"`
				PrevBatch string        `json:"prev_batch"`
			} `json:"timeline"`
			Ephemeral struct {
				Events []event.Event `json:"events"`
			} `json:"ephemeral"`
			UnreadNotifications struct {
				HighlightCount    int `json:"highlight_count"`
				NotificationCount int `json:"notification_count"`
			}
		} `json:"join"`
		Invite map[string]struct {
			State struct {
				Events []event.Event
			} `json:"invite_state"`
		} `json:"invite"`
	} `json:"rooms"`
}

type TurnServer struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	TTL      int      `json:"ttl"`
	URIs     []string `json:"uris"`
}
