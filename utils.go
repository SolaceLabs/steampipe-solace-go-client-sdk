package solacesdk

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/url"
	"reflect"
)

func GetMaxAndLimit(maxItems int) (int, int) {
	if maxItems < 0 {
		maxItems = math.MaxInt
	}

	return maxItems, int(math.Min(float64(defaultPageSize), float64(maxItems)))
}

func ColumnsToParams(params *url.Values, columns []string) {
	for _, c := range columns {
		params.Add("cols", c)
	}
}

func SetQueryParam(obj any, field string, value any) {
	ref := reflect.ValueOf(obj)

	// if its a pointer, resolve its value
	if ref.Kind() == reflect.Ptr {
		ref = reflect.Indirect(ref)
	}

	if ref.Kind() == reflect.Interface {
		ref = ref.Elem()
	}

	// should double check we now have a struct (could still be anything)
	if ref.Kind() != reflect.Struct {
		log.Fatal("unexpected type")
	}

	prop := ref.FieldByName(field)
	prop.Set(reflect.ValueOf(value))
}

func PrettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}
