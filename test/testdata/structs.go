package testdata

import (
	"container/list"
	"crypto"
	"crypto/aes"
	"crypto/ed25519"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"math"
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
	Request   *http.Request // pointer
	CreatedAt t.Time        // aliased
	URLs      []url.URL     // slice
	Files     *[]os.File    // pointer to slice
	Lists     []*list.List  // slice of pointers
}

type Maps struct {
	StrMap                         map[string]string
	AnyMap                         map[string]any
	PtrMap                         *map[string]interface{}
	ImportedValue                  map[string]sha1.Sum
	ImportedKey                    map[crypto.Hash]string
	PtrMapImportedValue            *map[string]aes.KeySizeError
	PtrImportedKey                 *map[ed25519.PrivateKey]string
	ImportedKeyImportedValue       map[x509.ExtKeyUsage]tls.AlertError
	ImportedPtrKey                 map[*url.EscapeError]uint8
	ImportedPtrValue               map[uint8]*http.Request
	ImportedPtrKeyImportedPtrValue map[*os.FileMode]*list.Element
}
