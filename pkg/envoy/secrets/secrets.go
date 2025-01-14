package secrets

import (
	"strings"

	"github.com/openservicemesh/osm/pkg/service"
)

const (
	// namespaceNameSeparator used upon marshalling/unmarshalling MeshService to a string
	// or viceversa
	namespaceNameSeparator = "/"
)

// UnmarshalSDSCert parses the SDS resource name and returns an SDSCert object and an error if any
// Examples:
// 1. Unmarshalling 'service-cert:foo/bar' returns SDSCert{CertType: service-cert, Name: foo/bar}, nil
// 2. Unmarshalling 'root-cert-for-mtls-inbound:foo/bar' returns SDSCert{CertType: root-cert-for-mtls-inbound, Name: foo/bar}, nil
// 3. Unmarshalling 'invalid-cert' returns nil, error
func UnmarshalSDSCert(str string) (*SDSCert, error) {
	var ret SDSCert

	// Check separators, ignore empty string fields
	slices := strings.Split(str, Separator)
	if len(slices) != 2 {
		return nil, errInvalidCertFormat
	}

	// Make sure the slices are not empty. Split might actually leave empty slices.
	for _, sep := range slices {
		if len(sep) == 0 {
			return nil, errInvalidCertFormat
		}
	}

	// Check valid certType
	ret.CertType = SDSCertType(slices[0])
	if _, ok := validCertTypes[ret.CertType]; !ok {
		return nil, errInvalidCertFormat
	}

	ret.Name = slices[1]

	return &ret, nil
}

// GetMeshService unmarshals a NamespaceService type from a SDSCert name
func (sdsc *SDSCert) GetMeshService() (*service.MeshService, error) {
	slices := strings.Split(sdsc.Name, namespaceNameSeparator)
	if len(slices) != 2 {
		return nil, service.ErrInvalidMeshServiceFormat
	}

	// Make sure the slices are not empty. Split might actually leave empty slices.
	for _, sep := range slices {
		if len(sep) == 0 {
			return nil, service.ErrInvalidMeshServiceFormat
		}
	}

	return &service.MeshService{
		Namespace: slices[0],
		Name:      slices[1],
	}, nil
}
