package db

//KVDBType for user to choose which key-value database
type KVDBType int

const (
	//KVDBLevel get level database
	KVDBLevel = iota
	//KVDBRedis get redis database
	KVDBRedis
)

//GetKVDBInstance returns specified key-value database instance
func GetKVDBInstance(t KVDBType) KeyValueDatabase {
	switch t {
	case KVDBLevel:
		return GetLevelDBInstance()
	case KVDBRedis:
		return GetRedisDBInstance()
	default:
		return GetLevelDBInstance()
	}
}
