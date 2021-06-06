package cache

// Item ...
type Item struct {
	Data    interface{} `json:"data,omitempty"`
	Expires int64       `json:"expires,omitempty"`
}
