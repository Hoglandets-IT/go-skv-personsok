package personsok

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"net/http"
	"os"

	"github.com/hoglandets-it/goskv/models"
	gofncfg "github.com/scheiblingco/gofn/cfgtools"
)

type SkvClient struct {
	Url string `json:"url"`

	CertPath string `json:"certPath"`
	CertPem  []byte `json:"certPem"`
	KeyPath  string `json:"keyPath"`
	KeyPem   []byte `json:"keyPem"`

	OrgId   string `json:"orgId"`
	OrderId string `json:"orderId"`

	client *http.Client
}

func (c *SkvClient) Init() error {
	certPem, err := os.ReadFile(c.CertPath)
	if err != nil {
		return err
	}

	keyPem, err := os.ReadFile(c.KeyPath)
	if err != nil {
		return err
	}

	c.CertPem = certPem
	c.KeyPem = keyPem

	cert, err := tls.X509KeyPair(c.CertPem, c.KeyPem)
	if err != nil {
		return err
	}

	c.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
				ClientAuth:   tls.RequireAnyClientCert,
				VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
					return nil
				},
				InsecureSkipVerify: true,
			},
		},
	}

	return nil
}

func (c *SkvClient) SearchPerson(personId string, cache bool) (*models.PersonpostResponseEnvelope, error) {
	req, err := models.NewPersonpostV4Request(c.OrgId, c.OrderId, personId)
	if err != nil {
		panic(err)
	}
	xmlBytes, err := xml.MarshalIndent(req, "", "    ")
	if err != nil {
		return nil, err
	}

	hr, err := http.NewRequest("POST", c.Url, bytes.NewBuffer(xmlBytes))
	if err != nil {
		return nil, err
	}

	hr.Header.Set("Content-Type", "text/xml; charset=utf-8")
	// hr.Header.Set("Connection", "close")

	resp, err := c.client.Do(hr)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	obj := models.PersonpostResponseEnvelope{}
	deco := xml.NewDecoder(resp.Body)

	if cache {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		os.WriteFile("response.xml", buf.Bytes(), 0644)

		deco = xml.NewDecoder(buf)
	}

	err = deco.Decode(&obj)
	return &obj, err
}

func NewClientFromConfigJson(configPath string) (*SkvClient, error) {
	config := SkvClient{}

	if err := gofncfg.LoadJsonConfig(configPath, &config); err != nil {
		return nil, err
	}

	err := config.Init()
	if err != nil {
		return nil, err
	}

	return &config, nil
}
