
package main

import (
	"encoding/json"
	"math"
	"math/rand"
	"syscall/js"
)

type GameState struct {
	X        float64     `json:"x"`
	Z        float64     `json:"z"`
	Score    int         `json:"score"`
	Health   int         `json:"health"`
	Won      bool        `json:"won"`
	Crystals [][2]float64 `json:"crystals"`
}

var state GameState

func clamp(v, min, max float64) float64 {
	return math.Max(min, math.Min(max, v))
}

func snapshot() string {
	data, _ := json.Marshal(state)
	return string(data)
}

func update(_ js.Value, args []js.Value) any {
	if len(args) < 3 {
		return snapshot()
	}
	dt := clamp(args[0].Float(), 0, 0.05)
	x := args[1].Float()
	z := args[2].Float()
	length := math.Hypot(x, z)
	if length > 1 {
		x /= length
		z /= length
	}
	speed := 6.0
	if !state.Won {
		state.X = clamp(state.X+x*speed*dt, -18, 18)
		state.Z = clamp(state.Z+z*speed*dt, -18, 18)
	}
	remaining := state.Crystals[:0]
	for _, crystal := range state.Crystals {
		distance := math.Hypot(
			state.X-crystal[0],
			state.Z-crystal[1],
		)
		if distance < 1.15 {
			state.Score++
		} else {
			remaining = append(remaining, crystal)
		}
	}
	state.Crystals = remaining
	if state.Score >= 10 {
		state.Won = true
	}
	return snapshot()
}

func reset(_ js.Value, _ []js.Value) any {
	state = GameState{
		X:      0,
		Z:      0,
		Health: 3,
	}
	r := rand.New(rand.NewSource(42))
	for i := 0; i < 10; i++ {
		x := r.Float64()*30 - 15
		z := r.Float64()*30 - 15
		if math.Hypot(x, z) < 4 {
			x += 7
		}
		state.Crystals = append(
			state.Crystals,
			[2]float64{x, z},
		)
	}
	return snapshot()
}

func main() {
	js.Global().Set("goUpdate", js.FuncOf(update))
	js.Global().Set("goReset", js.FuncOf(reset))
	reset(js.Undefined(), nil)
	select {}
}
