package secrets

import (
	"fmt"
	"testing"

	tassert "github.com/stretchr/testify/assert"
	trequire "github.com/stretchr/testify/require"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/openservicemesh/osm/pkg/service"
)

var _ = Describe("Test secert tools", func() {
	Context("Test UnmarshalSDSCert()", func() {
		It("Interface marshals and unmarshals preserving the exact same data", func() {
			InitialObj := SDSCert{
				CertType: ServiceCertType,
				Name:     "test-namespace/test-service",
			}

			// Marshal/stringify it
			marshaledStr := InitialObj.String()

			// Unmarshal it back from the string
			finalObj, _ := UnmarshalSDSCert(marshaledStr)

			// First and final object must be equal
			Expect(*finalObj).To(Equal(InitialObj))
		})

		It("returns service cert", func() {
			actual, err := UnmarshalSDSCert("service-cert:namespace-test/blahBlahBlahCert")
			Expect(err).ToNot(HaveOccurred())
			Expect(actual.CertType).To(Equal(ServiceCertType))
			Expect(actual.Name).To(Equal("namespace-test/blahBlahBlahCert"))
		})
		It("returns root cert for mTLS", func() {
			actual, err := UnmarshalSDSCert("root-cert-for-mtls-outbound:namespace-test/blahBlahBlahCert")
			Expect(err).ToNot(HaveOccurred())
			Expect(actual.CertType).To(Equal(RootCertTypeForMTLSOutbound))
			Expect(actual.Name).To(Equal("namespace-test/blahBlahBlahCert"))

		})

		It("returns root cert for non-mTLS", func() {
			actual, err := UnmarshalSDSCert("root-cert-https:namespace-test/blahBlahBlahCert")
			Expect(err).ToNot(HaveOccurred())
			Expect(actual.CertType).To(Equal(RootCertTypeForHTTPS))
			Expect(actual.Name).To(Equal("namespace-test/blahBlahBlahCert"))
		})

		It("returns an error (invalid formatting)", func() {
			_, err := UnmarshalSDSCert("blahBlahBlahCert")
			Expect(err).To(HaveOccurred())
		})

		It("returns an error (invalid formatting)", func() {
			_, err := UnmarshalSDSCert("blahBlahBlahCert:moreblabla/amazingservice:bla")
			Expect(err).To(HaveOccurred())
		})

		It("returns an error (missing cert type)", func() {
			_, err := UnmarshalSDSCert("blahBlahBlahCert/service")
			Expect(err).To(HaveOccurred())
		})

		It("returns an error (invalid serv type)", func() {
			_, err := UnmarshalSDSCert("revoked-cert:blah/BlahBlahCert")
			Expect(err).To(HaveOccurred())
		})

		It("returns an error (invalid mtls cert type)", func() {
			_, err := UnmarshalSDSCert("oot-cert-for-mtls-diagonalstream:blah/BlahBlahCert")
			Expect(err).To(HaveOccurred())
		})

		It("returns an error (empty slice)", func() {
			_, err := UnmarshalSDSCert(":")
			Expect(err).To(HaveOccurred())
		})
	})
})

func TestUnmarshalMeshService(t *testing.T) {
	assert := tassert.New(t)
	require := trequire.New(t)

	namespace := "randomNamespace"
	serviceName := "randomServiceName"
	meshService := &service.MeshService{
		Namespace: namespace,
		Name:      serviceName,
	}
	str := meshService.String()
	fmt.Println(str)

	testCases := []struct {
		name        string
		expectedErr bool
		sdsCert     SDSCert
	}{
		{
			name:        "successfully unmarshal service",
			expectedErr: false,
			sdsCert: SDSCert{
				Name: "randomNamespace/randomServiceName",
			},
		},
		{
			name:        "incomplete namespaced service name 1",
			expectedErr: true,
			sdsCert: SDSCert{
				Name: "/svnc",
			},
		},
		{
			name:        "incomplete namespaced service name 2",
			expectedErr: true,
			sdsCert: SDSCert{
				Name: "svnc/",
			},
		},
		{
			name:        "incomplete namespaced service name 3",
			expectedErr: true,
			sdsCert: SDSCert{
				Name: "/svnc/",
			},
		},
		{
			name:        "incomplete namespaced service name 3",
			expectedErr: true,
			sdsCert: SDSCert{
				Name: "/",
			},
		},
		{
			name:        "incomplete namespaced service name 3",
			expectedErr: true,
			sdsCert: SDSCert{
				Name: "",
			},
		},
		{
			name:        "incomplete namespaced service name 3",
			expectedErr: true,
			sdsCert: SDSCert{
				Name: "test",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.sdsCert.GetMeshService()
			if tc.expectedErr {
				assert.NotNil(err)
			} else {
				require.Nil(err)
				assert.Equal(meshService, actual)
			}
		})
	}
}
