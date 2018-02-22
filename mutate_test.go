package streamio_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"

	"github.com/jessejlt/streamio"
)

func TestMutate(t *testing.T) {
	t.Parallel()

	// Generate our Messages
	generations := 10
	messages := make([]streamio.Message, generations)
	for i := 0; i < generations; i++ {
		messages[i].Generation = i
		messages[i].Payload = i
	}

	// Prepare Messages for transport
	data, err := json.Marshal(messages)
	if err != nil {
		t.Fatal(err)
	}

	// Start our pipeline
	r := bytes.NewReader(data)
	w := new(bytes.Buffer)
	ctx := context.Background()
	if err := streamio.LinearPipeline(ctx, r, w); err != nil {

		if err != io.EOF {
			t.Fatal(err)
		}
	}

	dec := json.NewDecoder(w)
	var wrote int
	for wrote = 0; dec.More(); wrote++ {

		var m streamio.Message
		if err := dec.Decode(&m); err != nil {
			t.Fatal(err)
		}

		if have, expect := m.Payload, m.Generation+1; have != expect {
			t.Fatalf("Mutate: Have=%d, Expect=%d", have, expect)
		}
	}

	if have, expect := wrote, generations; have != expect {
		t.Fatalf("Len: Have=%d, Expect=%d", have, expect)
	}
}
