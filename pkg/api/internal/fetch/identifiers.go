package fetch

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"github.com/johannessarpola/go-lutakko-gigs/pkg/api/models"
	"reflect"
)

func createEventID(event models.Event) (string, error) {
	wid := withoutID(event)
	return hashStruct(wid)
}

// hashStruct calculates the SHA256 hash of the given struct and returns it as a hex string,
// skipping the key field.
func hashStruct(data interface{}) (string, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return "", err
	}

	hash := sha256.Sum256(buf.Bytes())
	return hex.EncodeToString(hash[:]), nil
}

func withoutID(data interface{}) interface{} {
	return filterField(data, "Id")
}

// filterField creates a copy of the struct with the key field zeroed out.
func filterField(data interface{}, skipField string) interface{} {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	newStruct := reflect.New(t).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.Name != skipField {
			newStruct.Field(i).Set(v.Field(i))
		}
	}

	return newStruct.Interface()
}
