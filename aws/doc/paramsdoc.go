package awsdoc

func TemplateParamsDoc(templateDef, param string) (string, bool) {
	if doc, ok := manualParamsDoc[templateDef][param]; ok {
		return doc, ok
	}
	doc, ok := generatedParamsDoc[templateDef][param]
	return doc, ok
}

var manualParamsDoc = map[string]map[string]string{
	"createloadbalancer": {
		"scheme": "The routing range of the loadbalancer: Internet-facing or internal.",
	},
}
