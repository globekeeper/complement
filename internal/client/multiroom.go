package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/matrix-org/complement/internal/b"
	"github.com/tidwall/gjson"
)

type FakeMrd struct {
	Foo string
}

func (c *CSAPI) SendMultiRoom(t *testing.T, dataType string, data interface{}) {
	t.Helper()
	paths := []string{"_matrix", "client", "v3", "multiroom", dataType}
	c.MustDo(t, "POST", paths, data)
}

func (c *CSAPI) SendMultiRoomVisibility(t *testing.T, dataType string, roomId string, expire time.Time) {
	t.Helper()
	c.SendEventSynced(t, roomId, b.Event{
		StateKey: &c.UserID,
		Type:     dataType,
		Content: map[string]interface{}{
			"expire_ts": expire.Unix(),
		},
	})
}

func (c *CSAPI) SendMultiRoomVisibilityOff(t *testing.T, dataType string, roomId string) {
	t.Helper()
	c.SendEventSynced(t, roomId, b.Event{
		StateKey: &c.UserID,
		Type:     dataType,
		Content: map[string]interface{}{
			"hidden": true,
		},
	})
}

func SyncMultiRoom(userID, dataType string, data *FakeMrd) SyncCheckOpt {
	return func(clientUserID string, topLevelSyncJSON gjson.Result) error {
		key := "multiroom." + GjsonEscape(userID) + "." + GjsonEscape(dataType)
		keyContent := key + ".content"
		mrContent := topLevelSyncJSON.Get(keyContent)
		if !mrContent.Exists() {
			return fmt.Errorf("key %s does not exist, sync body: %s", keyContent, topLevelSyncJSON.Raw)
		}
		keyTimestamp := key + ".timestamp"
		mrTimestamp := topLevelSyncJSON.Get(keyTimestamp)
		if !mrTimestamp.Exists() {
			return fmt.Errorf("key %s does not exist, sync body: %s", keyTimestamp, topLevelSyncJSON.Raw)
		}
		if mrTimestamp.Num == 0 {
			return fmt.Errorf("got timestamp equal 0")
		}
		str := mrContent.Get("Foo").String()
		if str != data.Foo {
			return fmt.Errorf("SyncMultiRoom: got %s, wanted %s, sync body: %s", str, data.Foo, topLevelSyncJSON.Raw)
		}
		return nil
	}
}

func SyncNoMultiRoom(userID, dataType string) SyncCheckOpt {
	return func(clientUserID string, topLevelSyncJSON gjson.Result) error {
		key := "multiroom." + GjsonEscape(userID) + "." + GjsonEscape(dataType)
		mrd := topLevelSyncJSON.Get(key)
		if mrd.Exists() {
			return fmt.Errorf("key %s exist, expected to not exist, sync body: %s", key, topLevelSyncJSON.Raw)
		}
		return nil
	}
}

func (c *CSAPI) MustSyncCheck(t *testing.T, syncReq SyncReq, checks ...SyncCheckOpt) string {
	res, nextBatch := c.MustSync(t, syncReq)
	for _, check := range checks {
		err := check(c.UserID, res)
		if err != nil {
			t.Fatalf("MustSyncCheck failed: %s", err)
		}
	}
	return nextBatch
}
