package csapi_tests

import (
	"testing"
	"time"

	"github.com/matrix-org/complement/internal/b"
	"github.com/matrix-org/complement/internal/client"
)

var dataType = "connect.multiroom.location"
var dataMrd = &client.FakeMrd{Foo: "bar"}

func TestMultiRoom(t *testing.T) {
	deployment := Deploy(t, b.BlueprintMultiRoom)
	defer deployment.Destroy(t)
	alice := deployment.Client(t, "hs1", "@alice:hs1")
	bob := deployment.Client(t, "hs1", "@bob:hs1")
	t.Run("multiroom data does not pass to sender when visibility is off", func(t *testing.T) {
		alice.SendMultiRoom(t, dataType, dataMrd)
		alice.MustSyncUntil(t, client.SyncReq{}, client.SyncNoMultiRoom(alice.UserID, dataType))

	})
	t.Run("multiroom data do not pass to others when visibility is off", func(t *testing.T) {
		alice.SendMultiRoom(t, dataType, dataMrd)
		bob.MustSyncCheck(t, client.SyncReq{}, client.SyncNoMultiRoom(alice.UserID, dataType))
	})
	t.Run("multiroom data pass to sender when visibility is on", func(t *testing.T) {
		roomID := alice.CreateRoom(t, map[string]interface{}{"preset": "public_chat"})
		alice.MustSyncUntil(t, client.SyncReq{}, client.SyncJoinedTo(alice.UserID, roomID))
		alice.SendMultiRoomVisibility(t, dataType, roomID, time.Now().Add(time.Hour))
		alice.SendMultiRoom(t, dataType, dataMrd)
		alice.MustSyncUntil(t, client.SyncReq{}, client.SyncMultiRoom(alice.UserID, dataType, dataMrd))
	})
	t.Run("multiroom data does not pass to users not joined to room", func(t *testing.T) {
		roomID := alice.CreateRoom(t, map[string]interface{}{"preset": "public_chat"})
		alice.MustSyncUntil(t, client.SyncReq{}, client.SyncJoinedTo(alice.UserID, roomID))
		alice.SendMultiRoomVisibility(t, dataType, roomID, time.Now().Add(time.Hour))
		alice.SendMultiRoom(t, dataType, dataMrd)
		alice.MustSyncUntil(t, client.SyncReq{}, client.SyncMultiRoom(alice.UserID, dataType, dataMrd))
		bob.MustSyncCheck(t, client.SyncReq{}, client.SyncNoMultiRoom(alice.UserID, dataType))
	})
	t.Run("multiroom data pass to other users when visibility is on and does not when visibility is off", func(t *testing.T) {
		roomID := alice.CreateRoom(t, map[string]interface{}{"preset": "public_chat"})
		alice.MustSyncUntil(t, client.SyncReq{}, client.SyncJoinedTo(alice.UserID, roomID))
		bob.MustDoFunc(t, "POST", []string{"_matrix", "client", "v3", "join", roomID})
		bob.MustSyncUntil(t, client.SyncReq{}, client.SyncJoinedTo(bob.UserID, roomID))
		alice.SendMultiRoomVisibility(t, dataType, roomID, time.Now().Add(time.Hour))
		alice.SendMultiRoom(t, dataType, dataMrd)
		alice.MustSyncUntil(t, client.SyncReq{}, client.SyncMultiRoom(alice.UserID, dataType, dataMrd))
		bob.MustSyncUntil(t, client.SyncReq{}, client.SyncMultiRoom(alice.UserID, dataType, dataMrd))
		alice.SendMultiRoomVisibilityOff(t, dataType, roomID)
		bob.MustSyncCheck(t, client.SyncReq{}, client.SyncNoMultiRoom(alice.UserID, dataType))
	})
}
