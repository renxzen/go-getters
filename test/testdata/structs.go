package testdata

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	t "time"
)

type Example struct {
	Name   string
	Value  int
	Active bool
}

type Container struct {
	Example    Example
	ExamplePtr *Example
}

type Pointer struct {
	Name  *string
	Age   *int
	Score *float64
	Flag  *bool
}

type Slices struct {
	SlicePtr  *[]Example
	Slice     []Example
	SliceInt  []int
	SliceStr  []string
	SliceBool []bool
}

type DynamicImports struct {
	// Context context.Context // TODO: interface detection
	Request    *http.Request   // pointer
	CreatedAt  t.Time          // aliased
	URLs       []url.URL       // slice
	Files      *[]os.File      // pointer to slice
	Formatters []*fmt.Stringer // slice of pointers
}
