package streamio

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
)

// MessageGenerator reads a stream of bytes
// produces and produces a stream of Messages.
func MessageGenerator(ctx context.Context, r io.Reader, w io.Writer) error {

	// Create our streaming decoder
	dec := json.NewDecoder(r)

	// Generate Messages
	writer := bufio.NewWriter(w)

	if _, err := dec.Token(); err != nil {
		return err
	}

	// While we have Messages
	for dec.More() && ctx.Err() == nil {

		var m Message
		err := dec.Decode(&m)
		if err != nil {
			return err
		}

		// Write our Message downstream.
		data, err := json.Marshal(m)
		if err != nil {
			return err
		}

		// Now write our Message
		if _, err := writer.Write(data); err != nil {
			return err
		}
		if err := writer.WriteByte('\n'); err != nil {
			return err
		}
		if err := writer.Flush(); err != nil {
			return err
		}
	}

	if _, err := dec.Token(); err != nil {
		return err
	}

	return ctx.Err()
}
