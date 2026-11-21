package defaultpostgresql

import (
	"time"

	"gorm.io/gorm"
)

// DemoUser 映射 PostgreSQL 默认实例的示例用户表。
type DemoUser struct {
	ID         int64   `gorm:"column:id;type:bigint;generated:identity;primaryKey;autoIncrement;not null;comment:自增主键ID"`
	UserNumber string  `gorm:"column:user_number;type:varchar(20);not null;uniqueIndex:uk_demo_users_user_number;comment:用户编号"`
	Nickname   string  `gorm:"column:nickname;type:varchar(50);not null;default:'';comment:昵称"`
	Phone      *string `gorm:"column:phone;type:varchar(20);default:null;uniqueIndex:uk_demo_users_phone;comment:手机号"`
	Avatar     string  `gorm:"column:avatar;type:varchar(512);not null;default:'';comment:头像地址"`
	Bio        string  `gorm:"column:bio;type:varchar(512);not null;default:'';comment:个人简介"`
	// Balance 用字符串接收定点值，避免浮点舍入。
	Balance     string         `gorm:"column:balance;type:numeric(12,2);not null;default:0.00;comment:示例账户余额"`
	Preferences map[string]any `gorm:"column:preferences;type:jsonb;serializer:json;comment:用户偏好设置"`
	Remark      *string        `gorm:"column:remark;type:text;comment:用户备注"`
	IsDisabled  bool           `gorm:"column:is_disabled;type:boolean;not null;default:false;comment:是否禁用账号"`
	LastLoginAt *time.Time     `gorm:"column:last_login_at;type:timestamptz(3);default:null;comment:最后登录时间"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamptz(3);not null;default:CURRENT_TIMESTAMP(3);comment:创建时间"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamptz(3);not null;default:CURRENT_TIMESTAMP(3);comment:更新时间"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz(3);default:null;index:idx_demo_users_deleted_at;comment:软删除时间"`
}

// TableName 返回示例用户表的固定物理表名。
func (DemoUser) TableName() string {
	return "demo_users"
}
