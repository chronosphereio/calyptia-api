package types

type AWSCustomer struct {
	AccountID   string `json:"accountID" yaml:"accountID"`
	Identifier  string `json:"identifier" yaml:"identifier"`
	ProductCode string `json:"productCode" yaml:"productCode"`
}

type CreateAWSContractFromToken struct {
	ProjectID              string `json:"-"`
	AmazonMarketplaceToken string `json:"amazonMarketplaceToken"`
}
