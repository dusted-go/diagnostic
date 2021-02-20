package trace

import (
	"bytes"
	"context"
	rng "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
)

// --------------------------------
// ID
// --------------------------------

// ID is the unique identifier for an entire trace.
type ID [16]byte

var emptyID ID

// String returns the hex string representation of the ID.
func (id ID) String() string {
	return hex.EncodeToString(id[:])
}

// IsValid checks if the trace ID does not consist of zeros only.
func (id ID) IsValid() bool {
	return !bytes.Equal(id[:], emptyID[:])
}

// --------------------------------
// SpanID
// --------------------------------

// SpanID is the unique identifier for a span of a trace.
type SpanID [8]byte

var emptySpanID SpanID

// String returns the hex string representation of the ID.
func (id SpanID) String() string {
	return hex.EncodeToString(id[:])
}

// Decimal returns the decimal representation of the ID.
func (id SpanID) Decimal() uint64 {
	return binary.LittleEndian.Uint64(id[:])
}

// IDGenerator allows custom generators for TraceID and SpanID.
type IDGenerator interface {
	NewTraceIDs() (ID, SpanID)
	NewSpanID() SpanID
}

// IsValid checks if the span ID does not consist of zeros only.
func (id SpanID) IsValid() bool {
	return !bytes.Equal(id[:], emptySpanID[:])
}

// --------------------------------
// Parsers
// --------------------------------

func decodeHex(value string, dest []byte) error {
	for _, r := range value {
		switch {
		case 'a' <= r && r <= 'f':
			continue
		case '0' <= r && r <= '9':
			continue
		default:
			return errors.New("invalid hexadecimal value")
		}
	}
	decoded, err := hex.DecodeString(value)
	if err != nil {
		return err
	}
	copy(dest, decoded)
	return nil
}

// ParseID returns a trace ID from a hexadecimal string if it meets the W3C specification.
// See more at: https://www.w3.org/TR/trace-context/#trace-id
func ParseID(value string) (ID, error) {
	id := ID{}
	if len(value) != 32 {
		return id, errors.New("cannot parse trace ID because the string value has an invalid length (must be 32 characters long)")
	}

	if err := decodeHex(value, id[:]); err != nil {
		return id, err
	}

	if !id.IsValid() {
		return id, errors.New("invalid/empty trace ID")
	}
	return id, nil
}

// ParseOpenTelemetrySpanID returns a span ID from a hexadecimal string if it meets the W3C specification.
// See more at: https://www.w3.org/TR/trace-context/#parent-id
func ParseOpenTelemetrySpanID(value string) (SpanID, error) {
	id := SpanID{}
	if len(value) != 16 {
		return id, errors.New("cannot parse span ID because the string value has an invalid length (must be 16 characters long)")
	}

	if err := decodeHex(value, id[:]); err != nil {
		return id, err
	}

	if !id.IsValid() {
		return id, errors.New("invalid/empty span ID")
	}
	return id, nil
}

// ParseGoogleCloudSpanID returns a span ID from a string holding an unsigned int64 value.
func ParseGoogleCloudSpanID(value string) (SpanID, error) {
	id := SpanID{}

	num, err := strconv.ParseUint(value, 0, 64)
	if err != nil {
		return id, fmt.Errorf("error paring Google Cloud SpanID to uint64: %w", err)
	}

	binary.LittleEndian.PutUint64(id[:], num)

	if !id.IsValid() {
		return id, errors.New("invalid/empty span ID")
	}
	return id, nil
}

// --------------------------------
// Default Generator
// --------------------------------

type randomIDGenerator struct {
	sync.Mutex
	randSource *rand.Rand
}

// NewSpanID returns a non-zero span ID from a randomly-chosen sequence.
func (gen *randomIDGenerator) NewSpanID() SpanID {
	gen.Lock()
	defer gen.Unlock()
	sid := SpanID{}
	gen.randSource.Read(sid[:])
	return sid
}

// NewTraceIDs returns a non-zero trace ID and a non-zero span ID from a randomly-chosen sequence.
func (gen *randomIDGenerator) NewTraceIDs() (ID, SpanID) {
	gen.Lock()
	defer gen.Unlock()
	tid := ID{}
	gen.randSource.Read(tid[:])
	sid := SpanID{}
	gen.randSource.Read(sid[:])
	return tid, sid
}

func defaultIDGenerator() IDGenerator {
	gen := &randomIDGenerator{}
	var rngSeed int64
	_ = binary.Read(rng.Reader, binary.LittleEndian, &rngSeed)
	gen.randSource = rand.New(rand.NewSource(rngSeed))
	return gen
}

// DefaultGenerator is the default trace ID and span ID generator.
var DefaultGenerator = defaultIDGenerator()

// --------------------------------
// Context Helpers
// --------------------------------

// Custom types to avoid key collisions in the context object.
type traceIDKey int
type spanIDKey int

// IDKey is the key that references the trace ID inside context.
const IDKey traceIDKey = 0

// SpanIDKey is the key that references the span ID inside context.
const SpanIDKey spanIDKey = 0

// TryGetID tries to get a previously saved trace ID.
func TryGetID(ctx context.Context) (ID, bool) {
	if ctx == nil {
		return ID{}, false
	}
	if traceID, ok := ctx.Value(IDKey).(ID); ok {
		return traceID, true
	}
	return ID{}, false
}

// TryGetSpanID tries to get a previously saved span ID.
func TryGetSpanID(ctx context.Context) (SpanID, bool) {
	if ctx == nil {
		return SpanID{}, false
	}
	if spanID, ok := ctx.Value(SpanIDKey).(SpanID); ok {
		return spanID, true
	}
	return SpanID{}, false
}

// Context adds a trace ID and span ID to the current context.
func Context(ctx context.Context, traceID ID, spanID SpanID) context.Context {
	ctx = context.WithValue(ctx, IDKey, traceID)
	ctx = context.WithValue(ctx, SpanIDKey, spanID)
	return ctx
}
