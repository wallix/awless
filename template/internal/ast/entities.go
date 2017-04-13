package ast

type Entity string

const (
	UnknownEntity Entity = "unknown"
	NoneEntity    Entity = "none"

	Vpc             Entity = "vpc"
	Subnet          Entity = "subnet"
	Instance        Entity = "instance"
	Volume          Entity = "volume"
	Tag             Entity = "tag"
	Secgroup        Entity = "secgroup"
	Keypair         Entity = "keypair"
	Internetgateway Entity = "internetgateway"
	Routetable      Entity = "routetable"
	Route           Entity = "route"
	Loadbalancer    Entity = "loadbalancer"
	Listener        Entity = "listener"
	Targetgroup     Entity = "targetgroup"
	Database        Entity = "database"
	Dbsubnetgroup   Entity = "dbsubnetgroup"

	Zone   Entity = "zone"
	Record Entity = "record"

	User      Entity = "user"
	Group     Entity = "group"
	Role      Entity = "role"
	Policy    Entity = "policy"
	Accesskey Entity = "accesskey"

	Bucket        Entity = "bucket"
	Storageobject Entity = "storageobject"

	Subscription Entity = "subscription"
	Topic        Entity = "topic"
	Queue        Entity = "queue"
)

var entities = map[Entity]struct{}{
	NoneEntity:      {},
	Vpc:             {},
	Subnet:          {},
	Instance:        {},
	Volume:          {},
	Tag:             {},
	Secgroup:        {},
	Keypair:         {},
	Internetgateway: {},
	Routetable:      {},
	Route:           {},
	Loadbalancer:    {},
	Listener:        {},
	Targetgroup:     {},
	Database:        {},
	Dbsubnetgroup:   {},
	Zone:            {},
	Record:          {},
	User:            {},
	Group:           {},
	Role:            {},
	Policy:          {},
	Accesskey:       {},
	Bucket:          {},
	Storageobject:   {},
	Subscription:    {},
	Topic:           {},
	Queue:           {},
}

func IsInvalidEntity(s string) bool {
	_, ok := entities[Entity(s)]
	return !ok
}
