package haiconf

type CommandArgs map[string]interface{}

type Commander interface {
	SetDefault() error
	SetUserConfig(CommandArgs) error
	Run() error
}
