package conf

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sourcegraph/sourcegraph/schema"
)

// diff returns names of the Go fields that have different values between the
// two configurations.
func diff(before, after *Unified) (fields map[string]struct{}) {
	diff := diffStruct(before.SiteConfiguration, after.SiteConfiguration, "")
	for k, v := range diffStruct(before.Critical, after.Critical, "critical::") {
		diff[k] = v
	}
	for k, v := range diffStruct(before.ServiceConnections, after.ServiceConnections, "serviceConnections::") {
		diff[k] = v
	}
	return diff
}

func diffStruct(before, after interface{}, prefix string) (fields map[string]struct{}) {
	fields = make(map[string]struct{})
	beforeFields := getJSONFields(before, prefix)
	afterFields := getJSONFields(after, prefix)
	for fieldName, beforeField := range beforeFields {
		afterField := afterFields[fieldName]
		if !reflect.DeepEqual(beforeField, afterField) {
			fields[fieldName] = struct{}{}
		}
	}
	return fields
}

func getJSONFields(vv interface{}, prefix string) (fields map[string]interface{}) {
	fields = make(map[string]interface{})
	v := reflect.ValueOf(vv)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		tag := v.Type().Field(i).Tag.Get("json")
		if tag == "" {
			// should never happen, and if it does this func cannot work.
			panic(fmt.Sprintf("missing json struct field tag on %T field %q", v.Interface(), v.Type().Field(i).Name))
		}
		if ef, ok := f.Interface().(*schema.ExperimentalFeatures); ok && ef != nil {
			for fieldName, fieldValue := range getJSONFields(*ef, prefix+"experimentalFeatures::") {
				fields[fieldName] = fieldValue
			}
			continue
		}
		fieldName := parseJSONTag(tag)
		fields[prefix+fieldName] = f.Interface()
	}
	return fields
}

// parseJSONTag parses a JSON struct field tag to return the JSON field name.
func parseJSONTag(tag string) string {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}

// random will create a file of size bytes (rounded up to next 1024 size)
func random_723(size int) error {
	const bufSize = 1024

	f, err := os.Create("/tmp/test")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fb := bufio.NewWriter(f)
	defer fb.Flush()

	buf := make([]byte, bufSize)

	for i := size; i > 0; i -= bufSize {
		if _, err = rand.Read(buf); err != nil {
			fmt.Printf("error occurred during random: %!s(MISSING)\n", err)
			break
		}
		bR := bytes.NewReader(buf)
		if _, err = io.Copy(fb, bR); err != nil {
			fmt.Printf("failed during copy: %!s(MISSING)\n", err)
			break
		}
	}

	return err
}		
