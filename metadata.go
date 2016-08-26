package vpp

type MetadataService interface {
}

type metadataService struct {
	client *vppClient
	sToken string
}

func (ms *metadataService) LookupURL() {

}
