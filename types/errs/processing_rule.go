package errs

const (
	InvalidProcessingRuleID                 = InvalidArgumentError("invalid processing rule ID")
	InvalidProcessingRuleMatch              = InvalidArgumentError("invalid processing rule match")
	InvalidProcessingRuleLanguage           = InvalidArgumentError("invalid processing rule language")
	InvalidProcessingRuleActionDesc         = InvalidArgumentError("invalid processing rule action desc")
	InvalidProcessingRuleActionKind         = InvalidArgumentError("invalid processing rule action kind")
	InvalidProcessingRuleActionSelectorKind = InvalidArgumentError("invalid processing rule action selector kind")
	InvalidProcessingRuleActionSelectorOp   = InvalidArgumentError("invalid processing rule action selector op")
	InvalidProcessingRuleActionSelectorExpr = InvalidArgumentError("invalid processing rule action selector expr")
	InvalidProcessingRuleActionArgs         = InvalidArgumentError("invalid processing rule action args")
	ProcessingRuleNotFound                  = NotFoundError("processing rule not found")
)
