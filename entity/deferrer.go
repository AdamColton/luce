package entity

type Deferrer interface {
	DeferSave(er Referer, saveFn func() error)
	DeferCacheClear(er Referer)
}

type DefaultDeferrer struct{}

func (dd DefaultDeferrer) DeferSave(er Referer, saveFn func() error) {
	saveFn()
}

func (dd DefaultDeferrer) DeferCacheClear(er Referer) {
	// cache is never cleared
}

var DeferStrategy Deferrer = DefaultDeferrer{}
