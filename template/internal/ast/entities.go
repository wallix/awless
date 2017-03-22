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
	Securitygroup   Entity = "securitygroup"
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
	NoneEntity:      struct{}{},
	Vpc:             struct{}{},
	Subnet:          struct{}{},
	Instance:        struct{}{},
	Volume:          struct{}{},
	Tag:             struct{}{},
	Securitygroup:   struct{}{},
	Keypair:         struct{}{},
	Internetgateway: struct{}{},
	Routetable:      struct{}{},
	Route:           struct{}{},
	Loadbalancer:    struct{}{},
	Listener:        struct{}{},
	Targetgroup:     struct{}{},
	Database:        struct{}{},
	Dbsubnetgroup:   struct{}{},
	Zone:            struct{}{},
	Record:          struct{}{},
	User:            struct{}{},
	Group:           struct{}{},
	Role:            struct{}{},
	Policy:          struct{}{},
	Accesskey:       struct{}{},
	Bucket:          struct{}{},
	Storageobject:   struct{}{},
	Subscription:    struct{}{},
	Topic:           struct{}{},
	Queue:           struct{}{},
}

func IsInvalidEntity(s string) bool {
	_, ok := entities[Entity(s)]
	return !ok
}
