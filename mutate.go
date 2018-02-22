package streamio

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
)

// Mutate ...
func Mutate(ctx context.Context, r io.Reader, w io.Writer) error {

	reader := bufio.NewReader(r)
	dec := json.NewEncoder(w)
	for {

		select {

		case <-ctx.Done():
			return ctx.Err()

		default:

			// Read a Message
			data, err := reader.ReadBytes('\n')
			if err != nil {
				return err
			}

			// Parse our Message
			var m Message
			if err := json.Unmarshal(data, &m); err != nil {
				return err
			}

			// Perform mutation
			m.Payload++

			// Send our mutation downstream
			if err := dec.Encode(m); err != nil {
				return err
			}
		}
	}
}
