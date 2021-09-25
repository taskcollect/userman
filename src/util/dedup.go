package util

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/buger/jsonparser"
)

// removes all default values from a json string
// real is the json string to be cleaned
// template is the json string with default values
// noNewKeys dictates whether any keys in the real not in the default should error out
// strictType dictates whether any keys in the real with a different type than the default should error out
func RemoveDefaultKeys(real, template []byte, noNewKeys bool, strictType bool) ([]byte, error) {
	err := jsonparser.ObjectEach(real,
		func(key []byte, rValue []byte, rType jsonparser.ValueType, offset int) error {
			dValue, dType, _, err := jsonparser.Get(template, string(key))

			if err != nil && dType != jsonparser.NotExist {
				// only filter parsing errors
				return err
			}

			if dType == jsonparser.NotExist {
				// default value does not exist for its counterpart in the real
				if noNewKeys {
					return fmt.Errorf("key %s of real was not in default", string(key))
				}
				return nil
			}

			if strictType && dType != rType {
				// default value is of different type
				return fmt.Errorf("key %s of type %s does not match templated type %s", string(key), rType, dType)
			}

			if bytes.Equal(rValue, dValue) {
				// the value is the same, remove it from the real
				real = jsonparser.Delete(real, string(key))
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return real, nil
}

func AddDefaultKeys(overrides, template []byte, noNewKeys bool) ([]byte, error) {
	// iterate through all of the keys of the template
	err := jsonparser.ObjectEach(template,
		func(key []byte, tValue []byte, tType jsonparser.ValueType, offset int) error {
			// get the override's counterpart to the template key
			_, oType, _, err := jsonparser.Get(overrides, string(key))

			if err != nil && oType != jsonparser.NotExist {
				return err
			}

			if oType == jsonparser.NotExist {
				// key exists in template but not in overrides

				if tType == jsonparser.String {
					// add quotes to the value if it's a string
					tValue = []byte(strconv.Quote(string(tValue)))
				}

				overrides, err = jsonparser.Set(overrides, tValue, string(key))
				if err != nil {
					return err
				}
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return overrides, nil
}
