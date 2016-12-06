package driver

type Token int

const (
	UNKNOWN Token = iota

	actionStart

	CREATE
	DELETE
	VERIFY

	actionEnd

	resourceStart

	INSTANCE
	SUBNET
	VPC

	resourceEnd

	paramStart

	COUNT
	BASE
	TYPE
	CIDR
	REFERENCES
	REF

	paramEnd
)

func (t Token) IsAction() bool {
	return t > actionStart && t < actionEnd
}

func (t Token) IsResource() bool {
	return t > resourceStart && t < resourceEnd
}

func (t Token) IsParam() bool {
	return t > paramStart && t < paramEnd
}

func (t Token) String() string {
	switch t {
	case CREATE:
		return "CREATE"
	case DELETE:
		return "DELETE"
	case VERIFY:
		return "VERIFY"

	case INSTANCE:
		return "INSTANCE"
	case SUBNET:
		return "SUBNET"
	case VPC:
		return "VPC"

	case COUNT:
		return "COUNT"
	case BASE:
		return "BASE"
	case TYPE:
		return "TYPE"
	case CIDR:
		return "CIDR"
	case REFERENCES:
		return "REFERENCES"
	case REF:
		return "REF"

	default:
		return "UNKNOWN"
	}
}

func TokenFromString(s string) Token {
	switch s {
	case "CREATE":
		return CREATE
	case "DELETE":
		return DELETE
	case "VERIFY":
		return VERIFY

	case "INSTANCE":
		return INSTANCE
	case "SUBNET":
		return SUBNET
	case "VPC":
		return VPC

	case "COUNT":
		return COUNT
	case "BASE":
		return BASE
	case "TYPE":
		return TYPE
	case "CIDR":
		return CIDR
	case "REFERENCES":
		return REFERENCES
	case "REF":
		return REF

	default:
		return UNKNOWN
	}
}
