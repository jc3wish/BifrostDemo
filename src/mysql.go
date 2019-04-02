package src

import (
"github.com/jc3wish/Bifrost/Bristol/mysql"
"database/sql/driver"
)

func init(){

}

func NewMySQLConn(Uri string) *MySQLConn{
	db := &MySQLConn{
		Uri:Uri,
	}
	db.DBConnect()
	return db
}

type MySQLConn struct {
	Uri string
	db mysql.MysqlConnection
}

func(This *MySQLConn) DBConnect(){
	This.db = mysql.NewConnect(This.Uri)
}

func(This *MySQLConn) ExecSQL(sql string) error{
	p := make([]driver.Value, 0)
	_,err := This.db.Exec(sql,p)
	if err != nil{
		return err
	}
	return nil
}


func(This *MySQLConn) Close(){
	This.db.Close()
}
