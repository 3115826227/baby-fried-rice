package interfaces

type DataObject interface {
	TableName() string
	Get() interface{}
}
