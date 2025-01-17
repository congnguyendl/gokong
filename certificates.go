package gokong

import (
	"encoding/json"
	"fmt"
)

type CertificateClient interface {
	GetById(id string) (*Certificate, error)
	Create(certificateRequest *CertificateRequest) (*Certificate, error)
	DeleteById(id string) error
	List() (*Certificates, error)
	UpdateById(id string, certificateRequest *CertificateRequest) (*Certificate, error)
}

type certificateClient struct {
	config *Config
}

type CertificateRequest struct {
	Cert *string   `json:"cert,omitempty" yaml:"cert,omitempty"`
	Key  *string   `json:"key,omitempty" yaml:"key,omitempty"`
	Tags []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
	SNIs *[]string `json:"snis" yaml:"snis"`
}

type Certificate struct {
	Id   *string   `json:"id,omitempty" yaml:"id,omitempty"`
	Cert *string   `json:"cert,omitempty" yaml:"cert,omitempty"`
	Key  *string   `json:"key,omitempty" yaml:"key,omitempty"`
	Tags []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

type Certificates struct {
	Results []*Certificate `json:"data,omitempty" yaml:"data,omitempty"`
	Total   int            `json:"total,omitempty" yaml:"total,omitempty"`
}

const CertificatesPath = "/certificates/"

func (certificateClient *certificateClient) GetById(id string) (*Certificate, error) {
	r, body, errs := newGet(certificateClient.config, CertificatesPath+id).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get certificate, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	certificate := &Certificate{}
	err := json.Unmarshal([]byte(body), certificate)
	if err != nil {
		return nil, fmt.Errorf("could not parse certificate get response, error: %v", err)
	}

	if certificate.Id == nil {
		return nil, nil
	}

	return certificate, nil
}

func (certificateClient *certificateClient) Create(certificateRequest *CertificateRequest) (*Certificate, error) {
	r, body, errs := newPost(certificateClient.config, CertificatesPath).Send(certificateRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not create new certificate, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	createdCertificate := &Certificate{}
	err := json.Unmarshal([]byte(body), createdCertificate)
	if err != nil {
		return nil, fmt.Errorf("could not parse certificate creation response, error: %v", err)
	}

	if createdCertificate.Id == nil {
		return nil, fmt.Errorf("could not create certificate, error: %v", body)
	}

	return createdCertificate, nil
}

func (certificateClient *certificateClient) DeleteById(id string) error {
	r, body, errs := newDelete(certificateClient.config, CertificatesPath+id).End()
	if errs != nil {
		return fmt.Errorf("could not delete certificate, result: %v error: %v", r, errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return fmt.Errorf("not authorised, message from kong: %s", body)
	}

	return nil
}

func (certificateClient *certificateClient) List() (*Certificates, error) {
	r, body, errs := newGet(certificateClient.config, CertificatesPath).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get certificates, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	certificates := &Certificates{}
	err := json.Unmarshal([]byte(body), certificates)
	if err != nil {
		return nil, fmt.Errorf("could not parse certificates list response, error: %v", err)
	}

	return certificates, nil
}

func (certificateClient *certificateClient) UpdateById(id string, certificateRequest *CertificateRequest) (*Certificate, error) {
	r, body, errs := newPatch(certificateClient.config, CertificatesPath+id).Send(certificateRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not update certificate, error: %v", errs)
	}

	if r.StatusCode == 401 || r.StatusCode == 403 {
		return nil, fmt.Errorf("not authorised, message from kong: %s", body)
	}

	updatedCertificate := &Certificate{}
	err := json.Unmarshal([]byte(body), updatedCertificate)
	if err != nil {
		return nil, fmt.Errorf("could not parse certificate update response, error: %v", err)
	}

	if updatedCertificate.Id == nil {
		return nil, fmt.Errorf("could not update certificate, error: %v", body)
	}

	return updatedCertificate, nil
}
