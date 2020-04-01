/*+-----------------------------+
 *| Author: Zoueature           |
 *+-----------------------------+
 *| Email: zoueature@gmail.com  |
 *+-----------------------------+
 */
package bencode

import (
	"errors"
	"reflect"
	"strconv"
)

const (
	bDecodeInt      = 'i'
	bDecodeList     = 'l'
	bDecodeHash     = 'd'
	bDecodeString   = 's'
	bDecodeEnd      = 'e'
	bDecodeFormated = 'f'

	bEncodeTag = "bencode"
)

func bDecode(torrentContent []byte) (interface{}, []byte, error) {
	//bdecode
	var typeStack []byte
	var matchContainer []interface{}
	var strMatcher []byte  //存储匹配的字符串长度
	var numMatcher []byte  //存储匹配的数值
	var startNumMatch bool //标识是否开启数值匹配
	var firstType byte
	var matchInfoLevel, infoStartIndex, infoEndIndex int
	for i := 0; i < len(torrentContent); i++ {
		b := torrentContent[i]
		if b == bDecodeHash || b == bDecodeList {
			if i == 0 {
				firstType = b
			}
			typeStack = append(typeStack, b)
			if len(matchContainer) > 0 && matchInfoLevel == 0 {
				v, ok := matchContainer[len(matchContainer)-1].(string)
				if ok && v == "info" {
					matchInfoLevel++
					infoStartIndex = i
				}
			} else if matchInfoLevel > 0 {
				matchInfoLevel++
			}
		} else if b == bDecodeInt {
			startNumMatch = true
			if matchInfoLevel > 0 {
				matchInfoLevel++
			}
		} else if b >= '0' && b <= '9' {
			if startNumMatch {
				numMatcher = append(numMatcher, b)
			} else {
				strMatcher = append(strMatcher, b)
			}
		} else if b == ':' {
			//字符串长度值匹配结束
			strLenStr := string(strMatcher)
			strLen, err := strconv.Atoi(strLenStr)
			if err != nil {
				return nil, nil, err
			}
			str := torrentContent[i+1 : i+1+strLen]
			i += strLen
			matchContainer = append(matchContainer, string(str))
			strMatcher = append(strMatcher[0:0])
			typeStack = append(typeStack, bDecodeString)
		} else if b == bDecodeEnd {
			if matchInfoLevel == 1 {
				infoEndIndex = i + 1
				matchInfoLevel--
			} else if matchInfoLevel > 0 {
				matchInfoLevel--
			}
			if startNumMatch {
				//数值匹配
				matchContainer = append(matchContainer, string(numMatcher))
				startNumMatch = false
				numMatcher = append(numMatcher[0:0])
				typeStack = append(typeStack, bDecodeInt)
				continue
			}
			tmp := make([]interface{}, 0)
			var nowType byte
			typeLen := len(typeStack)
			var j int
			for j = 0; j < typeLen; j++ {
				nowType = typeStack[len(typeStack)-j-1]
				if nowType == bDecodeFormated || nowType == bDecodeInt || nowType == bDecodeString {
					tmp = append(tmp, matchContainer[len(matchContainer)-j-1])
				} else {
					break
				}
			}
			//if len(tmp) == 0 {
			//	return nil, nil, errors.New("format data error")
			//}

			matchContainer = append(matchContainer[:len(matchContainer)-j])
			typeStack = append(typeStack[:len(typeStack)-j-1])
			var data interface{}
			if nowType == bDecodeList {
				data = tmp
			} else if nowType == bDecodeHash {
				l := len(tmp)
				if l%2 != 0 {
					return nil, nil, errors.New("format map error, item num error ")
				}
				m := make(map[string]interface{})
				var key string
				for k := l; k > 0; k-- {
					index := k - 1
					if k%2 == 0 {
						var ok bool
						key, ok = tmp[index].(string)
						if !ok {
							return nil, nil, errors.New("format map error, trans to key string error ")
						}
					} else {
						m[key] = tmp[index]
					}
				}
				data = m
			}
			matchContainer = append(matchContainer, data)
			typeStack = append(typeStack, bDecodeFormated)
		}
	}
	if len(matchContainer) < 0 {
		return nil, nil, errors.New("error,  bdecode empty")
	}
	info := torrentContent[infoStartIndex:infoEndIndex]
	if firstType == bDecodeHash {
		return matchContainer[0], info, nil
	}
	return matchContainer, info, nil
}

func BDecode(encodeStr []byte, container interface{}) error {
	ct := reflect.ValueOf(container)
	if ct.Kind() != reflect.Ptr {
		return errors.New("not a ptr")
	}
	elm := ct.Elem()
	if elm.Kind() != reflect.Struct {
		return errors.New("only support ptr of struct, ptr of " + elm.Kind().String() + " gave")
	}
	bInfo, _, err := bDecode(encodeStr)
	if err != nil {
		return err
	}
	bInfoMap, ok := bInfo.(map[string]interface{})
	if !ok {
		return errors.New("decode to map error")
	}
	err = fillStruct(bInfoMap, reflect.TypeOf(container).Elem(), &elm)
	return err
}

func fillStruct(data map[string]interface{}, sct reflect.Type, scv *reflect.Value) error {
	fieldNum := scv.NumField()
	for i := 0; i < fieldNum; i++ {
		fieldType := sct.Field(i)
		//if fieldType.Anonymous {
		//	continue
		//}
		fieldTag := fieldType.Tag.Get(bEncodeTag)
		if fieldTag == "" {
			continue
		}
		fieldData, ok := data[fieldTag]
		if !ok {
			continue
		}
		fieldValue := scv.Field(i)
		switch fieldType.Type.Kind() {
		case reflect.Struct:
			fieldValueMap, succ := fieldData.(map[string]interface{})
			if !succ {
				return errors.New(fieldType.Name + " struct type match error")
			}
			err := fillStruct(fieldValueMap, fieldType.Type, &fieldValue)
			if err != nil {
				return err
			}
		case reflect.String:
			fieldValueString, succ := fieldData.(string)
			if !succ {
				return errors.New(fieldType.Name + " string type match error")
			}
			fieldValue.SetString(fieldValueString)
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			fieldValueString, succ := fieldData.(string)
			if !succ {
				return errors.New(fieldType.Name + " integer type match error")
			}
			fieldValuFloat, err := strconv.ParseFloat(fieldValueString, 64)
			if err != nil {
				return err
			}
			fieldValueInt := int64(fieldValuFloat)
			fieldValue.SetInt(fieldValueInt)
		case reflect.Slice:
			fieldValueSlice, succ := fieldData.([]interface{})
			if !succ {
				return errors.New(fieldType.Name + " slice type match error")
			}
			reflectValue := reflect.MakeSlice(fieldType.Type, 0, len(fieldValueSlice))
			err := fillSlice(fieldValueSlice, fieldType.Type, &reflectValue)
			if err != nil {
				return err
			}
			fieldValue.Set(reflectValue)
		case reflect.Interface:
			fieldValue.Set(reflect.ValueOf(fieldData))
		default:
			return errors.New(fieldType.Name + " error, unknow type")
		}
	}
	return nil
}

func fillSlice(data []interface{}, t reflect.Type, value *reflect.Value) error {
	switch t.Elem().Kind() {
	case reflect.Struct:
		for _, v := range data {
			if m, ok := v.(map[string]interface{}); ok {
				err := fillStruct(m, t, value)
				if err != nil {
					return err
				}
			} else {
				return errors.New("match to struct error")
			}
		}
	case reflect.Slice:
		//if value.Len() == 0 {
		//	*value = reflect.MakeSlice(t, len(data), len(data))
		//}
		for _, v := range data {
			if m, ok := v.([]interface{}); ok {
				//childValue := value.Index(index)
				tmp := reflect.MakeSlice(t.Elem(), 0, 0)
				err := fillSlice(m, t.Elem(), &tmp)
				if err != nil {
					return err
				}
				*value = reflect.Append(*value, tmp)
			} else {
				return errors.New("match to slice error")
			}
		}
	case reflect.String:
		for _, v := range data {
			if m, ok := v.(string); ok {
				*value = reflect.Append(*value, reflect.ValueOf(m))
			} else {
				return errors.New("match to string error")
			}
		}
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		for _, v := range data {
			if m, ok := v.(string); ok {
				i, err := strconv.ParseInt(m, 10, 64)
				if err != nil {
					return err
				}
				value.SetInt(i)
			} else {
				return errors.New("match to string error")
			}
		}
	case reflect.Interface:
		for _, v := range data {
			*value = reflect.Append(*value, reflect.ValueOf(v))
		}
	default:
		return errors.New("unknow type")
	}
	return nil
}

func BEncode(data interface{}) (string, error) {
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}
	if dataValue.Kind() != reflect.Struct {
		return "", errors.New("only support struct")
	}
	str, err := bEncodeStruct(reflect.TypeOf(data), dataValue)
	return str, err
}

func bEncodeStruct(t reflect.Type, v reflect.Value) (string, error) {
	encodeStr := "d"
	var str string
	var err error
	fieldNum := v.NumField()
	for i := 0; i < fieldNum; i++ {
		field := v.Field(i)
		switch field.Type().Kind() {
		case reflect.Slice:
			str, err = bEncodeSlice(field.Type(), field)
			if err != nil {
				return "", err
			}
		case reflect.Struct:
			str, err = bEncodeStruct(field.Type(), field)
			if err != nil {
				return "", err
			}
		case reflect.String, reflect.Interface:
			str, err = bEncodeString(field.Type(), field)
			if err != nil {
				return "", err
			}
		case reflect.Int:
			str, err = bEncodeInt(field.Type(), field)
			if err != nil {
				return "", err
			}
		}
		encodeStr += str
	}
	encodeStr += "e"
	return encodeStr, nil
}

func bEncodeSlice(t reflect.Type, v reflect.Value) (string, error) {
	encodeStr := "l"
	var str string
	var err error
	switch t.Elem().Kind() {
	case reflect.String, reflect.Interface:
		str, err = bEncodeString()
	}
	return "", nil
}

func bEncodeString(t reflect.Type, v reflect.Value) (string, error) {
	s, ok := v.Interface().(string)
	if !ok {
		return "", errors.New("not string")
	}
	length := len(s)
	return strconv.Itoa(length) + ":" + s, nil
}

func bEncodeInt(t reflect.Type, v reflect.Value) (string, error) {
	i, ok := v.Interface().(float64)
	if !ok {
		return "", errors.New("not a number")
	}
	s := strconv.FormatFloat(i, 'e', 0, 10)
	return "i" + s + "e", nil
}
