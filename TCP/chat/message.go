package chat

import "encoding/json"

type message struct {
	Content string // the data of the message
	RoomID  int    // the room the user is present in
	UserIP  string // the IP of the person sending the message
}

func (m *message) Marshal() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (m *message) Unmarshal(bytes []byte) error {
	err := json.Unmarshal(bytes, &m)
	return err
}
