package tables

type SessionUserRelation struct {
	SessionID int64  `gorm:"column:session_id;unique_index:session_user_relation" json:"session_id"`
	UserID    string `gorm:"column:user_id;unique_index:session_user_relation" json:"user_id"`
	JoinTime  int64  `gorm:"column:join_time;" json:"join_time"`
	Delete    bool   `gorm:"column:delete" json:"delete"`
}

func (table *SessionUserRelation) TableName() string {
	return "baby_im_session_user_rel"
}

type Session struct {
	ID                 int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name               string `gorm:"column:name;not null" json:"name"`
	SessionType        int32  `gorm:"column:session_type;" json:"session_type"`
	JoinPermissionType int32  `gorm:"column:join_permission_type" json:"join_permission_type"`
	Origin             string `gorm:"column:origin;not null" json:"origin"`
	CreateTime         int64  `json:"create_time"`
	UpdateTime         int64  `json:"update_time"`
}

func (table *Session) TableName() string {
	return "baby_im_session"
}

type MessageUserRelation struct {
	MessageID     int64  `gorm:"column:message_id;primaryKey" json:"message_id"`
	SessionID     int64  `gorm:"column:session_id;" json:"session_id"`
	Receive       string `gorm:"column:receive;primaryKey" json:"receive"`
	Read          bool   `gorm:"column:read;not null" json:"read"`
	SendTimestamp int64  `gorm:"column:send_timestamp;not null" json:"send_time"`
}

func (table *MessageUserRelation) TableName() string {
	return "baby_im_message_user_rel"
}

type Message struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SessionID     int64  `gorm:"column:session_id;primaryKey" json:"session_id"`
	MessageType   int32  `gorm:"column:message_type;not null" json:"message_type"`
	Send          string `gorm:"column:send;not null" json:"send"`
	Content       []byte `gorm:"column:content" json:"content"`
	SendTimestamp int64  `gorm:"column:send_timestamp;not null" json:"send_time"`
}

func (table *Message) TableName() string {
	return "baby_im_message"
}
