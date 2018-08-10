package event

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"regexp"
)

// these are only used internally for marshalling/unmarshalling
// be sure to add it here and in event.UnmarshalJSON
const (
	eventRoomAliases           = "m.room.aliases"
	eventRoomCanonicalAlias    = "m.room.canonical_alias"
	eventRoomCreate            = "m.room.create"
	eventRoomJoinRules         = "m.room.join_rules"
	eventRoomMember            = "m.room.member"
	eventRoomPowerLevels       = "m.room.power_levels"
	eventRoomRedaction         = "m.room.redaction"
	eventRoomMessage           = "m.room.message"
	eventRoomMessageFeedback   = "m.room.message.feedback"
	eventRoomName              = "m.room.name"
	eventRoomTopic             = "m.room.topic"
	eventRoomAvatar            = "m.room.avatar"
	eventRoomPinnedEvents      = "m.room.pinned_events"
	eventCallInvite            = "m.call.invite"
	eventCallCandidates        = "m.call.candidates"
	eventCallAnswer            = "m.call.answer"
	eventCallHangup            = "m.call.hangup"
	eventTyping                = "m.typing"
	eventReceipt               = "m.receipt"
	eventPresence              = "m.presence"
	eventRoomHistoryVisibility = "m.room.history_visibility"
	eventRoomThirdPartyInvite  = "m.room.third_party_invite"
	eventRoomGuestAccess       = "m.room.guest_access"
	eventDirect                = "m.direct"
	messageText                = "m.text"
	messageEmote               = "m.emote"
	messageNotice              = "m.notice"
	messageImage               = "m.image"
	messageVideo               = "m.video"
	messageFile                = "m.file"
	messageLocation            = "m.location"
	messageAudio               = "m.audio"
)

const (
	formatOrgMatrixCustomHtml = "org.matrix.custom.html"
)

// Event represents a single Matrix event.
type Event struct {
	StateKey  *string     `json:"state_key,omitempty"` // The state key for the event. Only present on State Events.
	Sender    string      `json:"sender"`              // The user ID of the sender of the event
	Type      string      `json:"type"`                // The event type
	Timestamp int64       `json:"origin_server_ts"`    // The unix timestamp when this message was sent by the origin server
	ID        string      `json:"event_id"`            // The unique ID of this event
	RoomID    string      `json:"room_id"`             // The room the event was sent to. May be nil (e.g. for presence)
	Content   interface{} `json:"content"`             // The JSON content of the event.
	Redacts   string      `json:"redacts,omitempty"`   // The event ID that was redacted if a m.room.redaction event
}

// jsonEvent is used while unmarshalling to access Content as RawMessage
type jsonEvent struct {
	StateKey  *string         `json:"state_key,omitempty"` // The state key for the event. Only present on State Events.
	Sender    string          `json:"sender"`              // The user ID of the sender of the event
	Type      string          `json:"type"`                // The event type
	Timestamp int64           `json:"origin_server_ts"`    // The unix timestamp when this message was sent by the origin server
	ID        string          `json:"event_id"`            // The unique ID of this event
	RoomID    string          `json:"room_id"`             // The room the event was sent to. May be nil (e.g. for presence)
	Content   json.RawMessage `json:"content"`             // The JSON content of the event.
}

// UnmarshalJSON unmarshals JSON data into an Event.
func (e *Event) UnmarshalJSON(data []byte) error {
	je := jsonEvent{}
	if err := json.Unmarshal(data, &je); err != nil {
		return err
	}

	e.StateKey = je.StateKey
	e.Sender = je.Sender
	e.Type = je.Type
	e.Timestamp = je.Timestamp
	e.ID = je.ID
	e.RoomID = je.RoomID

	// unmarshal the content into the matching type
	switch je.Type {
	case eventRoomAliases:
		x := RoomAliases{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventRoomCanonicalAlias:
		x := RoomCanonicalAlias{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventRoomCreate:
		x := RoomCreate{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventRoomJoinRules:
		x := RoomJoinRules{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventRoomMember:
		x := RoomMember{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventRoomPowerLevels:
		x := RoomPowerLevels{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventRoomRedaction:
		x := RoomRedaction{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventRoomMessage:
		// first unmarshal the message contents to a map so we can access
		// msgtype.
		messageMap := make(map[string]interface{})
		err := json.Unmarshal(je.Content, &messageMap)
		if err != nil {
			log.Fatal(err)
		}
		if x, ok := messageMap["msgtype"]; ok {
			if msgType, ok := x.(string); ok {
				switch msgType {
				case "m.text":
					x := TextMessage{}
					if err := json.Unmarshal(je.Content, &x); err != nil {
						return err
					}
					e.Content = x
				case "m.emote":
					x := EmoteMessage{}
					if err := json.Unmarshal(je.Content, &x); err != nil {
						return err
					}
					e.Content = x
				case "m.notice":
					x := NoticeMessage{}
					if err := json.Unmarshal(je.Content, &x); err != nil {
						return err
					}
					e.Content = x
				case "m.image":
					x := NoticeMessage{}
					if err := json.Unmarshal(je.Content, &x); err != nil {
						return err
					}
					e.Content = x
				case "m.file":
					x := FileMessage{}
					if err := json.Unmarshal(je.Content, &x); err != nil {
						return err
					}
					e.Content = x
				case "m.location":
					x := LocationMessage{}
					if err := json.Unmarshal(je.Content, &x); err != nil {
						return err
					}
					e.Content = x
				case "m.video":
					x := VideoMessage{}
					if err := json.Unmarshal(je.Content, &x); err != nil {
						return err
					}
					e.Content = x
				case "m.audio":
					x := AudioMessage{}
					if err := json.Unmarshal(je.Content, &x); err != nil {
						return err
					}
					e.Content = x
				default:
					return fmt.Errorf("unknown msgtype: %v", msgType)
				}
			}
		}
	case eventRoomMessageFeedback:
		x := RoomMessageFeedback{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventRoomName:
		x := RoomName{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventRoomTopic:
		x := RoomTopic{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventRoomAvatar:
		x := RoomAvatar{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventRoomPinnedEvents:
		x := RoomPinnedEvents{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventCallInvite:
		x := CallInvite{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventCallCandidates:
		x := CallCandidates{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventCallAnswer:
		x := CallAnswer{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventCallHangup:
		x := CallHangup{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventTyping:
		x := Typing{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventReceipt:
		x := Receipt{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventPresence:
		x := Presence{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventRoomHistoryVisibility:
		x := RoomHistoryVisibility{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventRoomThirdPartyInvite:
		x := RoomThirdPartyInvite{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventRoomGuestAccess:
		x := RoomAliases{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	case eventDirect:
		x := RoomAliases{}
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	default:
		x := make(map[string]interface{})
		if err := json.Unmarshal(je.Content, &x); err != nil {
			return err
		}
		e.Content = x
	}

	return nil
}

// RoomAliases is the Content of a "m.room.alias" message.
type RoomAliases struct {
	Aliases []string `json:"aliases"`
}

// RoomCanonicalAlias is the Content of a "m.room.canonical_alias" message.
type RoomCanonicalAlias struct {
	Alias string `json:"alias"`
}

// RoomCreate is the Content of a "m.room.create" message.
type RoomCreate struct {
	Creator  string `json:"creator"`
	Federate bool   `json:"m.federate,omitempty"`
}

// RoomJoinRules is the Content if a "m.room.join_rules" message.
type RoomJoinRules struct {
	JoinRule string `json:"join_rule"`
}

type RoomMember struct {
	AvatarURL        string `json:"avatar_url,omitempty"`
	Displayname      string `json:"displayname,omitempty"`
	Membership       string `json:"membership"`
	IsDirect         string `json:"is_direct"`
	ThirdPartyInvite bool   `json:"third_party_invite"`
}

type RoomPowerLevels struct {
	Ban           int            `json:"ban,omitempty"`
	Events        map[string]int `json:"events,omitempty"`
	EventsDefault int            `json:"events_default,omitempty"`
	Invite        int            `json:"invite,omitempty"`
	Kick          int            `json:"kick,omitempty"`
	Redact        int            `json:"redact,omitempty"`
	StateDefault  int            `json:"state_default,omitempty"`
	Users         map[string]int `json:"users,omitempty"`
	UsersDefault  int            `json:"users_default,omitempty"`
}

type RoomRedaction struct {
	Reason string `json:"reason"`
}

type RoomMessageFeedback struct {
	TargetEventId string `json:"target_event_id"`
	Type          string `json:"type"`
}

type RoomName struct {
	Name string `json:"name"`
}

type RoomTopic struct {
	Topic string `json:"topic"`
}

type RoomAvatar struct {
	Info ImageInfo `json:"info"`
	URL  string    `json:"url"`
}

type RoomPinnedEvents struct {
	Pinned []string `json:"pinned"`
}

type RoomHistoryVisibility struct {
	HistoryVisibility string `json:"history_visibility"`
}

type RoomThirdPartyInvite struct {
	Displayname    string `json:"display_name"`
	KeyValidityURL string `json:"key_validity_url"`
	PublicKey      string `json:"public_key"`
	PublicKeys     []struct {
		KeyValidityURL string `json:"key_validity_url,omitempty"`
		PublicKey      string `json:"public_key"`
	} `json:"public_keys,omitempty"`
}

type RoomGuestAccess struct {
	GuestAccess string `json:"guest_access"`
}

type CallInvite struct {
	CallID string `json:"call_id"`
	Offer  struct {
		Type string `json:"type"`
		SDP  string `json:"sdp"`
	} `json:"offer"`
	Version  int `json:"version"`
	Lifetime int `json:"lifetime"`
}

type CallCandidates struct {
	CallID     string `json:"call_id"`
	Candidates []struct {
		SdpMid        string `json:"sdpMid"`
		SdpMLineIndex int    `json:"sdpMLineIndex"`
		Candidate     string `json:"candidate"`
	} `json:"candidates"`
	Version int `json:"version"`
}

type CallAnswer struct {
	CallID string `json:"call_id"`
	Answer struct {
		Type string `json:"type"`
		Sdp  string `json:"sdp"`
	} `json:"answer"`
	Version int `json:"version"`
}

type CallHangup struct {
	CallID  string `json:"call_id"`
	Version int    `json:"version"`
}

type Typing struct {
	UserIDs []string `json:"user_ids"`
}

type Receipt map[string]struct {
	MRead map[string]struct {
		Ts int `json:"ts"`
	} `json:"m.read"`
}

type Presence struct {
	AvatarURL       string `json:"avatar_url"`
	Displayname     string `json:"displayname"`
	LastActiveAgo   int    `json:"last_active_ago"`
	Presence        string `json:"presence"`
	CurrentlyActive bool   `json:"currently_active"`
	UserID          string `json:"user_id"`
}

//type Direct map[string][]string

// TextMessage is the contents of a Matrix formated message event.
type TextMessage struct {
	Body          string `json:"body"`
	MsgType       string `json:"msgtype"`
	Format        string `json:"format"`
	FormattedBody string `json:"formatted_body"`
}

var htmlRegex = regexp.MustCompile("<[^<]+?>")

func (m TextMessage) body() string {
	switch m.Format {
	case formatOrgMatrixCustomHtml:
		return html.UnescapeString(htmlRegex.ReplaceAllLiteralString(m.FormattedBody, ""))
	default:
		return m.FormattedBody
	}
}

func (m *TextMessage) MarshalJSON() ([]byte, error) {
	m.MsgType = messageText

	if m.Format != "" {
		m.Body = m.body()
	}

	return json.Marshal(m)
}

type EmoteMessage struct {
	TextMessage
}

func (m *EmoteMessage) MarshalJSON() ([]byte, error) {
	m.MsgType = messageEmote

	if m.Format != "" {
		m.Body = m.TextMessage.body()
	}

	return json.Marshal(m)
}

type NoticeMessage struct {
	Body    string `json:"body"`
	MsgType string `json:"msgtype"`
}

func (m *NoticeMessage) MarshalJSON() ([]byte, error) {
	m.MsgType = messageNotice
	return json.Marshal(m)
}

// ImageInfo contains info about an image - http://matrix.org/docs/spec/client_server/r0.2.0.html#m-image
type ImageInfo struct {
	Height   uint   `json:"h,omitempty"`
	Width    uint   `json:"w,omitempty"`
	Mimetype string `json:"mimetype,omitempty"`
	Size     uint   `json:"size,omitempty"`
}

// VideoInfo contains info about a video - http://matrix.org/docs/spec/client_server/r0.2.0.html#m-video
type VideoInfo struct {
	Mimetype      string    `json:"mimetype,omitempty"`
	ThumbnailInfo ImageInfo `json:"thumbnail_info"`
	ThumbnailURL  string    `json:"thumbnail_url,omitempty"`
	Height        uint      `json:"h,omitempty"`
	Width         uint      `json:"w,omitempty"`
	Duration      uint      `json:"duration,omitempty"`
	Size          uint      `json:"size,omitempty"`
}

// VideoMessage is an m.video  - http://matrix.org/docs/spec/client_server/r0.2.0.html#m-video
type VideoMessage struct {
	Body    string    `json:"body"`
	MsgType string    `json:"msgtype"`
	URL     string    `json:"url"`
	Info    VideoInfo `json:"info"`
}

func (m *VideoMessage) MarshalJSON() ([]byte, error) {
	m.MsgType = messageVideo
	return json.Marshal(m)
}

// ImageMessage is an m.image event
type ImageMessage struct {
	Body    string    `json:"body"`
	MsgType string    `json:"msgtype"`
	URL     string    `json:"url"`
	Info    ImageInfo `json:"info"`
}

func (m *ImageMessage) MarshalJSON() ([]byte, error) {
	m.MsgType = messageImage
	return json.Marshal(m)
}

type FileInfo struct {
	Mimetype      string    `json:"mimetype,omitempty"`
	Size          int       `json:"size"`
	ThumbnailURL  string    `json:"thumbnail_url"`
	ThumbnailInfo ImageInfo `json:"thumbnail_info"`
}

type FileMessage struct {
	Body     string   `json:"body"`
	MsgType  string   `json:"msgtype"`
	Filename string   `json:"filename"`
	Info     FileInfo `json:"info"`
	URL      string   `json:"url"`
}

func (m *FileMessage) MarshalJSON() ([]byte, error) {
	m.MsgType = messageFile
	return json.Marshal(m)
}

type LocationInfo struct {
	ThumbnailURL  string    `json:"thumbnail_url"`
	ThumbnailInfo ImageInfo `json:"thumbnail_info"`
}

type LocationMessage struct {
	Body    string       `json:"body"`
	MsgType string       `json:"msgtype"`
	GeoURI  string       `json:"geo_uri"`
	Info    LocationInfo `json:"location_info"`
}

func (m *LocationMessage) MarshalJSON() ([]byte, error) {
	m.MsgType = messageLocation
	return json.Marshal(m)
}

type AudioInfo struct {
	Duration int    `json:"duration"`
	Mimetype string `json:"mimetype"`
	Size     int    `json:"size"`
}

type AudioMessage struct {
	Body    string    `json:"body"`
	MsgType string    `json:"msgtype"`
	Info    AudioInfo `json:"audio_info"`
	URL     string    `json:"url"`
}

func (m *AudioMessage) MarshalJSON() ([]byte, error) {
	m.MsgType = messageAudio
	return json.Marshal(m)
}
