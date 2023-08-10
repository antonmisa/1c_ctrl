// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import "time"

// Cluster -.
type Cluster struct {
	ID            string `json:"id"       rac:"cluster"                        example:"UUID"`
	Host          string `json:"host"     rac:"host"                           example:"localhost"`
	Port          string `json:"port"     rac:"port"                           example:"1541"`
	Name          string `json:"name"     rac:"name"                           example:"name"`
	Exp           int    `json:"exp"      rac:"expiration-timeout"             example:"1200"`
	LT            int    `json:"lt"       rac:"lifetime-limit"                 example:"3000"`
	MaxMemSize    int    `json:"mms"      rac:"max-memory-size"                example:"50000"`
	MaxMemTimeLim int    `json:"mmts"     rac:"max-memory-time-limit"          example:"600"`
	SecLevel      int    `json:"sl"       rac:"security-level"                 example:"0"`
	SesFTLevel    int    `json:"sftl"     rac:"session-fault-tolerance-level"  example:"0"`
	LBMode        string `json:"lb"       rac:"load-balancing-mode"            example:"perfomance"`
	ErrCountTh    int    `json:"errth"    rac:"errors-count-threshold"         example:"0"`
	KillPP        int    `json:"kpp"      rac:"kill-problem-process"           example:"0"`
}

// Infobase -.
type Infobase struct {
	ID   string `json:"id"    rac:"infobase"   example:"UUID"`
	Name string `json:"name"  rac:"name"       example:"name"`
	Desc string `json:"desc"  rac:"descr"      example:"comments"`
}

// Session -.
type Session struct {
	ID             string    `json:"id"              rac:"session"     example:"UUID"`
	SID            int       `json:"sid"             rac:"session-id"  example:"12345"`
	InfobaseID     string    `json:"ib"              rac:"infobase"    example:"UUID"`
	ConnectionID   string    `json:"conn"            rac:"connection"  example:"UUID"`
	ProcessID      string    `json:"proc"            rac:"process"     example:"UUID"`
	UserName       string    `json:"uname"           rac:"user-name"   example:"UserName"`
	Host           string    `json:"host"            rac:"host"        example:"Host"`
	AppID          string    `json:"appid"           rac:"app-id"      example:"1CV8"`
	Loc            string    `json:"loc"             rac:"locale"      example:"ru"`
	Started        time.Time `json:"started"         rac:"started-at"      example:"2023-08-10T14:04:43"`
	LastActive     time.Time `json:"active"          rac:"last-active-at"  example:"2023-08-10T14:04:43"`
	Hibernate      string    `json:"hib"             rac:"hibernate"       example:"yes/no"`
	HiberTime      int       `json:"hibtm"             rac:"passive-session-hibernate-time"  example:"1200"`
	HiberTermTime  int       `json:"hibterm"             rac:"hibernate-session-terminate-time" example:"3600"`
	BlockedDB      int       `json:"blockdb"             rac:"blocked-by-dbms"  example:"0"`
	BlockedLS      int       `json:"blockls"             rac:"blocked-by-ls"  example:"0"`
	Bytes          int       `json:"bytes"             rac:"bytes-all"  example:"12345"`
	Bytes5m        int       `json:"bytes5m"             rac:"bytes-last-5min"  example:"123"`
	Calls          int       `json:"calls"             rac:"calls-all"  example:"5"`
	Calls5m        int       `json:"calls5m"             rac:"calls-last-5min"  example:"2"`
	BytesDB        int       `json:"bytesdb"             rac:"dbms-bytes-all"  example:"123"`
	BytesDB5m      int       `json:"bytesdb5m"             rac:"dbms-bytes-last-5min"  example:"12"`
	DBProcInfo     string    `json:"dbproci"             rac:"db-proc-info"  example:""`
	DBProc         int       `json:"dbproc"             rac:"db-proc-took"  example:"123"`
	DBProcAt       string    `json:"dbprocat"             rac:"db-proc-took-at"  example:""`
	Duration       int       `json:"dur"             rac:"duration-all"  example:"100"`
	DurationDB     int       `json:"durdb"             rac:"duration-all-dbms"  example:"100"`
	DurationCur    int       `json:"durcur"             rac:"duration-current"  example:"80"`
	DurationCurDB  int       `json:"durcurdb"             rac:"duration-current-dbms"  example:"80"`
	Duration5m     int       `json:"dur5m"             rac:"duration-last-5min"  example:"100"`
	DurationDB5m   int       `json:"durdb5m"             rac:"duration-last-5min-dbms"  example:"100"`
	MemoryCur      int       `json:"memcur"             rac:"memory-current"  example:"12345"`
	Memory5m       int       `json:"mem5m"             rac:"memory-last-5min"  example:"1234"`
	Memory         int       `json:"mem"             rac:"memory-total"  example:"123456"`
	ReadCur        int       `json:"readcur"             rac:"read-current"  example:"5678"`
	Read5m         int       `json:"read5m"             rac:"read-last-5min"  example:"56"`
	Read           int       `json:"read"             rac:"read-total"  example:"56789"`
	WriteCur       int       `json:"writecur"             rac:"write-current"  example:"123"`
	Write5m        int       `json:"write5m"             rac:"write-last-5min"  example:"123"`
	Write          int       `json:"write"             rac:"write-total"  example:"123"`
	DurationSvcCur int       `json:"dursvccur"             rac:"duration-current-service"  example:"0"`
	DurationSvc5m  int       `json:"dursvc5m"             rac:"duration-last-5min-service"  example:"0"`
	DurationSvc    int       `json:"dursvc"             rac:"duration-all-service"  example:"0"`
	Svc            string    `json:"svc"             rac:"current-service-name"  example:"Name"`
	CPUCur         int       `json:"cpucur"             rac:"cpu-time-current"  example:"123"`
	CPU5m          int       `json:"cpu5m"             rac:"cpu-time-last-5min"  example:"12"`
	CPU            int       `json:"cpu"             rac:"cpu-time-total"  example:"1234"`
	Sep            string    `json:"sep"             rac:"data-separation"  example:""`
}

// Connection -.
type Connection struct {
	ID         string    `json:"id"          rac:"connection" example:"UUID"`
	CID        int       `json:"cid"         rac:"connection-id"  example:"12345"`
	InfobaseID string    `json:"ib"          rac:"infobase" example:"UUID"`
	ProcessID  string    `json:"proc"        rac:"process" example:"UUID"`
	Host       string    `json:"host"        rac:"host" example:"localhost"`
	AppID      string    `json:"appid"       rac:"application" example:"1CV8"`
	Connected  time.Time `json:"connected"   rac:"connected-at"      example:"2023-08-10T11:40:55"`
	SID        int       `json:"sid"         rac:"session-number"  example:"12345"`
	Blocked    int       `json:"blocked"     rac:"blocked-by-ls"  example:"0"`
}
