package chat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	expected := `{"Content":"test","RoomID":500,"UserIP":"127.0.0.1:8080"}`
	testMessage := message{
		Content: "test",
		RoomID:  500,
		UserIP:  "127.0.0.1:8080",
	}

	bytes, err := testMessage.Marshal()
	require.NoError(t, err)
	require.Equal(t, string(bytes), expected)
}

func TestUnmarshal(t *testing.T) {
	testData := []byte(`{"Content":"test","RoomID":500,"UserIP":"127.0.0.1:8080"}`)
	expected := message{
		Content: "test",
		RoomID:  500,
		UserIP:  "127.0.0.1:8080",
	}

	err := expected.Unmarshal(testData)
	require.NoError(t, err)
	require.Equal(t, expected, expected)
}
