package models

import (
	"encoding/json"
)

type KV struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}
