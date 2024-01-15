package utils

import (
	"fmt"
	"math"
	"reflect"
)

func AnyToFloat64(in interface{}) (float64, error) {
	var floatType = reflect.TypeOf(float64(0))
	v := reflect.ValueOf(in)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}

func Float64ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Float64Round(num*output)) / output
}

func Float64Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
