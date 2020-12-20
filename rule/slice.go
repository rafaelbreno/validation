package rule

import (
	"github.com/lanvard/support"
	"reflect"
)

type Slice struct{}

func (s Slice) Verify(value support.Value) error {
	kind := support.Kind(value.Raw())
	if kind != reflect.Array && kind != reflect.Slice {
		return MustBeASliceError
	}
	return nil
}
