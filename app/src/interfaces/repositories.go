package interfaces

type DbHandler interface {
	Execute(statement string, model interface{})
	Query(statement string, params map[string]interface{}) Row
	Trunc(table string)
}

type Row interface {
	Scan(dest ...interface{})
	Next() bool
}

type DbRepo struct {
	dbHandlers map[string]DbHandler
	dbHandler  DbHandler
}
