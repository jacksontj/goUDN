package udn

// Interface to implement which lets you inject in the middle of the UDN
// resolution to participate in the lookup. This is useful to abstract
// away various forms of storage/access-- such as databases, disks, network,
// etc. All by implementing this interface
type UDNGetter interface {
	GetUDN(keyParts []string) (val interface{}, resolved int, err error)
}
