package model

type Student struct {
	Name         string         `gorm:"column:name;type:VARCHAR(512)"`
	Age          int            `gorm:"column:age;type:tinyint(4)"`
}

func (Student) TableName() string {
	return "student"
}
