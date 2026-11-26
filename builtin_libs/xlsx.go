package builtin_libs

import (
	"io"
	"strconv"

	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
	. "github.com/zgg-lang/zgg-go/runtime"
)

var (
	xlsxFileClass ValueType
)

func libXlsx(c *Context) ValueObject {
	lib := NewObject()
	xlsxInitClasses()
	lib.SetMember("File", xlsxFileClass, c)
	lib.SetMember("open", NewNativeFunction("open", func(c *Context, _ Value, args []Value) Value {
		var fn ValueStr
		EnsureFuncParams(c, "openFile", args, ArgRuleRequired("filename", TypeStr, &fn))
		return NewObjectAndInit(xlsxFileClass, c, fn)
	}, "filename"), c)
	lib.SetMember("new", NewNativeFunction("newFile", func(c *Context, _ Value, args []Value) Value {
		return NewObjectAndInit(xlsxFileClass, c)
	}), c)
	return lib
}

type (
	xlsxFileReserved struct {
		file *excelize.File
	}
)

func xlsxInitClasses() {
	xlsxFileClass = NewClassBuilder("File").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			var file *excelize.File
			EnsureFuncParams(c, "File.__init__", args,
				NewOneOfHelper("from").
					On(TypeStr, func(a Value) {
						var err error
						filename := a.ToString(c)
						if file, err = excelize.OpenFile(filename); err != nil {
							c.RaiseRuntimeError("File.__init__: open xlsx file %s error %+v", filename, err)
						}
					}).
					On(ptablePTableClass, func(a Value) {
						file = excelize.NewFile()
						sheet := file.GetSheetName(0)
						xlsxPTableToSheet(c, file, sheet, a.(ValueObject))
					}).
					On(TypeGoValue, func(a Value) {
						gv := a.ToGoValue(c)
						var err error
						switch v := gv.(type) {
						case *excelize.File:
							file = v
						case io.Reader:
							if file, err = excelize.OpenReader(v); err != nil {
								c.RaiseRuntimeError("File.__init__: open xlsx from io.Reader error %+v", err)
							}
						default:
							c.RaiseRuntimeError("File.__init__: invalid input argument")
						}
					}).
					Default(func() {
						file = excelize.NewFile()
					}),
			)
			this.Reserved = &xlsxFileReserved{file: file}
		}).
		Method("sheetCount", func(c *Context, this ValueObject, args []Value) Value {
			f := this.Reserved.(*xlsxFileReserved).file
			return NewInt(int64(f.SheetCount))
		}).
		Method("sheetNames", func(c *Context, this ValueObject, args []Value) Value {
			f := this.Reserved.(*xlsxFileReserved).file
			return NewArrayByValues(lo.Map(f.GetSheetList(), func(n string, _ int) Value { return NewStr(n) })...)
		}).
		Method("sheetStrs", func(c *Context, this ValueObject, args []Value) Value {
			f := this.Reserved.(*xlsxFileReserved).file
			var (
				err  error
				data [][]string
			)
			EnsureFuncParams(c, "sheetData", args, NewOneOfHelper("sheet").
				On(TypeStr, func(sheet Value) {
					if data, err = f.GetRows(sheet.ToString(c)); err != nil {
						c.RaiseRuntimeError("sheetData: get rows from sheet %s error %+v", sheet.ToString(c), err)
					}
				}).
				On(TypeInt, func(sheet Value) {
					i := sheet.(ValueInt).AsInt()
					if i < 0 || i >= f.SheetCount {
						c.RaiseRuntimeError("sheetData: get rows from sheet %d error: index out of range", i)
					}
					name := f.GetSheetList()[i]
					if data, err = f.GetRows(name); err != nil {
						c.RaiseRuntimeError("sheetData: get rows from sheet %d error %+v", i, err)
					}
				}),
			)
			rv := NewArray(len(data))
			for _, row := range data {
				rv.PushBack(NewArrayByValues(lo.Map(row, func(v string, _ int) Value { return NewStr(v) })...))
			}
			return rv
		}).
		Method("sheetValues", func(c *Context, this ValueObject, args []Value) Value {
			f := this.Reserved.(*xlsxFileReserved).file
			var (
				err   error
				sheet string
				data  [][]string
			)
			EnsureFuncParams(c, "sheetValues", args, NewOneOfHelper("sheet").
				On(TypeStr, func(sv Value) {
					sheet = sv.ToString(c)
					if data, err = f.GetRows(sheet); err != nil {
						c.RaiseRuntimeError("sheetData: get rows from sheet %s error %+v", sheet, err)
					}
				}).
				On(TypeInt, func(sv Value) {
					i := sv.(ValueInt).AsInt()
					if i < 0 || i >= f.SheetCount {
						c.RaiseRuntimeError("sheetData: get rows from sheet %d error: index out of range", i)
					}
					sheet = f.GetSheetList()[i]
					if data, err = f.GetRows(sheet); err != nil {
						c.RaiseRuntimeError("sheetData: get rows from sheet %d error %+v", i, err)
					}
				}),
			)
			rv := NewArray(len(data))
			for i, row := range data {
				vrow := NewArray(len(row))
				for j, cell := range row {
					pos := xlsxPosFromXY(c, j, i)
					if cellType, err := f.GetCellType(sheet, pos); err != nil {
						c.RaiseRuntimeError("sheetValues get cellType pos %s:%s error %+v", sheet, pos, err)
					} else {
						switch cellType {
						case excelize.CellTypeBool:
							if cv, err := strconv.ParseBool(cell); err != nil {
								c.RaiseRuntimeError("sheetValues parse cell %s to bool error %+v", pos, err)
							} else {
								vrow.PushBack(NewBool(cv))
							}
						case excelize.CellTypeNumber:
							if ci, err := strconv.ParseInt(cell, 10, 64); err == nil {
								vrow.PushBack(NewInt(ci))
							} else if cf, err := strconv.ParseFloat(cell, 64); err == nil {
								vrow.PushBack(NewFloat(cf))
							} else {
								c.RaiseRuntimeError("sheetValues parse cell %s to number error %+v", pos, err)
							}
						default:
							vrow.PushBack(NewStr(cell))
						}
					}
				}
				rv.PushBack(vrow)
			}
			return rv
		}).
		Method("save", func(c *Context, this ValueObject, args []Value) Value {
			f := this.Reserved.(*xlsxFileReserved).file
			EnsureFuncParams(c, "save", args, NewOneOfHelper("dst").
				On(TypeStr, func(a Value) {
					if err := f.SaveAs(a.ToString(c)); err != nil {
						c.RaiseRuntimeError("save File to %s error %+v", a.ToString(c), err)
					}
				}))
			return Undefined()
		}).
		Build()
}

func xlsxPosFromXY(c *Context, x, y int) string {
	if p, err := excelize.CoordinatesToCellName(x+1, y+1); err != nil {
		c.RaiseRuntimeError("CoordinatesToCellName(%d, %d) error %+v", x, y)
		return ""
	} else {
		return p
	}
}

func xlsxPTableToSheet(c *Context, file *excelize.File, sheet string, ptable ValueObject) {
	table := c.InvokeMethod(ptable, "toArray", Args(NewBool(true))).(ValueArray)
	for i := 0; i < table.Len(); i++ {
		row := table.GetIndex(i, c).(ValueArray)
		for j := 0; j < row.Len(); j++ {
			pos := xlsxPosFromXY(c, j, i)
			switch v := row.GetIndex(j, c).(type) {
			case ValueInt:
				file.SetCellInt(sheet, pos, v.Value())
			case ValueFloat:
				file.SetCellFloat(sheet, pos, v.Value(), -1, 64)
			case ValueBool:
				file.SetCellBool(sheet, pos, v.Value())
			default:
				file.SetCellStr(sheet, pos, v.ToString(c))
			}
		}
	}
}
