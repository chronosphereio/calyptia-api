package errs

const (
	InvalidAWSMarketplaceCustomerToken = InvalidArgumentError("invalid AWS marketplace customer token")
	AWSContractAlreadyExists           = ConflictError("AWS contract already exists")
)
