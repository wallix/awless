package awsservices

type config map[string]interface{}

func (c config) region() string {
	if region, ok := c["aws.region"].(string); ok {
		return region
	}
	return ""
}

func (c config) profile() string {
	if profile, ok := c["aws.profile"].(string); ok {
		return profile
	}
	return ""
}

func (c config) getBool(key string, def bool) bool {
	if b, ok := c[key].(bool); ok {
		return b
	}
	return def
}
