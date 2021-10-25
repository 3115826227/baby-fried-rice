package tables

type Area struct {
	Code       string
	Name       string
	ParentCode string
	Level      int
}

func (table *Area) TableName() string {
	return "area"
}

func (table *Area) Get() interface{} {
	return *table
}
