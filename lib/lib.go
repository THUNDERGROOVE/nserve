package lib

import "time"

type TargetMetadata struct {
	retryCount int
	Started    time.Time
	Running    bool
}

type Target struct {
	Command   string   `json:"command"`
	Name      string   `json:"name, omitempty"`
	Arguments []string `json:"arguments, omitempty"`
	NoNginx   bool     `json:"nonginx"`
	// nginx shit
	ServerName string         `json:"server-name, omitempty"`
	Port       string         `json:"port, omitempty"`
	Location   string         `json:"location, omitempty"`
	Data       TargetMetadata `json:"-"`
	Running    bool           `json:"running, omitempty"`
}
