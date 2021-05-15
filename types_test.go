package notion

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	t.Run("MarshalAndUnmarshal", func(t *testing.T) {
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
		if !tm.Time.Equal(u.Time) {
			t.Errorf("Unmarshal: Not equal value with before marshal")
		}
	})

	t.Run("Unmarshal", func(t *testing.T) {
		u := &Time{}
		err := json.Unmarshal([]byte("{}"), u)
		if err != nil {
			t.Fatal(err)
		}
		if !u.Time.IsZero() {
			t.Errorf("\"{}\" has to be parsed Zero: %s", u.Time.Format(time.RFC3339))
		}
	})
}
