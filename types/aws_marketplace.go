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

type AWSMarketplaceSubscriptionNotification struct {
	Action                 AWSMarketplaceSubscriptionNotificationAction `json:"action"`
	CustomerIdentifier     string                                       `json:"customer-identifier"`
	ProductCode            string                                       `json:"product-code"`
	OfferIdentifier        string                                       `json:"offer-identifier"`
	IsFreeTrialTermPresent bool                                         `json:"isFreeTrialTermPresent"`
}

type AWSMarketplaceSubscriptionNotificationAction string

const (
	AWSMarketplaceSubscriptionNotificationActionSubscribeSuccess   AWSMarketplaceSubscriptionNotificationAction = "subscribe-success"
	AWSMarketplaceSubscriptionNotificationActionSubscribeFail      AWSMarketplaceSubscriptionNotificationAction = "subscribe-fail"
	AWSMarketplaceSubscriptionNotificationActionUnsubscribePending AWSMarketplaceSubscriptionNotificationAction = "unsubscribe-pending"
	AWSMarketplaceSubscriptionNotificationActionUnsubscribeSuccess AWSMarketplaceSubscriptionNotificationAction = "unsubscribe-success"
)
