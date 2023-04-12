package builtin_libs

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/samber/lo"
	. "github.com/zgg-lang/zgg-go/runtime"

	runewidth "github.com/mattn/go-runewidth"
)

var (
	ptablePTableClass ValueType
)

func libPtable(c *Context) ValueObject {
	lib := NewObject()
	lib.SetMember("PTable", ptablePTableClass, c)
	lib.SetMember("__call__", ptablePTableClass, c)
	lib.SetMember("fromCsvFile", ptableFromCsvFile, c)
	lib.SetMember("fromCsv", ptableFromCsv, c)
	lib.SetMember("query", NewNativeFunction("ptable.query", func(c *Context, this Value, args []Value) Value {
		return c.InvokeMethod(ptablePTableClass, "query", Args(args...))
	}), c)
	return lib
}

type ptableMeta struct {
	headers    []string
	colFormats []string
}

type ptableAsciiCellInfo struct {
	lines  []string
	widths []int
}

type ptableTextChars struct {
	topLeft     string
	top         string
	topMid      string
	topRight    string
	sepLeft     string
	sep         string
	sepMid      string
	sepRight    string
	dataLeft    string
	dataMid     string
	dataRight   string
	bottomLeft  string
	bottom      string
	bottomMid   string
	bottomRight string
}

var (
	ptableAsciiChars = ptableTextChars{
		topLeft:     "+-",
		top:         "-",
		topMid:      "-+-",
		topRight:    "-+",
		sepLeft:     "+-",
		sep:         "-",
		sepMid:      "-+-",
		sepRight:    "-+",
		dataLeft:    "| ",
		dataMid:     " | ",
		dataRight:   " |",
		bottomLeft:  "+-",
		bottom:      "-",
		bottomMid:   "-+-",
		bottomRight: "-+",
	}
	ptableUnicodeChars = ptableTextChars{
		topLeft:     "┏━",
		top:         "━",
		topMid:      "━┯━",
		topRight:    "━┓",
		sepLeft:     "┣━",
		sep:         "━",
		sepMid:      "━┿━",
		sepRight:    "━┫",
		dataLeft:    "┃ ",
		dataMid:     " │ ",
		dataRight:   " ┃",
		bottomLeft:  "┗━",
		bottom:      "━",
		bottomMid:   "━┷━",
		bottomRight: "━┛",
	}
)

type ptableTableInstance struct {
	name     string
	instance ValueObject
}

func initPTableClass() {
	filter := func(c *Context, filterFn ValueCallable, row []Value) bool {
		if filterFn == nil {
			return true
		}
		filterArg := NewArrayByValues(row...)
		c.Invoke(filterFn, nil, Args(filterArg))
		return c.RetVal.IsTrue()
	}
	formatText := func(c *Context, formatFn ValueCallable, meta *ptableMeta, i, j int, item Value) string {
		if formatFn != nil {
			c.Invoke(formatFn, nil, Args(item, NewInt(int64(i)), NewInt(int64(j))))
			return c.RetVal.ToString(c)
		} else if j < len(meta.colFormats) && meta.colFormats[j] != "" {
			return fmt.Sprintf(meta.colFormats[j], item.ToGoValue())
		} else {
			return item.ToString(c)
		}
	}
	getAlign := func(c *Context, alignConf Value, v Value, i, j int, defaultAlign string) string {
		switch ac := alignConf.(type) {
		case ValueCallable:
			c.Invoke(ac, nil, Args(v, NewInt(int64(i)), NewInt(int64(j))))
			return c.RetVal.ToString(c)
		case ValueStr:
			return ac.Value()
		case ValueArray:
			return ac.GetIndex(j, c).ToString(c)
		}
		return defaultAlign
	}
	getColor := func(c *Context, colorConf, v Value, i, j int) string {
		switch ac := colorConf.(type) {
		case ValueCallable:
			c.Invoke(ac, nil, Args(v, NewInt(int64(i)), NewInt(int64(j))))
			return c.RetVal.ToString(c)
		case ValueStr:
			return ac.Value()
		case ValueArray:
			return ac.GetIndex(j, c).ToString(c)
		}
		return ""
	}
	withColor := func(color, text string) string {
		if color == "" {
			return text
		}
		return fmt.Sprintf("\x1b[%s;1m%s\x1b[0m", color, text)
	}
	renderAscii := func(c *Context, this ValueObject, formatFn, filterFn ValueCallable, alignConf, colorConf Value, isMarkdown bool, tchars *ptableTextChars) Value {
		_meta := this.GetMember("_meta", c).ToGoValue().(*ptableMeta)
		_rows := this.GetMember("_rows", c).ToGoValue().([][]Value)
		rowsText := make([][]ptableAsciiCellInfo, 0, len(_rows))
		colMaxWidths := make([]int, len(_meta.headers))
		// init max widths
		for i, header := range _meta.headers {
			colMaxWidths[i] = runewidth.StringWidth(header)
		}
		rowLines := make([]int, 0, len(_rows))
		for i, row := range _rows {
			if !filter(c, filterFn, row) {
				continue
			}
			rowText := make([]ptableAsciiCellInfo, len(row))
			n := 0
			for j, item := range row {
				itemText := formatText(c, formatFn, _meta, i, j, item)
				var lines []string
				if isMarkdown {
					lines = []string{itemText}
				} else {
					lines = strings.Split(itemText, "\n")
				}
				ln := len(lines)
				if ln > n {
					n = ln
				}
				widths := make([]int, ln)
				for k, line := range lines {
					w := runewidth.StringWidth(line)
					widths[k] = w
					if j >= len(colMaxWidths) {
						colMaxWidths = append(colMaxWidths, w)
					} else if colMaxWidths[j] < w {
						colMaxWidths[j] = w
					}
				}
				rowText[j] = ptableAsciiCellInfo{lines: lines, widths: widths}
			}
			rowsText = append(rowsText, rowText)
			rowLines = append(rowLines, n)
		}
		var b strings.Builder
		if !IsUndefined(colorConf) {
			b.WriteString("\x1b[0m")
		}
		// head line
		if !isMarkdown {
			b.WriteString(tchars.topLeft)
			for j, w := range colMaxWidths {
				if j > 0 {
					b.WriteString(tchars.topMid)
				}
				b.WriteString(strings.Repeat(tchars.top, w))
			}
			b.WriteString(tchars.topRight)
			// headers
			if len(_meta.headers) > 0 {
				b.WriteString("\n")
				b.WriteString(tchars.dataLeft)
				for j, w := range colMaxWidths {
					if j > 0 {
						b.WriteString(tchars.dataMid)
					}
					var header string
					if j < len(_meta.headers) {
						header = _meta.headers[j]
					}
					b.WriteString(header)
					b.WriteString(strings.Repeat(" ", w-runewidth.StringWidth(header)))
				}
				b.WriteString(tchars.dataRight)
				b.WriteString("\n")
				b.WriteString(tchars.sepLeft)
				for j, w := range colMaxWidths {
					if j > 0 {
						b.WriteString(tchars.sepMid)
					}
					b.WriteString(strings.Repeat(tchars.sep, w))
				}
				b.WriteString(tchars.sepRight)
			}
		} else {
			if len(_meta.headers) > 0 {
				b.WriteString("| ")
				for j, w := range colMaxWidths {
					if j > 0 {
						b.WriteString(" | ")
					}
					var header string
					if j < len(_meta.headers) {
						header = _meta.headers[j]
					}
					b.WriteString(header)
					b.WriteString(strings.Repeat(" ", w-runewidth.StringWidth(header)))
				}
				b.WriteString(" |\n|")
				for j, w := range colMaxWidths {
					var align string
					if len(_rows) > 0 {
						align = getAlign(c, alignConf, _rows[0][j], 0, j, "")
					}
					l, r := " ", " "
					switch align {
					case "l":
						l = ":"
					case "r":
						r = ":"
					case "m":
						l, r = ":", ":"
					}
					if j > 0 {
						b.WriteString("|")
					}
					b.WriteString(l)
					b.WriteString(strings.Repeat("-", w))
					b.WriteString(r)
				}
				b.WriteString("|")
			}
		}
		for i, row := range rowsText {
			for k := 0; k < rowLines[i]; k++ {
				b.WriteString("\n")
				b.WriteString(tchars.dataLeft)
				for j, item := range row {
					if j > 0 {
						b.WriteString(tchars.dataMid)
					}
					w := colMaxWidths[j]
					if k >= len(item.lines) {
						b.WriteString(strings.Repeat(" ", w))
					} else {
						// data rows
						defaultAlign := "l"
						if !isMarkdown {
							switch _rows[i][j].(type) {
							case ValueInt, ValueFloat:
								defaultAlign = "r"
							}
						}
						align := getAlign(c, alignConf, _rows[i][j], i, j, defaultAlign)
						color := getColor(c, colorConf, _rows[i][j], i, j)
						line := withColor(color, item.lines[k])
						lw := item.widths[k]
						switch align {
						case "l":
							b.WriteString(line)
							b.WriteString(strings.Repeat(" ", w-lw))
						case "m":
							{
								sl := (w - lw) / 2
								b.WriteString(strings.Repeat(" ", sl))
								b.WriteString(line)
								b.WriteString(strings.Repeat(" ", (w-lw)-sl))
							}
						default:
							b.WriteString(strings.Repeat(" ", w-lw))
							b.WriteString(line)
						}
					}
				}
				b.WriteString(tchars.dataRight)
			}
		}
		// bottom line
		if !isMarkdown {
			b.WriteString("\n")
			b.WriteString(tchars.bottomLeft)
			for j, w := range colMaxWidths {
				if j > 0 {
					b.WriteString(tchars.bottomMid)
				}
				b.WriteString(strings.Repeat(tchars.bottom, w))
			}
			b.WriteString(tchars.bottomRight)
		}
		return NewStr(b.String())
	}
	var (
		tableNamesInSelectRegexp = regexp.MustCompile(`(?i)(?:FROM|JOIN)\s+(\w+(?:\s*,\s*\w+)*)`)
		tableNamesInBlock        = regexp.MustCompile(`\w+`)
	)
	findTables := func(c *Context, querySQL string) []ptableTableInstance {
		rv := make([]ptableTableInstance, 0)
		tableMap := map[string]ValueObject{}
		for _, m := range tableNamesInSelectRegexp.FindAllStringSubmatch(querySQL, -1) {
			for _, names := range m[1:] {
				for _, name := range tableNamesInBlock.FindAllString(names, -1) {
					if _, found := tableMap[name]; found {
						continue
					}
					v, exists := c.FindValue(name)
					if !exists {
						c.RaiseRuntimeError("cannot find table %s in context", name)
					}
					if !v.Type().IsSubOf(ptablePTableClass) {
						c.RaiseRuntimeError("variable %s is not a PTable", name)
					}
					tableMap[name] = v.(ValueObject)
				}
			}
		}
		for name, instance := range tableMap {
			rv = append(rv, ptableTableInstance{name: name, instance: instance})
		}
		return rv
	}
	addTableToDB := func(c *Context, tmpDB *sql.DB, table ValueObject, tableName string) {
		var (
			meta              = table.GetMember("_meta", c).ToGoValue().(*ptableMeta)
			rows              = table.GetMember("_rows", c).ToGoValue().([][]Value)
			createStmtBuilder strings.Builder
			insertStmtBuilder strings.Builder
		)
		createStmtBuilder.WriteString("CREATE TABLE ")
		createStmtBuilder.WriteString(tableName)
		createStmtBuilder.WriteString(" (")
		insertStmtBuilder.WriteString("INSERT INTO ")
		insertStmtBuilder.WriteString(tableName)
		insertStmtBuilder.WriteString(" VALUES (")
		for i, h := range meta.headers {
			if i >= len(rows[0]) {
				c.RaiseRuntimeError("PTable.query: not enough values in first row")
			}
			dbType := "INT"
			if len(rows) > 0 {
				v := rows[0][i]
				switch v.(type) {
				case ValueInt:
					dbType = "INT"
				case ValueStr:
					dbType = "VARCHAR(1024)"
				case ValueFloat:
					dbType = "REAL"
				case ValueBool:
					dbType = "TINYINT"
				default:
					c.RaiseRuntimeError("PTable.query: unsupported value type %s",
						v.Type().Name)
				}
			}
			if i > 0 {
				createStmtBuilder.WriteString(", ")
				insertStmtBuilder.WriteString(", ?")
			} else {
				insertStmtBuilder.WriteString("?")
			}
			createStmtBuilder.WriteRune('`')
			createStmtBuilder.WriteString(h)
			createStmtBuilder.WriteString("` ")
			createStmtBuilder.WriteString(dbType)
		}
		createStmtBuilder.WriteString(")")
		insertStmtBuilder.WriteString(")")
		if _, err := tmpDB.Exec(createStmtBuilder.String()); err != nil {
			c.RaiseRuntimeError("PTable.query: create temp table error %+v", err)
		}
		var (
			insertValues = make([]interface{}, len(meta.headers))
			insertSQL    = insertStmtBuilder.String()
		)
		for rowIndex, row := range rows {
			for i := range insertValues {
				insertValues[i] = row[i].ToGoValue()
			}
			if _, err := tmpDB.Exec(insertSQL, insertValues...); err != nil {
				c.RaiseRuntimeError("PTable.query: insert values in row %d error %+v",
					rowIndex, err)
			}
		}
	}
	query := func(c *Context, querySQL string, queryArgs []any) Value {
		tmpDB, err := sql.Open("sqlite", ":memory:")
		if err != nil {
			c.RaiseRuntimeError("PTable.query: open temp database error %+v", err)
		}
		defer tmpDB.Close()
		for _, table := range findTables(c, querySQL) {
			addTableToDB(c, tmpDB, table.instance, table.name)
		}
		retRows, err := tmpDB.Query(querySQL, queryArgs...)
		if err != nil {
			c.RaiseRuntimeError("PTable.query: query error %+v", err)
		}
		defer func() {
			if err := retRows.Close(); err != nil {
				c.RaiseRuntimeError("PTable.query: query close error %+v", err)
			}
		}()
		colTypes, err := retRows.ColumnTypes()
		if err != nil {
			c.RaiseRuntimeError("PTable.query: read rows get column types error %+v", err)
		}
		cols, err := retRows.Columns()
		if err != nil {
			c.RaiseRuntimeError("PTable.query: read rows get columns error %+v", err)
		}
		rv := NewObjectAndInit(ptablePTableClass, c, lo.Map(cols, func(n string, _ int) Value {
			return NewStr(n)
		})...)
		rvRows := make([][]Value, 0)
		for retRows.Next() {
			row := *(dbScanRowsToArray(c, retRows, colTypes, cols).Values)
			rvRows = append(rvRows, row)
		}
		rv.SetMember("_rows", NewGoValue(rvRows), c)
		return rv

	}
	ptablePTableClass = NewClassBuilder("PTable").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			_meta := &ptableMeta{
				headers: make([]string, len(args)),
			}
			for i, arg := range args {
				header := arg.ToString(c)
				_meta.headers[i] = header
			}
			this.SetMember("_meta", NewGoValue(_meta), c)
			_rows := make([][]Value, 0)
			this.SetMember("_rows", NewGoValue(_rows), c)
		}).
		Method("add", func(c *Context, this ValueObject, args []Value) Value {
			_rows := this.GetMember("_rows", c).ToGoValue().([][]Value)
			_rows = append(_rows, args)
			this.SetMember("_rows", NewGoValue(_rows), c)
			return this
		}).
		Method("addArray", func(c *Context, this ValueObject, args []Value) Value {
			_rows := this.GetMember("_rows", c).ToGoValue().([][]Value)
			var arr ValueArray
			EnsureFuncParams(c, "addArray", args, ArgRuleRequired("array", TypeArray, &arr))
			for i := 0; i < arr.Len(); i++ {
				row := arr.GetIndex(i, c)
				if rowArr, is := row.(ValueArray); !is {
					c.RaiseRuntimeError("PTable.addArray: every array items must be an array")
				} else {
					_rows = append(_rows, *rowArr.Values)
				}
			}
			this.SetMember("_rows", NewGoValue(_rows), c)
			return this
		}).
		Method("addObject", func(c *Context, this ValueObject, args []Value) Value {
			if len(args) < 1 {
				return this
			}
			o := args[0]
			_rows := this.GetMember("_rows", c).ToGoValue().([][]Value)
			var row []Value
			if len(args) == 1 {
				_meta := this.GetMember("_meta", c).ToGoValue().(*ptableMeta)
				row = make([]Value, len(_meta.headers))
				for i, h := range _meta.headers {
					row[i] = o.GetMember(h, c)
				}
			} else {
				row = make([]Value, len(args)-1)
				for i := range row {
					row[i] = o.GetMember(args[i+1].ToString(c), c)
				}
			}
			_rows = append(_rows, row)
			this.SetMember("_rows", NewGoValue(_rows), c)
			return this
		}).
		Method("addObjects", func(c *Context, this ValueObject, args []Value) Value {
			arr := c.MustArray(args[0])
			_rows := this.GetMember("_rows", c).ToGoValue().([][]Value)
			for i := 0; i < arr.Len(); i++ {
				o := arr.GetIndex(i, c)
				row := make([]Value, len(args)-1)
				for j := range row {
					row[j] = o.GetMember(args[j+1].ToString(c), c)
				}
				_rows = append(_rows, row)
			}
			this.SetMember("_rows", NewGoValue(_rows), c)
			return this
		}).
		Method("colFormat", func(c *Context, this ValueObject, args []Value) Value {
			var (
				col    ValueInt
				format ValueStr
			)
			EnsureFuncParams(c, "colFormat", args,
				ArgRuleRequired("col", TypeInt, &col),
				ArgRuleRequired("foramt", TypeStr, &format),
			)
			_meta := this.GetMember("_meta", c).ToGoValue().(*ptableMeta)
			for i := len(_meta.colFormats); i < col.AsInt()+1; i++ {
				_meta.colFormats = append(_meta.colFormats, "")
			}
			_meta.colFormats[col.AsInt()] = format.Value()
			return this
		}).
		Method("ascii", func(c *Context, this ValueObject, args []Value) Value {
			var (
				formatFn  ValueCallable
				filterFn  ValueCallable
				alignConf Value
				colorConf Value
			)
			if len(args) > 0 {
				opt := args[0]
				if f, ok := c.GetCallable(opt.GetMember("formatter", c)); ok {
					formatFn = f
				}
				if f, ok := c.GetCallable(opt.GetMember("filter", c)); ok {
					filterFn = f
				}
				alignConf = opt.GetMember("align", c)
				colorConf = opt.GetMember("color", c)
			}
			return renderAscii(c, this, formatFn, filterFn, alignConf, colorConf, false, &ptableAsciiChars)
		}).
		Method("txt", func(c *Context, this ValueObject, args []Value) Value {
			return c.InvokeMethod(this, "ascii", Args(args...))
		}).
		Method("text", func(c *Context, this ValueObject, args []Value) Value {
			return c.InvokeMethod(this, "ascii", Args(args...))
		}).
		Method("unicode", func(c *Context, this ValueObject, args []Value) Value {
			var (
				formatFn  ValueCallable
				filterFn  ValueCallable
				alignConf Value
				colorConf Value
			)
			if len(args) > 0 {
				opt := args[0]
				if f, ok := c.GetCallable(opt.GetMember("formatter", c)); ok {
					formatFn = f
				}
				if f, ok := c.GetCallable(opt.GetMember("filter", c)); ok {
					filterFn = f
				}
				alignConf = opt.GetMember("align", c)
				colorConf = opt.GetMember("color", c)
			}
			return renderAscii(c, this, formatFn, filterFn, alignConf, colorConf, false, &ptableUnicodeChars)
		}).
		Method("markdown", func(c *Context, this ValueObject, args []Value) Value {
			var (
				formatFn  ValueCallable
				filterFn  ValueCallable
				alignConf Value
				colorConf Value
			)
			if len(args) > 0 {
				opt := args[0]
				if f, ok := c.GetCallable(opt.GetMember("formatter", c)); ok {
					formatFn = f
				}
				if f, ok := c.GetCallable(opt.GetMember("filter", c)); ok {
					filterFn = f
				}
				alignConf = opt.GetMember("align", c)
				colorConf = opt.GetMember("color", c)
			}
			return renderAscii(c, this, formatFn, filterFn, alignConf, colorConf, true, &ptableAsciiChars)
		}).
		Method("md", func(c *Context, this ValueObject, args []Value) Value {
			return c.InvokeMethod(this, "markdown", Args(args...))
		}).
		Method("html", func(c *Context, this ValueObject, args []Value) Value {
			var b strings.Builder
			_meta := this.GetMember("_meta", c).ToGoValue().(*ptableMeta)
			_rows := this.GetMember("_rows", c).ToGoValue().([][]Value)
			var (
				formatFn ValueCallable
				filterFn ValueCallable
			)
			// var alignConf Value
			if len(args) > 0 {
				if f, ok := c.GetCallable(args[0].GetMember("formatter", c)); ok {
					formatFn = f
				}
				if f, ok := c.GetCallable(args[0].GetMember("filter", c)); ok {
					filterFn = f
				}
				// alignConf = args[0].GetMember("align", c)
			}
			colNum := len(_meta.headers)
			rowsTexts := make([][]string, 0, len(_rows))
			for i, row := range _rows {
				if !filter(c, filterFn, row) {
					continue
				}
				l := len(row)
				if l > colNum {
					colNum = l
				}
				rowTexts := make([]string, l)
				for j, item := range row {
					rowTexts[j] = formatText(c, formatFn, _meta, i, j, item)
				}
				rowsTexts = append(rowsTexts, rowTexts)
			}
			b.WriteString("<TABLE>\n")
			if len(_meta.headers) > 0 {
				b.WriteString("  <THEAD>\n")
				b.WriteString("    <TR>\n")
				for _, h := range _meta.headers {
					b.WriteString("      <TH>")
					b.WriteString(html.EscapeString(h))
					b.WriteString("</TH>\n")
				}
				b.WriteString("    </TR>\n")
				b.WriteString("  </THEAD>\n")
			}
			if len(rowsTexts) > 0 {
				b.WriteString("  <TBODY>\n")
				for _, row := range rowsTexts {
					b.WriteString("    <TR>\n")
					for _, item := range row {
						b.WriteString("      <TD>")
						b.WriteString(html.EscapeString(item))
						b.WriteString("</TD>\n")
					}
					b.WriteString("    </TR>\n")
				}
				b.WriteString("  </TBODY>\n")
			}
			b.WriteString("</TABLE>\n")
			return NewStr(b.String())
		}).
		Method("csv", func(c *Context, this ValueObject, args []Value) Value {
			_meta := this.GetMember("_meta", c).ToGoValue().(*ptableMeta)
			_rows := this.GetMember("_rows", c).ToGoValue().([][]Value)
			var (
				formatFn ValueCallable
				filterFn ValueCallable
			)
			// var alignConf Value
			if len(args) > 0 {
				if f, ok := c.GetCallable(args[0].GetMember("formatter", c)); ok {
					formatFn = f
				}
				if f, ok := c.GetCallable(args[0].GetMember("filter", c)); ok {
					filterFn = f
				}
			}
			colNum := len(_meta.headers)
			rowsTexts := make([][]string, 0, len(_rows))
			for i, row := range _rows {
				if !filter(c, filterFn, row) {
					continue
				}
				l := len(row)
				if l > colNum {
					colNum = l
				}
				rowTexts := make([]string, l)
				for j, item := range row {
					rowTexts[j] = formatText(c, formatFn, _meta, i, j, item)
				}
				rowsTexts = append(rowsTexts, rowTexts)
			}
			var b strings.Builder
			w := csv.NewWriter(&b)
			if len(_meta.headers) > 0 {
				w.Write(_meta.headers)
			}
			if len(_rows) > 0 {
				w.WriteAll(rowsTexts)
			}
			w.Flush()
			return NewStr(b.String())
		}).
		Method("toArray", func(c *Context, this ValueObject, args []Value) Value {
			var (
				includHeaders ValueBool
			)
			EnsureFuncParams(c, "PTable.toArray", args,
				ArgRuleOptional("includHeaders", TypeBool, &includHeaders, NewBool(false)),
			)
			incH := includHeaders.Value()
			_meta := this.GetMember("_meta", c).ToGoValue().(*ptableMeta)
			_rows := this.GetMember("_rows", c).ToGoValue().([][]Value)
			rowNum := len(_rows)
			if incH {
				rowNum++
			}
			rv := NewArray(rowNum)
			if incH {
				header := NewArray(len(_meta.headers))
				for _, h := range _meta.headers {
					header.PushBack(NewStr(h))
				}
				rv.PushBack(header)
			}
			for _, row := range _rows {
				rv.PushBack(NewArrayByValues(row...))
			}
			return rv
		}, "includeHeaders").
		Method("query", func(c *Context, this ValueObject, args []Value) Value {
			if len(args) < 1 {
				c.RaiseRuntimeError("PTable.query: requires at least 1 argument")
			}
			var (
				querySQL  = args[0].ToString(c)
				queryArgs = lo.Map(args[1:], func(v Value, _ int) any {
					return v.ToGoValue()
				})
			)
			c.PushStack()
			defer c.PopStack()
			c.SetLocalValue("this", this)
			return query(c, querySQL, queryArgs)
		}).
		StaticMethod("query", func(c *Context, this Value, args []Value) Value {
			if len(args) < 1 {
				c.RaiseRuntimeError("PTable.query: requires at least 1 argument")
			}
			var (
				querySQL  = args[0].ToString(c)
				queryArgs = lo.Map(args[1:], func(v Value, _ int) any {
					return v.ToGoValue()
				})
			)
			return query(c, querySQL, queryArgs)
		}).
		Method("__str__", func(c *Context, this ValueObject, args []Value) Value {
			return c.InvokeMethod(this, "ascii", Args(args...))
		}).
		Build()
}

var ptableFromCsvFile = NewNativeFunction("ptable.fromCsvFile", func(c *Context, this Value, args []Value) Value {
	var (
		filename ValueStr
	)
	EnsureFuncParams(c, "ptable.fromCsvFile", args,
		ArgRuleRequired("filename", TypeStr, &filename),
	)
	name := filename.Value()
	f, err := os.Open(name)
	if err != nil {
		c.RaiseRuntimeError("open csv file %s error %s", name, err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		c.RaiseRuntimeError("read csv file %s error %s", name, err)
	}
	if len(rows) == 0 {
		c.RaiseRuntimeError("read csv file %s empty", name)
	}
	tableArgs := make([]Value, 0, len(rows[0]))
	for _, v := range rows[0] {
		tableArgs = append(tableArgs, NewStr(v))
	}
	rv := NewObjectAndInit(ptablePTableClass, c, tableArgs...)
	for i := 1; i < len(rows); i++ {
		tableArgs = make([]Value, 0, len(rows[i]))
		for _, v := range rows[i] {
			tableArgs = append(tableArgs, NewStr(v))
		}
		c.InvokeMethod(rv, "add", Args(tableArgs...))
	}
	return rv
})

var ptableFromCsv = NewNativeFunction("ptable.fromCsv", func(c *Context, this Value, args []Value) Value {
	var (
		argFilename ValueStr
		argUrl      ValueStr
		argContent  ValueStr
	)
	EnsureFuncParams(c, "ptable.fromCsv", args,
		ArgRuleOptional("filename", TypeStr, &argFilename, NewStr("")),
		ArgRuleOptional("url", TypeStr, &argUrl, NewStr("")),
		ArgRuleOptional("content", TypeStr, &argContent, NewStr("")),
	)
	var contentReader io.Reader
	if content := argContent.Value(); content != "" {
		contentReader = strings.NewReader(content)
	} else if url := argUrl.Value(); url != "" {
		r, err := http.Get(url)
		if err != nil {
			c.RaiseRuntimeError("load csv from url %s error %s", url, err)
		}
		defer r.Body.Close()
		contentReader = r.Body
	} else if name := argFilename.Value(); name != "" {
		f, err := os.Open(name)
		if err != nil {
			c.RaiseRuntimeError("open csv file %s error %s", name, err)
		}
		defer f.Close()
		contentReader = f
	}
	rows, err := csv.NewReader(contentReader).ReadAll()
	if err != nil {
		c.RaiseRuntimeError("read csv error %s", err)
	}
	if len(rows) == 0 {
		c.RaiseRuntimeError("read csv empty")
	}
	tableArgs := make([]Value, 0, len(rows[0]))
	for _, v := range rows[0] {
		tableArgs = append(tableArgs, NewStr(v))
	}
	rv := NewObjectAndInit(ptablePTableClass, c, tableArgs...)
	for i := 1; i < len(rows); i++ {
		tableArgs = make([]Value, 0, len(rows[i]))
		for _, v := range rows[i] {
			tableArgs = append(tableArgs, NewStr(v))
		}
		c.InvokeMethod(rv, "add", Args(tableArgs...))
	}
	return rv
}, "filename", "url", "content")

func init() {
	initPTableClass()
}
