package notion

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	tm := Time{Time: time.Date(2020, 5, 3, 14, 15, 30, 0, time.Local)}

	b, err := json.Marshal(tm)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(b, []byte("\"2020-05-03T14:15:30")) {
		t.Errorf("Marshal: Invalid format: %s", string(b))
	}

	u := &Time{}
	err = json.Unmarshal(b, u)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
}
