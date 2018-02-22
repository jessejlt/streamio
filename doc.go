// Package streamio demonstrates how to use
// the io package to define streaming pipelines.
package streamio

// Message is the data we'll use within our
// processing pipelines.
type Message struct {
	Generation int
	Payload    int
}
