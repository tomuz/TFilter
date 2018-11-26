package TFilter

import (
	"reflect"
	"time"
)

const (
	EQ = "EQ"
	LT = "LT"
	GT = "GT"
)

func (tf *TFilter) EQ(key string, value interface{}) *TFilter {
	return tf.searcher(EQ, key, value)
}
func (tf *TFilter) LT(key string, value interface{}) *TFilter {
	return tf.searcher(LT, key, value)
}
func (tf *TFilter) GT(key string, value interface{}) *TFilter {
	return tf.searcher(GT, key, value)
}

func (tf *TFilter) searcher(operation string, key string, value interface{}) *TFilter {
	auxObj := []interface{}{}
	maxChan := tf.size / maxChannels(tf.size)

	resChan := make(chan []interface{}, maxChan)
	for i := 0; i < tf.size; i += maxChan {
		switch operation {
		case EQ:
			go searchEQ(tf.objs[i:i+maxChan], key, value, resChan)
		case LT:
			go searchLT(tf.objs[i:i+maxChan], key, value, resChan)
		case GT:
			go searchGT(tf.objs[i:i+maxChan], key, value, resChan)
		}
	}

	for i := 0; i < tf.size; i += maxChan {
		select {
		case objs := <-resChan:
			auxObj = append(auxObj, objs...)
		}
	}
	tf.objs = tf.objs[0:len(auxObj)]
	copy(auxObj, tf.objs)

	return tf
}

func searchEQ(objs []interface{}, key string, value interface{}, resChan chan []interface{}) {
	auxObj := []interface{}{}
	for _, obj := range objs {
		RValue := reflect.ValueOf(obj)
		RType := reflect.TypeOf(obj)

		for i := 0; i < RValue.NumField(); i++ {
			fieldRType := RType.Field(i)
			fieldRValue := RValue.Field(i)

			zero := reflect.Zero(fieldRValue.Type()).Interface()
			isZero := reflect.DeepEqual(fieldRValue.Interface(), zero)
			if !isZero && fieldRType.Tag.Get("key") == key {
				if fieldRValue.Interface() == value {
					auxObj = append(auxObj, obj)
				}
			}
		}
	}
	resChan <- auxObj
}

func searchLT(objs []interface{}, key string, value interface{}, resChan chan []interface{}) {
	auxObj := []interface{}{}
	for _, obj := range objs {
		RValue := reflect.ValueOf(obj)
		RType := reflect.TypeOf(obj)

		for i := 0; i < RValue.NumField(); i++ {
			fieldRType := RType.Field(i)
			fieldRValue := RValue.Field(i)

			zero := reflect.Zero(fieldRValue.Type()).Interface()
			isZero := reflect.DeepEqual(fieldRValue.Interface(), zero)
			if !isZero && fieldRType.Tag.Get("key") == key {
				shouldAppend := false
				switch fieldRValue.Interface().(type) {
				case string:
					objVal := fieldRValue.Interface().(string)
					val := value.(string)
					if len(objVal) > len(val) {
						shouldAppend = true
					}
				case float32, float64:
					objVal := fieldRValue.Interface().(float64)
					val := value.(float64)
					if objVal > val {
						shouldAppend = true
					}
				case int, int16, int32, int64:
					objVal := fieldRValue.Interface().(int64)
					val := value.(int64)
					if objVal > val {
						shouldAppend = true
					}
				case time.Time:
					objVal := fieldRValue.Interface().(time.Time)
					val := value.(time.Time)
					if val.Before(objVal) {
						shouldAppend = true
					}

				}
				if shouldAppend {
					auxObj = append(auxObj, obj)

				}
			}
		}
	}
	resChan <- auxObj
}

func searchGT(objs []interface{}, key string, value interface{}, resChan chan []interface{}) {
	auxObj := []interface{}{}
	for _, obj := range objs {
		RValue := reflect.ValueOf(obj)
		RType := reflect.TypeOf(obj)

		for i := 0; i < RValue.NumField(); i++ {
			fieldRType := RType.Field(i)
			fieldRValue := RValue.Field(i)

			zero := reflect.Zero(fieldRValue.Type()).Interface()
			isZero := reflect.DeepEqual(fieldRValue.Interface(), zero)
			if !isZero && fieldRType.Tag.Get("key") == key {
				shouldAppend := false
				switch fieldRValue.Interface().(type) {
				case string:
					objVal := fieldRValue.Interface().(string)
					val := value.(string)
					if len(objVal) < len(val) {
						shouldAppend = true
					}
				case float32, float64:
					objVal := fieldRValue.Interface().(float64)
					val := value.(float64)
					if objVal < val {
						shouldAppend = true
					}
				case int, int16, int32, int64:
					objVal := fieldRValue.Interface().(int64)
					val := value.(int64)
					if objVal < val {
						shouldAppend = true
					}
				case time.Time:
					objVal := fieldRValue.Interface().(time.Time)
					val := value.(time.Time)
					if val.After(objVal) {
						shouldAppend = true
					}

				}
				if shouldAppend {
					auxObj = append(auxObj, obj)

				}
			}
		}
	}
	resChan <- auxObj
}
