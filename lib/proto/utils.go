package proto

import (
	"context"
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-rod/rod/lib/utils"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Client interface to send the request.
// So that this lib doesn't handle anything has side effect.
type Client interface {
	Call(ctx context.Context, sessionID, methodName string, params json.RawMessage) (res []byte, err error)
}

// Payload interface returns the name of the event, such as "Page.loadEventFired"
type Payload interface {
	// MethodName is called method name is because the json-schema definition of it is "method".
	// And "eventName" is already used by a lot of existing fields.
	MethodName() string
}

// Caller interface to get the context of the request
type Caller interface {
	// CallContext returns ctx, client, and the sessionID
	CallContext() (context.Context, Client, string)
}

// Call method with request and response containers.
func Call(method string, req, res interface{}, caller Caller) error {
	ctx, client, id := caller.CallContext()

	payload, err := Normalize(req)
	if err != nil {
		return err
	}

	bin, err := client.Call(ctx, id, method, payload)
	if err != nil {
		return err
	}

	if res != nil {
		err = json.Unmarshal(bin, res)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetType from method name of this package,
// such as proto.GetType("Page.enable") will return the type of proto.PageEnable
func GetType(methodName string) reflect.Type {
	return types[methodName]
}

// ParseMethodName to domain and name
func ParseMethodName(method string) (domain, name string) {
	arr := strings.Split(method, ".")
	return arr[0], arr[1]
}

// Normalizable interface to transform the params into the correct data structure before being sent by the client.
// Because the json-schema doesn't cover all the type constrains of the protocol, we need this extra layer to do
// the normalization.
// Such as when send mouse wheel events, the deltaX and deltaY can't be omitted. The json-schema is wrong for them.
type Normalizable interface {
	Normalize() (json.RawMessage, error)
}

// Normalize the method payload
func Normalize(m interface{}) (json.RawMessage, error) {
	n, ok := m.(Normalizable)
	if ok {
		return n.Normalize()
	}
	return json.Marshal(m)
}

// JSON represent a JSON value
type JSON struct {
	gjson.Result
	raw []byte
}

// NewJSON json object
func NewJSON(val interface{}) JSON {
	j := JSON{}
	if b, ok := val.([]byte); ok {
		j.raw = b
	} else {
		j.raw = utils.MustToJSONBytes(val)
	}
	j.Result = gjson.ParseBytes(j.raw)
	return j
}

// UnmarshalJSON interface
func (j *JSON) UnmarshalJSON(b []byte) error {
	j.raw = b
	j.Result = gjson.ParseBytes(b)
	return nil
}

// MarshalJSON interface
func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j.raw) == 0 {
		return []byte("null"), nil
	}
	return j.raw, nil
}

// Join elements
func (j JSON) Join(sep string) string {
	list := []string{}

	for _, el := range j.Array() {
		list = append(list, el.String())
	}

	return strings.Join(list, sep)
}

// Set a json value for the specified path and return the new JSON.
func (j JSON) Set(path string, value interface{}) (JSON, error) {
	b, _ := j.MarshalJSON()
	b, err := sjson.SetBytes(b, path, value)
	return NewJSON(b), err
}

// TimeSinceEpoch UTC time in seconds, counted from January 1, 1970.
type TimeSinceEpoch struct {
	time.Time
}

// UnmarshalJSON interface
func (t *TimeSinceEpoch) UnmarshalJSON(b []byte) error {
	t.Time = (time.Unix(0, 0)).Add(
		time.Duration(gjson.ParseBytes(b).Float() * float64(time.Second)),
	)
	return nil
}

// MarshalJSON interface
func (t TimeSinceEpoch) MarshalJSON() ([]byte, error) {
	d := float64(t.Time.UnixNano()) / float64(time.Second)
	return json.Marshal(d)
}

// MonotonicTime Monotonically increasing time in seconds since an arbitrary point in the past.
type MonotonicTime struct {
	time.Duration
}

// UnmarshalJSON interface
func (t *MonotonicTime) UnmarshalJSON(b []byte) error {
	t.Duration = time.Duration(gjson.ParseBytes(b).Float() * float64(time.Second))
	return nil
}

// MarshalJSON interface
func (t MonotonicTime) MarshalJSON() ([]byte, error) {
	d := float64(t.Duration) / float64(time.Second)
	return json.Marshal(d)
}

var _ Normalizable = InputDispatchMouseEvent{}

// Normalize interface
// All the fields are staticcally typed, this function can't have error.
func (e InputDispatchMouseEvent) Normalize() (json.RawMessage, error) {
	data, _ := json.Marshal(e)

	if e.Type == InputDispatchMouseEventTypeMouseWheel {
		data, _ = sjson.SetBytes(data, "deltaX", e.DeltaX)
		data, _ = sjson.SetBytes(data, "deltaY", e.DeltaY)
	}

	return data, nil
}

// PatternToReg FetchRequestPattern.URLPattern to regular expression
func PatternToReg(pattern string) string {
	if pattern == "" {
		return ""
	}

	pattern = " " + pattern
	pattern = regexp.MustCompile(`([^\\])\*`).ReplaceAllString(pattern, "$1.*")
	pattern = regexp.MustCompile(`([^\\])\?`).ReplaceAllString(pattern, "$1.")

	return `\A` + strings.TrimSpace(pattern) + `\z`
}

// Point from the origin (0, 0)
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Len is the number of vertices
func (q DOMQuad) Len() int {
	return len(q) / 2
}

// Each point
func (q DOMQuad) Each(fn func(pt Point, i int)) {
	for i := 0; i < q.Len(); i++ {
		fn(Point{q[i*2], q[i*2+1]}, i)
	}
}

// Center of the polygon
func (q DOMQuad) Center() Point {
	var x, y float64
	q.Each(func(pt Point, _ int) {
		x += pt.X
		y += pt.Y
	})
	return Point{x / float64(q.Len()), y / float64(q.Len())}
}

// OnePointInside the shape
func (res *DOMGetContentQuadsResult) OnePointInside() *Point {
	if len(res.Quads) == 0 {
		return nil
	}

	center := res.Quads[0].Center()

	return &center
}

// MoveTo X and Y to x and y
func (p *InputTouchPoint) MoveTo(x, y float64) {
	p.X = x
	p.Y = y
}
