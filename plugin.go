package venom

type IPlugin interface {
	OnStart(*Config)
	OnDestroy(*Config)
}
