package db

type KeyValueDatabase interface {
	//Set(key []byte, value interface{}) error
	//Get(key []byte) (interface{}, error)
	KVDBName() string
}

//type RelationDatabase interface {
//}
