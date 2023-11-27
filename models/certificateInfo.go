package models

type CertificateInfo struct {
	CertFilePath string // X509 cert file in .pem format
	KeyFilePath  string // X509 key file in .pem format
	Truststore   string // CA authorities in .pem format
}
