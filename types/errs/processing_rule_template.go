package errs

const (
	ProcessingRuleTemplateAlreadyExists          = ConflictError("processing rule template already exists")
	ProcessingRuleTemplateNotFound               = NotFoundError("processing rule template not found")
	InvalidProcessingRuleTemplateName            = InvalidArgumentError("invalid processing rule template name")
	InvalidProcessingRuleTemplateInput           = InvalidArgumentError("invalid processing rule template input")
	InvalidProcessingRuleTemplatePipelineVersion = InvalidArgumentError("invalid processing rule template pipeline version")
	InvalidProcessingRuleDef                     = InvalidArgumentError("invalid processing rule template definition")
	InvalidProcessingRuleTemplateID              = InvalidArgumentError("invalid processing rule template ID")
)
