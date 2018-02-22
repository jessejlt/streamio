package streamio

import (
	"context"
	"io"
	"log"

	"golang.org/x/sync/errgroup"
)

// LinearPipeline streams Messages through a single Mutate
// cycle and back through the supplied Writer.
func LinearPipeline(ctx context.Context, r io.Reader, w io.Writer) error {

	group, ctx := errgroup.WithContext(ctx)

	// pw is our Generator's output,
	// pr is our Mutator's input
	pr, pw := io.Pipe()

	// Generate our Messages
	group.Go(func() error {

		if err := MessageGenerator(ctx, r, pw); err != nil {
			log.Printf("%v\n", err)
			return err
		}

		return pw.Close()
	})

	// Mutate Messages as they generated
	group.Go(func() error {

		if err := Mutate(ctx, pr, w); err != nil {
			log.Printf("%v\n", err)
			return err
		}
		return nil
	})

	return group.Wait()
}
