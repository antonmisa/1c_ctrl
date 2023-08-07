// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

// Cluster -.
type Cluster struct {
	ID   string `json:"id"       example:"UUID like"`
	Host string `json:"host"  example:"localhost"`
	Port string `json:"port"     example:"1541"`
	Name string `json:"name"  example:"name as text"`
}

// Infobase -.
type Infobase struct {
	ID   string `json:"id"       example:"UUID like"`
	Name string `json:"name"  example:"name as text"`
	Desc string `json:"desc"  example:"some comments"`
}

// Session -.
type Session struct {
	ID           string `json:"id"       example:"UUID like"`
	InfobaseID   string `json:"idbase"  example:"UUID of infobase"`
	ConnectionID string `json:"idbconn"  example:"UUID of connection"`
	ProcessID    string `json:"idproc"  example:"UUID of process"`
	UserName     string `json:"uname"  example:"Name of the user"`
	Host         string `json:"host"  example:"Host of the user"`
	AppID        string `json:"appid"  example:"Application identifier"`
}

// Connection -.
type Connection struct {
	ID         string `json:"id"       example:"UUID like"`
	InfobaseID string `json:"idbase"  example:"UUID of infobase"`
	ProcessID  string `json:"idproc"  example:"UUID of process"`
	Host       string `json:"host"  example:"Host of the user"`
	AppID      string `json:"appid"  example:"Application identifier"`
}
