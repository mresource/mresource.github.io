//go:build js && wasm

// Command wasm exposes the progress package to the browser as
// window.mresourceProgress.decide(jsonString) → jsonString.
//
// Input:  {"plan":{"baseline":42,"goal":70,"weeks":12},"scores":[43,44],"changes":[4]}
// Output: {"decision":"collect_more","sinceChange":2,"aim":[42,44.33,…]} or {"error":"…"}
package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/mresource/mresource.github.io/go/progress"
)

type request struct {
	Plan    progress.Plan `json:"plan"`
	Scores  []float64     `json:"scores"`
	Changes []int         `json:"changes"`
}

type response struct {
	progress.Result
	Aim   []float64 `json:"aim,omitempty"` // aim value for week 0..max(weeks, len(scores))
	Error string    `json:"error,omitempty"`
}

func decide(_ js.Value, args []js.Value) (result any) {
	defer func() {
		if r := recover(); r != nil {
			result = encode(response{Error: "internal error"})
		}
	}()
	if len(args) != 1 {
		return encode(response{Error: "decide expects one JSON string"})
	}
	var req request
	if err := json.Unmarshal([]byte(args[0].String()), &req); err != nil {
		return encode(response{Error: "invalid JSON: " + err.Error()})
	}
	res, err := progress.Decide(req.Plan, req.Scores, req.Changes)
	if err != nil {
		return encode(response{Error: err.Error()})
	}
	out := response{Result: res}
	n := max(req.Plan.Weeks, len(req.Scores))
	for w := 0; w <= n; w++ {
		out.Aim = append(out.Aim, req.Plan.AimAt(float64(w)))
	}
	return encode(out)
}

func encode(r response) string {
	b, err := json.Marshal(r)
	if err != nil {
		return `{"error":"could not encode result"}`
	}
	return string(b)
}

func main() {
	js.Global().Set("mresourceProgress", js.ValueOf(map[string]any{
		"decide": js.FuncOf(decide),
	}))
	js.Global().Call("dispatchEvent", js.Global().Get("Event").New("mresource-progress-ready"))
	select {} // keep the Go runtime alive for callbacks
}
