package udn

import "reflect"

var getInterface reflect.Type
var setInterface reflect.Type

func init() {
	getInterface = reflect.TypeOf((*UDNGetter)(nil)).Elem()
	setInterface = reflect.TypeOf((*UDNSetter)(nil)).Elem()
}

// Interface to implement which lets you inject in the middle of the UDN
// resolution to participate in the lookup. This is useful to abstract
// away various forms of storage/access-- such as databases, disks, network,
// etc. All by implementing this interface
type UDNGetter interface {
	GetUDN(keyParts []string) (val interface{}, resolved int, err error)
}

// Setter interface, it will be passed a list of keys-- if the value is set
// the error must be nil-- otherwise we assume nothing was set
type UDNSetter interface {
	SetUDN(keyParts []string, val interface{}) error
}
