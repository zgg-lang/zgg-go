package builtin_libs

import (
	"database/sql"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"

	. "github.com/zgg-lang/zgg-go/runtime"
)

var (
	dbQueryResultClass  ValueType
	dbSessionClass      ValueType
	dbActiveRecordClass ValueType
)

type (
	dbDialect interface {
		ShowTablesSQL() string
		Quote(name string) string
	}

	dbCommonDialect struct{}
	dbMysqlDialect  struct{}
)

func (dbCommonDialect) Quote(name string) string { return name }
func (dbCommonDialect) ShowTablesSQL() string    { return "SHOW TABLES" }
func (dbMysqlDialect) Quote(name string) string  { return fmt.Sprintf("`%s`", name) }
func (dbMysqlDialect) ShowTablesSQL() string     { return "SHOW TABLES" }

var (
	dbDialectMap = map[string]dbDialect{
		"mysql": dbMysqlDialect{},
	}
)

func libDb(*Context) ValueObject {
	lib := NewObject()
	queryResultClass := initQueryResultClass()
	dbClass := initDatabaseClass(queryResultClass)
	lib.SetMember("Database", dbClass, nil)
	lib.SetMember("QueryResult", queryResultClass, nil)
	lib.SetMember("open", NewNativeFunction("open", func(c *Context, this Value, args []Value) Value {
		var (
			engine ValueStr
			uri    ValueStr
		)
		if len(args) == 1 {
			dsn := args[0].ToString(c)
			dsnUrl, err := url.Parse(dsn)
			if err != nil {
				c.OnRuntimeError("db.open parse dsn error %s", err)
				return nil
			}
			engine = NewStr(dsnUrl.Scheme)
			var uriStr string
			if u := dsnUrl.User; u != nil {
				uriStr += u.Username()
				if p, ok := u.Password(); ok {
					uriStr += ":" + p
				}
				uriStr += "@"
			}
			uriStr += fmt.Sprintf("tcp(%s)%s", dsnUrl.Host, dsnUrl.Path)
			if q := dsnUrl.RawQuery; q != "" {
				uriStr += "?" + q
			}
			uri = NewStr(uriStr)
		} else {
			EnsureFuncParams(c, "db.open", args,
				ArgRuleRequired{"engine", TypeStr, &engine},
				ArgRuleRequired{"uri", TypeStr, &uri},
			)
		}
		driverFound := false
		for _, d := range sql.Drivers() {
			if engine.Value() == d {
				driverFound = true
				break
			}
		}
		if !driverFound {
			c.OnRuntimeError("db.open unexpected driver %s", engine.Value())
			return nil
		}
		db, err := sql.Open(engine.Value(), uri.Value())
		if err != nil {
			c.OnRuntimeError("db.open fail %s", err)
			return nil
		}
		dialect, dialectFound := dbDialectMap[engine.Value()]
		if !dialectFound {
			dialect = dbCommonDialect{}
		}
		return NewObjectAndInit(dbClass, c, NewGoValue(db), NewGoValue(dialect))
	}), nil)
	return lib
}

func dbScanRowsToArray(c *Context, rows *sql.Rows, colTypes []*sql.ColumnType, cols []string) (ret ValueArray) {
	var err error
	if colTypes == nil {
		colTypes, err = rows.ColumnTypes()
		if err != nil {
			c.OnRuntimeError("QueryResult.__init__ get column types error %s", err)
			return
		}
	}
	fields := make([]interface{}, len(colTypes))
	for i, ct := range colTypes {
		fields[i] = reflect.New(ct.ScanType()).Interface()
	}
	if err := rows.Scan(fields...); err != nil {
		c.OnRuntimeError("QueryResult.next scan fields error %s", err)
		return
	}
	item := NewArray()
	for i := range cols {
		if fields[i] == nil {
			item.PushBack(Nil())
			continue
		}
		set := false
		switch fv := fields[i].(type) {
		case *sql.NullInt32:
			if fv.Valid {
				item.PushBack(NewInt(int64(fv.Int32)))
			} else {
				item.PushBack(Nil())
			}
			set = true
		case *sql.NullInt64:
			if fv.Valid {
				item.PushBack(NewInt(fv.Int64))
			} else {
				item.PushBack(Nil())
			}
			set = true
		case *sql.NullFloat64:
			if fv.Valid {
				item.PushBack(NewFloat(fv.Float64))
			} else {
				item.PushBack(Nil())
			}
			set = true
		case *sql.NullBool:
			if fv.Valid {
				item.PushBack(NewBool(fv.Bool))
			} else {
				item.PushBack(Nil())
			}
			set = true
		case *sql.NullString:
			if fv.Valid {
				item.PushBack(NewStr(fv.String))
			} else {
				item.PushBack(Nil())
			}
			set = true
		case *sql.NullTime:
			if fv.Valid {
				item.PushBack(NewObjectAndInit(timeClass, c, NewInt(fv.Time.UnixNano())))
			} else {
				item.PushBack(Nil())
			}
			set = true
		case *sql.RawBytes:
			switch colTypes[i].DatabaseTypeName() {
			case "DECIMAL":
				{
					v, err := strconv.ParseFloat(string(*fv), 64)
					if err != nil {
						c.OnRuntimeError("parse db DECIMAL value %s err %s", string(*fv), err)
					}
					item.PushBack(NewFloat(v))
					set = true
				}
			case "BLOB":
				{
					item.PushBack(NewBytes([]byte(*fv)))
					set = true
				}
			default:
				{
					item.PushBack(NewStr(string(*fv)))
					set = true
				}
			}
		default:
		}
		if !set {
			item.PushBack(FromGoValue(reflect.ValueOf(fields[i]).Elem(), c))
		}
	}
	ret = item
	return
}

func dbScanRowsToObject(c *Context, rows *sql.Rows, colTypes []*sql.ColumnType, cols []string) (ret ValueObject) {
	var err error
	if colTypes == nil {
		colTypes, err = rows.ColumnTypes()
		if err != nil {
			c.OnRuntimeError("QueryResult.__init__ get column types error %s", err)
			return
		}
	}
	fields := make([]interface{}, len(colTypes))
	for i, ct := range colTypes {
		fields[i] = reflect.New(ct.ScanType()).Interface()
	}
	if err := rows.Scan(fields...); err != nil {
		c.OnRuntimeError("QueryResult.next scan fields error %s", err)
		return
	}
	item := NewObject()
	for i, colName := range cols {
		if fields[i] == nil {
			item.SetMember(colName, Nil(), c)
			continue
		}
		set := false
		switch fv := fields[i].(type) {
		case *sql.NullInt32:
			if fv.Valid {
				item.SetMember(colName, NewInt(int64(fv.Int32)), c)
			} else {
				item.SetMember(colName, Nil(), c)
			}
			set = true
		case *sql.NullInt64:
			if fv.Valid {
				item.SetMember(colName, NewInt(fv.Int64), c)
			} else {
				item.SetMember(colName, Nil(), c)
			}
			set = true
		case *sql.NullFloat64:
			if fv.Valid {
				item.SetMember(colName, NewFloat(fv.Float64), c)
			} else {
				item.SetMember(colName, Nil(), c)
			}
			set = true
		case *sql.NullBool:
			if fv.Valid {
				item.SetMember(colName, NewBool(fv.Bool), c)
			} else {
				item.SetMember(colName, Nil(), c)
			}
			set = true
		case *sql.NullString:
			if fv.Valid {
				item.SetMember(colName, NewStr(fv.String), c)
			} else {
				item.SetMember(colName, Nil(), c)
			}
			set = true
		case *sql.NullTime:
			if fv.Valid {
				item.SetMember(colName, NewObjectAndInit(timeClass, c, NewInt(fv.Time.UnixNano())), c)
			} else {
				item.SetMember(colName, Nil(), c)
			}
			set = true
		case *sql.RawBytes:
			switch colTypes[i].DatabaseTypeName() {
			case "DECIMAL":
				{
					v, err := strconv.ParseFloat(string(*fv), 64)
					if err != nil {
						c.OnRuntimeError("parse db DECIMAL value %s err %s", string(*fv), err)
					}
					item.SetMember(colName, NewFloat(v), c)
					set = true
				}
			case "BLOB":
				{
					item.SetMember(colName, NewBytes([]byte(*fv)), c)
					set = true
				}
			default:
				{
					item.SetMember(colName, NewStr(string(*fv)), c)
					set = true
				}
			}
		default:
		}
		if !set {
			item.SetMember(colName, FromGoValue(reflect.ValueOf(fields[i]).Elem(), c), c)
		}
	}
	ret = item
	return
}

func initQueryResultClass() ValueType {
	dbQueryResultClass = NewClassBuilder("QueryResult").
		Constructor(func(c *Context, thisObj ValueObject, args []Value) {
			rows := args[0].ToGoValue().(*sql.Rows)
			colTypes, err := rows.ColumnTypes()
			if err != nil {
				c.OnRuntimeError("QueryResult.__init__ get column types error %s", err)
				return
			}
			thisObj.SetMember("_rows", args[0], c)
			thisObj.SetMember("_colTypes", NewGoValue(colTypes), c)
			return
		}).
		Method("each", func(c *Context, this ValueObject, args []Value) Value {
			c.AssertArgNum(len(args), 1, 1, "QueryResult.each")
			callback := c.MustCallable(args[0])
			cts := this.GetMember("_colTypes", c).ToGoValue().([]*sql.ColumnType)
			rows := this.GetMember("_rows", c).ToGoValue().(*sql.Rows)
			cols, err := rows.Columns()
			if err != nil {
				c.OnRuntimeError("QueryResult.next get columns error %s", err)
				return nil
			}
			// for i, ct := range cts {
			// 	fmt.Println(cols[i], ct.DatabaseTypeName(), ct.ScanType())
			// }
			for rows.Next() {
				row := dbScanRowsToObject(c, rows, cts, cols)
				c.Invoke(callback, Undefined(), Args(row))
			}
			return Undefined()
		}).
		Method("__iter__", func(c *Context, this ValueObject, args []Value) Value {
			cts := this.GetMember("_colTypes", c).ToGoValue().([]*sql.ColumnType)
			rows := this.GetMember("_rows", c).ToGoValue().(*sql.Rows)
			cols, err := rows.Columns()
			if err != nil {
				c.OnRuntimeError("QueryResult.next get columns error %s", err)
				return nil
			}
			return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
				if !rows.Next() {
					return NewArrayByValues(Nil(), NewBool(false))
				}
				row := dbScanRowsToObject(c, rows, cts, cols)
				return NewArrayByValues(row, NewBool(true))
			})
		}).
		Method("all", func(c *Context, this ValueObject, args []Value) Value {
			cts := this.GetMember("_colTypes", c).ToGoValue().([]*sql.ColumnType)
			rows := this.GetMember("_rows", c).ToGoValue().(*sql.Rows)
			cols, err := rows.Columns()
			if err != nil {
				c.OnRuntimeError("QueryResult.next get columns error %s", err)
				return nil
			}
			rv := NewArray()
			for rows.Next() {
				row := dbScanRowsToObject(c, rows, cts, cols)
				rv.PushBack(row)
			}
			return rv
		}).
		Method("allArray", func(c *Context, this ValueObject, args []Value) Value {
			cts := this.GetMember("_colTypes", c).ToGoValue().([]*sql.ColumnType)
			rows := this.GetMember("_rows", c).ToGoValue().(*sql.Rows)
			cols, err := rows.Columns()
			if err != nil {
				c.OnRuntimeError("QueryResult.allArray get columns error %s", err)
				return nil
			}
			rv := NewArray()
			for rows.Next() {
				row := dbScanRowsToArray(c, rows, cts, cols)
				rv.PushBack(row)
			}
			return rv
		}).
		Method("allOne", func(c *Context, this ValueObject, args []Value) Value {
			cts := this.GetMember("_colTypes", c).ToGoValue().([]*sql.ColumnType)
			rows := this.GetMember("_rows", c).ToGoValue().(*sql.Rows)
			cols, err := rows.Columns()
			if err != nil {
				c.OnRuntimeError("QueryResult.allOne get columns error %s", err)
				return nil
			}
			rv := NewArray()
			for rows.Next() {
				row := dbScanRowsToArray(c, rows, cts, cols)
				rv.PushBack(row.GetIndex(0, c))
			}
			return rv
		}).
		Method("next", func(c *Context, this ValueObject, args []Value) Value {
			rows := this.GetMember("_rows", c).ToGoValue().(*sql.Rows)
			if !rows.Next() {
				return Nil()
			}
			cts := this.GetMember("_colTypes", c).ToGoValue().([]*sql.ColumnType)
			cols, err := rows.Columns()
			if err != nil {
				c.OnRuntimeError("QueryResult.next get columns error %s", err)
				return nil
			}
			return dbScanRowsToObject(c, rows, cts, cols)
		}).
		Method("nextArray", func(c *Context, this ValueObject, args []Value) Value {
			rows := this.GetMember("_rows", c).ToGoValue().(*sql.Rows)
			if !rows.Next() {
				return Nil()
			}
			cts := this.GetMember("_colTypes", c).ToGoValue().([]*sql.ColumnType)
			cols, err := rows.Columns()
			if err != nil {
				c.OnRuntimeError("QueryResult.nextArray get columns error %s", err)
				return nil
			}
			return dbScanRowsToArray(c, rows, cts, cols)
		}).
		Method("nextOne", func(c *Context, this ValueObject, args []Value) Value {
			rows := this.GetMember("_rows", c).ToGoValue().(*sql.Rows)
			if !rows.Next() {
				return Nil()
			}
			cts := this.GetMember("_colTypes", c).ToGoValue().([]*sql.ColumnType)
			cols, err := rows.Columns()
			if err != nil {
				c.OnRuntimeError("QueryResult.nextArray get columns error %s", err)
				return nil
			}
			return dbScanRowsToArray(c, rows, cts, cols).GetIndex(0, c)
		}).
		Method("toTable", func(c *Context, this ValueObject, args []Value) Value {
			cts := this.GetMember("_colTypes", c).ToGoValue().([]*sql.ColumnType)
			rows := this.GetMember("_rows", c).ToGoValue().(*sql.Rows)
			cols, err := rows.Columns()
			if err != nil {
				c.OnRuntimeError("QueryResult.next get columns error %s", err)
				return nil
			}
			colsValue := make([]Value, len(cols))
			for i, col := range cols {
				colsValue[i] = NewStr(col)
			}
			table := NewObjectAndInit(ptablePTableClass, c, colsValue...)
			for rows.Next() {
				row := dbScanRowsToArray(c, rows, cts, cols)
				c.InvokeMethod(table, "add", func() []Value {
					r := make([]Value, len(cols))
					for i := range cols {
						r[i] = row.GetIndex(i, c)
					}
					return r
				})
			}
			return table
		}).
		Method("close", func(c *Context, this ValueObject, args []Value) Value {
			rows := this.GetMember("_rows", c).ToGoValue().(*sql.Rows)
			if err := rows.Close(); err != nil {
				c.OnRuntimeError("QeuryResult.close error %s", err)
				return nil
			}
			return Undefined()
		}).
		Build()
	return dbQueryResultClass
}

func initDatabaseClass(queryResultClass ValueType) ValueType {
	return NewClassBuilder("Database").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			this.SetMember("_db", args[0], c)
			this.SetMember("_dialect", args[1], c)
		}).
		Method("m", func(c *Context, this ValueObject, args []Value) Value {
			var table ValueStr
			requiredArgs := args
			if len(args) > 1 {
				requiredArgs = args[:1]
			}
			EnsureFuncParams(c, "Database.m", requiredArgs,
				ArgRuleRequired{"table", TypeStr, &table},
			)
			rv := NewObjectAndInit(dbActiveRecordClass, c, this, table)
			if len(args) > 1 {
				c.InvokeMethod(rv, "and", Args(args[1:]...))
			}
			return rv
		}).
		Method("query", func(c *Context, this ValueObject, args []Value) Value {
			var (
				querySql ValueStr
			)
			EnsureFuncParams(c, "Database.query", args,
				ArgRuleRequired{"querySql", TypeStr, &querySql},
			)
			db := this.GetMember("_db", c)
			queryArgs := make([]interface{}, len(args)-1)
			for i := range queryArgs {
				queryArgs[i] = args[i+1].ToGoValue()
			}
			rows, err := db.ToGoValue().(*sql.DB).Query(querySql.Value(), queryArgs...)
			if err != nil {
				c.OnRuntimeError("Database.query query fail %s", err)
				return nil
			}
			return NewObjectAndInit(queryResultClass, c, NewGoValue(rows))
		}).
		Method("tables", func(c *Context, this ValueObject, args []Value) Value {
			db := this.GetMember("_db", c)
			dialect := this.GetMember("_dialect", c).ToGoValue().(dbDialect)
			rows, err := db.ToGoValue().(*sql.DB).Query(dialect.ShowTablesSQL())
			if err != nil {
				c.OnRuntimeError("Database.query query fail %s", err)
				return nil
			}
			res := NewObjectAndInit(queryResultClass, c, NewGoValue(rows))
			return c.InvokeMethod(res, "allOne", NoArgs)
		}).
		Method("execute", func(c *Context, this ValueObject, args []Value) Value {
			var (
				execSql ValueStr
			)
			EnsureFuncParams(c, "Database.execute", args,
				ArgRuleRequired{"querySql", TypeStr, &execSql},
			)
			db := this.GetMember("_db", c)
			execArgs := make([]interface{}, len(args)-1)
			for i := range execArgs {
				execArgs[i] = args[i+1].ToGoValue()
			}
			res, err := db.ToGoValue().(*sql.DB).Exec(execSql.Value(), execArgs...)
			if err != nil {
				c.OnRuntimeError("Database.execute execute fail %s", err)
				return nil
			}
			rv := NewObject()
			if affetcted, err := res.RowsAffected(); err != nil {
				c.OnRuntimeError("Database.execute get affected fail %s", err)
				return nil
			} else {
				rv.SetMember("affected", NewInt(affetcted), c)
			}
			if lastInsertID, err := res.LastInsertId(); err != nil {
				c.OnRuntimeError("Database.execute get lastInsertID fail %s", err)
				return nil
			} else {
				rv.SetMember("lastInsertID", NewInt(lastInsertID), c)
			}
			return rv
		}).
		Method("close", func(c *Context, this ValueObject, args []Value) Value {
			db := this.GetMember("_db", c).ToGoValue().(*sql.DB)
			if err := db.Close(); err != nil {
				c.OnRuntimeError("Database.close error %s", err)
				return nil
			}
			return Undefined()
		}).
		Method("newSession", func(c *Context, this ValueObject, args []Value) Value {
			db := this.GetMember("_db", c)
			return NewObjectAndInit(dbSessionClass, c, db, this.GetMember("_dialect", c))
		}).
		Method("atom", func(c *Context, this ValueObject, args []Value) Value {
			c.AssertArgNum(len(args), 1, 1, "Database.atom")
			fn := c.MustCallable(args[0])
			db := this.GetMember("_db", c)
			success := false
			session := NewObjectAndInit(dbSessionClass, c, db, this.GetMember("_dialect", c))
			defer func() {
				if success {
					c.InvokeMethod(session, "commit", NoArgs)
				} else {
					c.InvokeMethod(session, "rollback", NoArgs)
				}
			}()
			c.Invoke(fn, nil, Args(session))
			success = c.RetVal.IsTrue()
			return Undefined()
		}).
		Build()
}

func initDatabaseSessionClass() ValueType {
	globalSpId := int64(0)
	return NewClassBuilder("Session").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			parent := args[0].ToGoValue()
			var tx *sql.Tx
			var err error
			switch session := parent.(type) {
			case *sql.DB:
				{
					tx, err = session.Begin()
					if err != nil {
						c.OnRuntimeError("Session.__init__: begin transaction error %s", err)
					}
				}
			case *sql.Tx:
				{
					tx = session
					spId := fmt.Sprintf("_zggdb_sp_%d", atomic.AddInt64(&globalSpId, 1))
					_, err = tx.Exec(fmt.Sprintf("SAVEPOINT %s", spId))
					if err != nil {
						c.OnRuntimeError("Session.__init__: begin transaction error %s", err)
					}
					this.SetMember("__spId", NewStr(spId), c)
				}
			default:
				c.OnRuntimeError("unexpected session parent type")
			}
			this.SetMember("_tx", NewGoValue(tx), c)
			this.SetMember("_dialect", args[1], c)
		}).
		Method("m", func(c *Context, this ValueObject, args []Value) Value {
			var table ValueStr
			requiredArgs := args
			if len(args) > 1 {
				requiredArgs = args[:1]
			}
			EnsureFuncParams(c, "Session.m", requiredArgs,
				ArgRuleRequired{"table", TypeStr, &table},
			)
			rv := NewObjectAndInit(dbActiveRecordClass, c, this, table)
			if len(args) > 1 {
				c.InvokeMethod(rv, "and", Args(args[1:]...))
			}
			return rv
		}).
		Method("query", func(c *Context, this ValueObject, args []Value) Value {
			var (
				querySql ValueStr
			)
			EnsureFuncParams(c, "Session.query", args,
				ArgRuleRequired{"querySql", TypeStr, &querySql},
			)
			tx := this.GetMember("_tx", c).ToGoValue().(*sql.Tx)
			queryArgs := make([]interface{}, len(args)-1)
			for i := range queryArgs {
				queryArgs[i] = args[i+1].ToGoValue()
			}
			rows, err := tx.Query(querySql.Value(), queryArgs...)
			if err != nil {
				c.OnRuntimeError("Database.query query fail %s", err)
				return nil
			}
			return NewObjectAndInit(dbQueryResultClass, c, NewGoValue(rows))
		}).
		Method("execute", func(c *Context, this ValueObject, args []Value) Value {
			var (
				execSql ValueStr
			)
			EnsureFuncParams(c, "Session.execute", args,
				ArgRuleRequired{"querySql", TypeStr, &execSql},
			)
			tx := this.GetMember("_tx", c).ToGoValue().(*sql.Tx)
			execArgs := make([]interface{}, len(args)-1)
			for i := range execArgs {
				execArgs[i] = args[i+1].ToGoValue()
			}
			res, err := tx.Exec(execSql.Value(), execArgs...)
			if err != nil {
				c.OnRuntimeError("Session.execute execute fail %s", err)
				return nil
			}
			rv := NewObject()
			if affetcted, err := res.RowsAffected(); err != nil {
				c.OnRuntimeError("Session.execute get affected fail %s", err)
				return nil
			} else {
				rv.SetMember("affected", NewInt(affetcted), c)
			}
			if lastInsertID, err := res.LastInsertId(); err != nil {
				c.OnRuntimeError("Session.execute get lastInsertID fail %s", err)
				return nil
			} else {
				rv.SetMember("lastInsertID", NewInt(lastInsertID), c)
			}
			return rv
		}).
		Method("commit", func(c *Context, this ValueObject, args []Value) Value {
			tx := this.GetMember("_tx", c).ToGoValue().(*sql.Tx)
			if spId, ok := this.GetMember("__spId", c).(ValueStr); ok {
				if _, err := tx.Exec(fmt.Sprintf("RELEASE SAVEPOINT %s", spId.Value())); err != nil {
					c.OnRuntimeError("Session.commit: %s", err)
					return nil
				}
			} else {
				if err := tx.Commit(); err != nil {
					c.OnRuntimeError("Session.commit: %s", err)
					return nil
				}
			}
			return Undefined()
		}).
		Method("rollback", func(c *Context, this ValueObject, args []Value) Value {
			tx := this.GetMember("_tx", c).ToGoValue().(*sql.Tx)
			if spId, ok := this.GetMember("__spId", c).(ValueStr); ok {
				if _, err := tx.Exec(fmt.Sprintf("ROLLBACK TO %s", spId.Value())); err != nil {
					c.OnRuntimeError("Session.rollback: %s", err)
					return nil
				}
			} else {
				if err := tx.Rollback(); err != nil {
					c.OnRuntimeError("Session.rollback: %s", err)
					return nil
				}
			}
			return Undefined()
		}).
		Method("atom", func(c *Context, this ValueObject, args []Value) Value {
			c.AssertArgNum(len(args), 1, 1, "Session.atom")
			fn := c.MustCallable(args[0])
			tx := this.GetMember("_tx", c)
			success := false
			session := NewObjectAndInit(dbSessionClass, c, tx, this.GetMember("_dialect", c))
			defer func() {
				if success {
					c.InvokeMethod(session, "commit", NoArgs)
				} else {
					c.InvokeMethod(session, "rollback", NoArgs)
				}
			}()
			c.Invoke(fn, nil, Args(session))
			success = c.RetVal.IsTrue()
			return Undefined()
		}).
		Build()
}

func initDatabaseActiveRecordClass() ValueType {
	_buildWhere := func(c *Context, builder *strings.Builder, filters ValueArray) {
		if fn := filters.Len(); fn > 0 {
			for i := 0; i < fn; i++ {
				if i == 0 {
					builder.WriteString(" WHERE (")
				} else {
					builder.WriteString(" AND (")
				}
				f := filters.GetIndex(i, c)
				builder.WriteString(f.ToString(c))
				builder.WriteString(")")
			}
		}
	}
	_buildOrderBy := func(c *Context, builder *strings.Builder, orderBys ValueArray) {
		if fn := orderBys.Len(); fn > 0 {
			for i := 0; i < fn; i++ {
				if i == 0 {
					builder.WriteString(" ORDER BY ")
				} else {
					builder.WriteString(", ")
				}
				f := orderBys.GetIndex(i, c)
				builder.WriteString(f.ToString(c))
			}
		}
	}
	var _addFiltersByKV func(c *Context, dialect dbDialect, key string, val Value, filters, sqlArgs ValueArray)
	_addFiltersByKV = func(c *Context, dialect dbDialect, key string, val Value, filters, sqlArgs ValueArray) {
		if paramNum := strings.Count(key, "?"); paramNum > 0 {
			if varr, isArr := val.(ValueArray); isArr {
				if varr.Len() != paramNum {
					c.OnRuntimeError("db filter %s has %d placeholder(s), but get %d parameter(s)", key, paramNum, varr.Len())
				}
				filters.PushBack(NewStr(key))
				for i := 0; i < paramNum; i++ {
					sqlArgs.PushBack(varr.GetIndex(i, c))
				}
			} else {
				if paramNum != 1 {
					c.OnRuntimeError("db filter %s has %d placeholder(s), but get %d parameter(s)", key, paramNum, 1)
				}
				filters.PushBack(NewStr(key))
				sqlArgs.PushBack(val)
			}
		} else {
			switch v := val.(type) {
			case ValueUndefined:
				filters.PushBack(NewStr(key))
			case ValueArray:
				for i := 0; i < v.Len(); i++ {
					_addFiltersByKV(c, dialect, key, v.GetIndex(i, c), filters, sqlArgs)
				}
			case ValueCallable:
				c.Invoke(v, nil, Args(NewGoValue(dialect), NewStr(key), filters, sqlArgs))
			default:
				filters.PushBack(NewStr("%s = ?", dialect.Quote(key)))
				sqlArgs.PushBack(v)
			}
		}
	}
	_buildOffsetLimit := func(c *Context, builder *strings.Builder, this ValueObject, sqlArgs ValueArray) {
		switch offset := this.GetMember("_offset", c).(type) {
		case ValueInt:
			switch limit := this.GetMember("_limit", c).(type) {
			case ValueInt:
				builder.WriteString(" LIMIT ?, ?")
				sqlArgs.PushBack(offset)
				sqlArgs.PushBack(limit)
			case ValueUndefined:
				builder.WriteString(" LIMIT ?, 18446744073709551610")
				sqlArgs.PushBack(offset)
			}
		case ValueUndefined:
			switch limit := this.GetMember("_limit", c).(type) {
			case ValueInt:
				builder.WriteString(" LIMIT ?")
				sqlArgs.PushBack(limit)
			case ValueUndefined:
			}
		}

	}
	doQuery := func(c *Context, this ValueObject, args []Value) ValueObject {
		conn := c.MustObject(this.GetMember("conn", c))
		table := c.MustStr(this.GetMember("table", c))
		dialect := this.GetMember("_dialect", c).ToGoValue().(dbDialect)
		filters := c.MustArray(this.GetMember("filters", c))
		sqlArgs := c.MustArray(this.GetMember("sqlArgs", c))
		var sqlBuilder strings.Builder
		sqlBuilder.WriteString("SELECT")
		if len(args) > 0 {
			for i, arg := range args {
				field := arg.ToString(c)
				if !strings.ContainsAny(field, "`( ") {
					field = dialect.Quote(field)
				}
				if i > 0 {
					sqlBuilder.WriteString(fmt.Sprintf(", %s", field))
				} else {
					sqlBuilder.WriteString(fmt.Sprintf(" %s", field))
				}
			}
		} else {
			sqlBuilder.WriteString(" *")
		}
		sqlBuilder.WriteString(fmt.Sprintf(" FROM %s", dialect.Quote(table)))
		_buildWhere(c, &sqlBuilder, filters)
		_buildOrderBy(c, &sqlBuilder, c.MustArray(this.GetMember("orderBys", c)))
		_buildOffsetLimit(c, &sqlBuilder, this, sqlArgs)
		sql := sqlBuilder.String()
		if this.GetMember("_showSql", c).IsTrue() {
			fmt.Fprintln(c.Stdout, sql, sqlArgs.ToGoValue())
		}
		c.InvokeMethod(conn, "query", func() []Value {
			rv := make([]Value, sqlArgs.Len()+1)
			rv[0] = NewStr(sql)
			for i, a := range *sqlArgs.Values {
				rv[i+1] = a
			}
			return rv
		})
		return c.MustObject(c.RetVal)
	}
	return NewClassBuilder("ActiveRecord").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			var (
				conn  ValueObject
				table ValueStr
			)
			EnsureFuncParams(c, "ActiveRecord.__init__", args,
				ArgRuleRequired{"conn", TypeObject, &conn},
				ArgRuleRequired{"table", TypeStr, &table},
			)
			this.SetMember("conn", conn, c)
			this.SetMember("_dialect", conn.GetMember("_dialect", c), c)
			this.SetMember("table", table, c)
			this.SetMember("filters", NewArray(), c)
			this.SetMember("sqlArgs", NewArray(), c)
			this.SetMember("orderBys", NewArray(), c)
			this.SetMember("_showSql", NewBool(false), c)
		}).
		Method("showSql", func(c *Context, this ValueObject, args []Value) Value {
			if len(args) == 0 {
				this.SetMember("_showSql", NewBool(true), c)
			} else {
				this.SetMember("_showSql", NewBool(args[0].IsTrue()), c)
			}
			return this
		}).
		Method("and", func(c *Context, this ValueObject, args []Value) Value {
			if n := len(args); n > 0 {
				dialect := this.GetMember("_dialect", c).ToGoValue().(dbDialect)
				filters := c.MustArray(this.GetMember("filters", c))
				sqlArgs := c.MustArray(this.GetMember("sqlArgs", c))
				if filterObj, ok := args[0].(ValueObject); ok && n == 1 {
					filterObj.Iterate(func(k string, v Value) {
						_addFiltersByKV(c, dialect, k, v, filters, sqlArgs)
					})
				} else {
					filters.PushBack(args[0])
					for i := 1; i < n; i++ {
						sqlArgs.PushBack(args[i])
					}
				}
			}
			return this
		}).
		Method("asc", func(c *Context, this ValueObject, args []Value) Value {
			orderBys := c.MustArray(this.GetMember("orderBys", c))
			dialect := this.GetMember("_dialect", c).ToGoValue().(dbDialect)
			for _, arg := range args {
				field := arg.ToString(c)
				if !strings.ContainsAny(field, "`( ") {
					field = dialect.Quote(field)
				}
				orderBys.PushBack(NewStr("%s ASC", field))
			}
			return this
		}).
		Method("desc", func(c *Context, this ValueObject, args []Value) Value {
			orderBys := c.MustArray(this.GetMember("orderBys", c))
			dialect := this.GetMember("_dialect", c).ToGoValue().(dbDialect)
			for _, arg := range args {
				field := arg.ToString(c)
				if !strings.ContainsAny(field, "`( ") {
					field = dialect.Quote(field)
				}
				orderBys.PushBack(NewStr("%s DESC", field))
			}
			return this
		}).
		Method("limit", func(c *Context, this ValueObject, args []Value) Value {
			var v ValueInt
			EnsureFuncParams(c, "limit", args, ArgRuleRequired{"limit", TypeInt, &v})
			this.SetMember("_limit", v, c)
			return this
		}).
		Method("offset", func(c *Context, this ValueObject, args []Value) Value {
			var v ValueInt
			EnsureFuncParams(c, "offset", args, ArgRuleRequired{"offset", TypeInt, &v})
			this.SetMember("_offset", v, c)
			return this
		}).
		Method("find", func(c *Context, this ValueObject, args []Value) Value {
			res := doQuery(c, this, args)
			records := c.InvokeMethod(res, "all", Args(args...))
			c.InvokeMethod(res, "close", NoArgs)
			return records
		}).
		Method("findOne", func(c *Context, this ValueObject, args []Value) Value {
			this.SetMember("_limit", NewInt(1), c)
			res := doQuery(c, this, args)
			records := c.InvokeMethod(res, "all", Args(args...))
			c.InvokeMethod(res, "close", NoArgs)
			if arr, ok := records.(ValueArray); ok && arr.Len() > 0 {
				return arr.GetIndex(0, c)
			}
			return Nil()
		}).
		Method("one", func(c *Context, this ValueObject, args []Value) Value {
			return c.InvokeMethod(this, "findOne", Args(args...))
		}).
		Method("all", func(c *Context, this ValueObject, args []Value) Value {
			return c.InvokeMethod(this, "find", Args(args...))
		}).
		Method("count", func(c *Context, this ValueObject, args []Value) Value {
			var countField ValueStr
			EnsureFuncParams(c, "ActiveRecord.count", args,
				ArgRuleOptional{"countField", TypeStr, &countField, NewStr("(1)")},
			)
			args = []Value{countField}
			res := doQuery(c, this, args)
			records := c.InvokeMethod(res, "allArray", Args(args...))
			c.InvokeMethod(res, "close", NoArgs)
			if arr, ok := records.(ValueArray); ok && arr.Len() > 0 {
				return arr.GetIndex(0, c).(ValueArray).GetIndex(0, c)
			}
			return NewInt(0)
		}).
		Method("toTable", func(c *Context, this ValueObject, args []Value) Value {
			res := doQuery(c, this, args)
			records := c.InvokeMethod(res, "toTable", Args(args...))
			c.InvokeMethod(res, "close", NoArgs)
			return records
		}).
		Method("update", func(c *Context, this ValueObject, args []Value) Value {
			var updates ValueObject
			EnsureFuncParams(c, "ActiveRecord.update", args, ArgRuleRequired{"updates", TypeObject, &updates})
			conn := c.MustObject(this.GetMember("conn", c))
			dialect := this.GetMember("_dialect", c).ToGoValue().(dbDialect)
			table := c.MustStr(this.GetMember("table", c))
			filters := c.MustArray(this.GetMember("filters", c))
			sqlArgs := c.MustArray(this.GetMember("sqlArgs", c))
			realArgs := make([]Value, 1, 10+sqlArgs.Len())
			var sqlBuilder strings.Builder
			sqlBuilder.WriteString(fmt.Sprintf("UPDATE %s", dialect.Quote(table)))
			updateAdded := false
			updates.Iterate(func(k string, v Value) {
				if !updateAdded {
					sqlBuilder.WriteString(fmt.Sprintf(" SET %s = ?", dialect.Quote(k)))
					updateAdded = true
				} else {
					sqlBuilder.WriteString(fmt.Sprintf(", %s = ?", dialect.Quote(k)))
				}
				realArgs = append(realArgs, v)
			})
			if !updateAdded {
				return NewInt(0)
			}
			_buildWhere(c, &sqlBuilder, filters)
			_buildOrderBy(c, &sqlBuilder, c.MustArray(this.GetMember("orderBys", c)))
			_buildOffsetLimit(c, &sqlBuilder, this, sqlArgs)
			sql := sqlBuilder.String()
			realArgs[0] = NewStr(sql)
			realArgs = append(realArgs, *sqlArgs.Values...)
			if this.GetMember("_showSql", c).IsTrue() {
				fmt.Fprintln(c.Stdout, sql, NewArrayByValues(realArgs[1:]...).ToGoValue())
			}
			c.InvokeMethod(conn, "execute", Args(realArgs...))
			res := c.MustObject(c.RetVal)
			return res.GetMember("affected", c)
		}).
		Method("add", func(c *Context, this ValueObject, args []Value) Value {
			if len(args) < 1 {
				c.OnRuntimeError("requires at least 1 argument(s)")
			}
			table := c.MustStr(this.GetMember("table", c))
			dialect := this.GetMember("_dialect", c).ToGoValue().(dbDialect)
			var sqlBuilder strings.Builder
			sqlBuilder.WriteString("INSERT INTO ")
			sqlBuilder.WriteString(dialect.Quote(table))
			sqlBuilder.WriteString(" (")
			fieldMap := map[string]bool{}
			fields := []string{}
			for _, arg := range args {
				item := c.MustObject(arg)
				item.Iterate(func(key string, value Value) {
					if _, added := fieldMap[key]; !added {
						fieldMap[key] = true
						if len(fields) > 0 {
							sqlBuilder.WriteString(", ")
						}
						sqlBuilder.WriteString(dialect.Quote(key))
						fields = append(fields, key)
					}
				})
			}
			execArgs := make([]Value, 1+len(args)*len(fields))
			sqlBuilder.WriteString(")")
			p := 1
			for i, arg := range args {
				if i == 0 {
					sqlBuilder.WriteString(" VALUES (")
				} else {
					sqlBuilder.WriteString("), (")
				}
				for j, field := range fields {
					if j == 0 {
						sqlBuilder.WriteString("?")
					} else {
						sqlBuilder.WriteString(", ?")
					}
					execArgs[p] = arg.GetMember(field, c)
					p++
				}
			}
			sqlBuilder.WriteString(")")
			sql := sqlBuilder.String()
			execArgs[0] = NewStr(sql)
			if this.GetMember("_showSql", c).IsTrue() {
				fmt.Fprintln(c.Stdout, sql, NewArrayByValues(execArgs[1:]...).ToGoValue())
			}
			conn := c.MustObject(this.GetMember("conn", c))
			return c.InvokeMethod(conn, "execute", Args(execArgs...))
		}).
		Build()
}

type _dbopDef struct {
	op      string
	pattern string
	n       int
}

func libDbOp(*Context) ValueObject {
	lib := NewObject()
	defs := []_dbopDef{
		{"gt", "%s > ?", 1},
		{"ge", "%s >= ?", 1},
		{"lt", "%s < ?", 1},
		{"le", "%s <= ?", 1},
		{"eq", "%s = ?", 1},
		{"ne", "%s <> ?", 1},
	}
	getArgs := func(c *Context, name string, args []Value) (dbDialect, ValueStr, ValueArray, ValueArray) {
		var (
			dialect GoValue
			field   ValueStr
			filters ValueArray
			sqlArgs ValueArray
		)
		EnsureFuncParams(c, name, args,
			ArgRuleRequired{"dialect", TypeGoValue, &dialect},
			ArgRuleRequired{"field", TypeStr, &field},
			ArgRuleRequired{"filters", TypeArray, &filters},
			ArgRuleRequired{"sqlArgs", TypeArray, &sqlArgs},
		)
		return dialect.ToGoValue().(dbDialect), field, filters, sqlArgs
	}
	for _, dd := range defs {
		d := dd
		lib.SetMember(d.op, NewNativeFunction(d.op, func(c *Context, this Value, fargs []Value) Value {
			c.AssertArgNum(len(fargs), d.n, d.n, d.op)
			return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
				dialect, field, filters, sqlArgs := getArgs(c, d.op, args)
				filters.PushBack(NewStr(d.pattern, dialect.Quote(field.Value())))
				for _, fa := range fargs {
					sqlArgs.PushBack(fa)
				}
				return Undefined()
			})
		}), nil)
	}
	lib.SetMember("in_", NewNativeFunction("in_", func(c *Context, this Value, fargs []Value) Value {
		if len(fargs) == 0 {
			return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
				_, _, filters, _ := getArgs(c, "in_", args)
				filters.PushBack(NewStr("1 = 2"))
				return Undefined()
			})
		}
		return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
			dialect, field, filters, sqlArgs := getArgs(c, "in_", args)
			var fb strings.Builder
			fb.WriteString(fmt.Sprintf("%s in (?", dialect.Quote(field.Value())))
			sqlArgs.PushBack(fargs[0])
			for i := 1; i < len(fargs); i++ {
				fb.WriteString(", ?")
				sqlArgs.PushBack(fargs[i])
			}
			fb.WriteRune(')')
			filters.PushBack(NewStr(fb.String()))
			return Undefined()
		})
	}), nil)
	lib.SetMember("contains", NewNativeFunction("contains", func(c *Context, this Value, fargs []Value) Value {
		var subs ValueStr
		EnsureFuncParams(c, "contains", fargs, ArgRuleRequired{"subs", TypeStr, &subs})
		return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
			dialect, field, filters, sqlArgs := getArgs(c, "contains", args)
			filters.PushBack(NewStr("%s like ?", dialect.Quote(field.ToString(c))))
			sqlArgs.PushBack(NewStr("%%%s%%", subs.Value()))
			return Undefined()
		})
	}), nil)
	lib.SetMember("startsWith", NewNativeFunction("startsWith", func(c *Context, this Value, fargs []Value) Value {
		var subs ValueStr
		EnsureFuncParams(c, "startsWith", fargs, ArgRuleRequired{"subs", TypeStr, &subs})
		return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
			dialect, field, filters, sqlArgs := getArgs(c, "startsWith", args)
			filters.PushBack(NewStr("%s like ?", dialect.Quote(field.ToString(c))))
			sqlArgs.PushBack(NewStr("%s%%", subs.Value()))
			return Undefined()
		})
	}), nil)
	lib.SetMember("endsWith", NewNativeFunction("endsWith", func(c *Context, this Value, fargs []Value) Value {
		var subs ValueStr
		EnsureFuncParams(c, "endsWith", fargs, ArgRuleRequired{"subs", TypeStr, &subs})
		return NewNativeFunction("", func(c *Context, this Value, args []Value) Value {
			dialect, field, filters, sqlArgs := getArgs(c, "endsWith", args)
			filters.PushBack(NewStr("%s like ?", dialect.Quote(field.ToString(c))))
			sqlArgs.PushBack(NewStr("%%%s", subs.Value()))
			return Undefined()
		})
	}), nil)
	return lib
}

func init() {
	dbSessionClass = initDatabaseSessionClass()
	dbActiveRecordClass = initDatabaseActiveRecordClass()
}
