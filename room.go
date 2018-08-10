package gomatrix

import "github.com/rbns/gomatrix/event"

// Room represents a single Matrix room.
type Room struct {
	ID    string
	State map[string]map[string]*event.Event
}

// UpdateState updates the room's current state with the given Event. This will clobber events based
// on the type/state_key combination.
func (room Room) UpdateState(e *event.Event) {
	_, exists := room.State[e.Type]
	if !exists {
		room.State[e.Type] = make(map[string]*event.Event)
	}
	room.State[e.Type][*e.StateKey] = e
}

// GetStateEvent returns the state event for the given type/state_key combo, or nil.
func (room Room) GetStateEvent(eventType string, stateKey string) *event.Event {
	stateEventMap, _ := room.State[eventType]
	event, _ := stateEventMap[stateKey]
	return event
}

// GetMembershipState returns the membership state of the given user ID in this room. If there is
// no entry for this member, 'leave' is returned for consistency with left users.
func (room Room) GetMembershipState(userID string) string {
	e := room.GetStateEvent("m.room.member", userID)
	if e != nil {
		if t, ok := e.Content.(*event.RoomMember); ok {
			return t.Membership
		}
	}
	return "leave"
}

// NewRoom creates a new Room with the given ID
func NewRoom(roomID string) *Room {
	// Init the State map and return a pointer to the Room
	return &Room{
		ID:    roomID,
		State: make(map[string]map[string]*event.Event),
	}
}
