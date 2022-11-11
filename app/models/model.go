package models

// BaseModel 模型基类
type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id,omitempty"`
}

// CommonTimestampsField 时间戳
type CommonTimestampsField struct {
	//CreatedAt time.Time `gorm:"column:created_at;index;" json:"created_at,omitempty"`
	//UpdatedAt time.Time `gorm:"column:updated_at;index;" json:"updated_at,omitempty"`
}

type delStatus struct {
	YES, NO uint8
}

var DelStatus = delStatus{
	YES: 1,
	NO:  0,
}

type SqlOpt struct {
	Field         string
	Where         map[string]interface{}
	Order         string
	Join          string
	LeftJoinTable string
	LeftJoinOn    string
	Group         string
	Limit         int
	Offset        int
}
