package utils


import (
	"reflect"
)

// Array 数组
type Array []interface{}

// ToStringSlice convert to string clice
func (array Array) ToStringSlice() []string {
	l := make([]string, 0)
	for _, v := range array {
		l = append(l, v.(string))
	}

	return l
}

// ToIntSlice Convert to int slice
func (array Array) ToIntSlice() []int {
	l := make([]int, 0)
	for _, v := range array {
		l = append(l, v.(int))
	}

	return l
}

// ToFloat32Slice Convert to int slice
func (array Array) ToFloat32Slice() []float32 {
	l := make([]float32, 0)
	for _, v := range array {
		l = append(l, v.(float32))
	}

	return l
}

// ToFloat64Slice Convert to int slice
func (array Array) ToFloat64Slice() []float64 {
	l := make([]float64, 0)
	for _, v := range array {
		l = append(l, v.(float64))
	}

	return l
}

// IsExistItem 判断slice中是否存在某个item
func IsExistItem(value interface{}, array interface{}) bool {
    switch reflect.TypeOf(array).Kind() {
    case reflect.Slice:
        s := reflect.ValueOf(array)
        for i := 0; i < s.Len(); i++ {
            if reflect.DeepEqual(value, s.Index(i).Interface()) {
                return true
            }
        }
    }
    return false
}

// RemoveDuplicates 去掉slice中的重复项
func RemoveDuplicates(array interface{}) Array {
    temp := make(map[interface{}]bool, 0)
    result := make([]interface{}, 0)
	av := reflect.ValueOf(array)
    for i := 0; i < av.Len(); i++ {
        if _, ok := temp[av.Index(i).Interface()]; !ok {
            temp[av.Index(i).Interface()] = true
            result = append(result, av.Index(i).Interface())
        }
    }

    return result
}

// Union 并集, 两个slice A，B，把他们所有的元素合并在一起组成的slice
func Union(a, b interface{}) Array {
	result := make([]interface{}, 0)
	temp := make(map[interface{}]bool, 0)

	
	av := reflect.ValueOf(a)
	for i := 0; i < av.Len(); i++ {
		if _, ok := temp[av.Index(i).Interface()]; ok {
			continue
		}
		temp[av.Index(i).Interface()] = true
		result = append(result, av.Index(i).Interface())
	}

	bv := reflect.ValueOf(b)
	for i := 0; i < bv.Len(); i++ {
		if _, ok := temp[bv.Index(i).Interface()]; ok {
			continue
		}
		temp[bv.Index(i).Interface()] = true
		result = append(result, bv.Index(i).Interface())
	}
	return result
}

// Intersection 交集
func Intersection(a, b interface{}) Array {
    var inter []interface{}
    temp := make(map[interface{}]bool)

	av := reflect.ValueOf(a)
    for i := 0; i < av.Len(); i++ {
        if _, ok := temp[av.Index(i).Interface()]; !ok {
            temp[av.Index(i).Interface()] = true
        }
    }
	bv := reflect.ValueOf(b)
    for i := 0; i < bv.Len(); i++ {
        if _, ok := temp[bv.Index(i).Interface()]; ok {
            inter = append(inter, bv.Index(i).Interface())
        }
    }

    return inter
}

// Difference 差集
func Difference(a, b interface{}) Array {
	result := make([]interface{}, 0)
	temp := map[interface{}]struct{}{}

	bs := reflect.ValueOf(b)
	for i := 0; i < bs.Len(); i++ {
		if _, ok := temp[bs.Index(i).Interface()]; !ok {
			temp[bs.Index(i).Interface()] = struct{}{}
		}
	}

	as := reflect.ValueOf(a)
	for i := 0; i < as.Len(); i++ {
		if _, ok := temp[as.Index(i).Interface()]; !ok {
			result = append(result, as.Index(i).Interface())
		}
	}

	return result
}