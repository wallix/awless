package awsspec

type Definition struct {
	Action, Entity, Api         string
	RequiredParams, ExtraParams []string
}

func AWSLookupDefinitions(key string) (t Definition, ok bool) {
	t, ok = AWSTemplatesDefinitions[key]
	return
}
