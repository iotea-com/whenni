package certificates

type CaCert struct {
	CaCertPem string `json:"caCertPem"`
}

type ClientKeypair struct {
	ClientCertPem string `json:"clientCertPem,omitempty"`
	ClientCertKey string `json:"clientCertKey,omitempty"`
}
