package builtin_libs

import (
	. "github.com/zgg-lang/zgg-go/runtime"
)

var chartChartType ValueType

func libChart(c *Context) ValueObject {
	lib := NewObject()
	return lib
}

func init() {
	chartChartType = NewClassBuilder("Chart").
		Constructor(func(c *Context, this ValueObject, args []Value) {
			this.SetMember("__allSeries", NewGoValue([]ValueArray{}), c)
		}).
		Method("addSeries", func(c *Context, this ValueObject, args []Value) Value {
			var series ValueArray
			EnsureFuncParams(c, "Chart.addSeries", args,
				ArgRuleRequired("series", TypeArray, &series),
			)
			allSeries := this.GetMember("__allSeries", c).ToGoValue().([]ValueArray)
			allSeries = append(allSeries, series)
			this.SetMember("__allSeries", NewGoValue(allSeries), c)
			return this
		}).
		// .Method("bar", func (c *Context, this ValueObject, args []Value) Value {
		// 	s := []chart.Series{}
		// 	allSeries := this.GetMember("__allSeries").ToGoValue().([]ValueArray)

		// 	ch := chart.Chart{
		// 		Series: []chart.Series
		// 	}
		// }).
		Build()
}
