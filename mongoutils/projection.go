package mongoutils

import (
	"reflect" //nolint:depguard // reflect is required in order to compute the projection.
	"strings"
	"sync"

	"github.com/siddhant2408/golang-libraries/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// Projection returns a projection for the given value.
//
// It is forbidden to update the returned value, because it is shared in a cache.
func Projection(v interface{}) interface{} {
	t := reflect.TypeOf(v)
	if res, ok := projectionCache.Load(t); ok {
		return res
	}
	res := projectionType(t)
	projectionCache.Store(t, res)
	return res
}

var projectionCache sync.Map

// ProjectionNoCache is identical to Projection except:
//  - it uses no cache
//  - it returns the result as a bson.M
//  - it is allowed to update the returned value (since it uses no cache)
func ProjectionNoCache(v interface{}) bson.M {
	return projectionType(reflect.TypeOf(v))
}

func projectionType(t reflect.Type) bson.M {
	var res bson.M
	switch t.Kind() {
	case reflect.Array, reflect.Slice, reflect.Ptr:
		res = projectionType(t.Elem())
	case reflect.Struct:
		res = projectionStruct(t)
	default:
		panic(errors.Newf("unsupported type %v", t))
	}
	return res
}

func projectionStruct(t reflect.Type) bson.M {
	sl := make(bson.M)
	for i, c := 0, t.NumField(); i < c; i++ {
		name := projectionStructField(t.Field(i))
		if name != "" {
			sl[name] = 1
		}
	}
	return sl
}

func projectionStructField(sf reflect.StructField) string {
	if sf.PkgPath != "" {
		// It means that the field is unexported.
		return ""
	}
	name := projectionStructFieldTag(sf)
	if name == "" {
		// bson uses the lower case struct field name if not specified.
		name = strings.ToLower(sf.Name)
	}
	return name
}

func projectionStructFieldTag(sf reflect.StructField) string {
	tag, ok := sf.Tag.Lookup("bson")
	if !ok {
		return ""
	}
	return strings.TrimSpace(strings.Split(tag, ",")[0])
}
