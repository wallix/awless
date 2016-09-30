package main

import (
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/deploy"
	"github.com/wallix/awless/groups"
	"github.com/wallix/awless/provider"
	"github.com/wallix/awless/provider/awsprovider"
	"github.com/wallix/awless/reflectline"
	"github.com/wallix/awless/users"
)

func init() {
	provider.RegisterBuilder("aws", &awsprovider.ProviderBuilder{})
}

func main() {
	c := config.NewConfig("/home/quentin/.awless/default")
	n := reflectline.NewShellGroup()
	n.AddNode("users",
		reflectline.MustReflectGroup(users.NewCommands(c)))
	n.AddNode("groups",
		reflectline.MustReflectGroup(groups.NewCommands(c)))
	n.AddNode("provider",
		reflectline.MustReflectGroup(provider.NewCommands(c)))
	n.AddNode("deploy",
		reflectline.MustReflectGroup(deploy.NewCommands(c)))
	n.AddNode("config",
		reflectline.MustReflectGroup(config.NewCommands(c)))
	reflectline.NewShell(n).Run()
}
