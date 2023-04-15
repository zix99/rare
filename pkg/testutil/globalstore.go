package testutil

var globalStoreStack []func()

func StoreGlobal[T any](ptr *T) {
	curr := *ptr
	globalStoreStack = append(globalStoreStack, func() {
		*ptr = curr
	})
}

func SwitchGlobal[T any](ptr *T, newVal T) {
	StoreGlobal(ptr)
	*ptr = newVal
}

func RevertGlobals() {
	for i := len(globalStoreStack) - 1; i >= 0; i-- {
		globalStoreStack[i]()
	}
	globalStoreStack = nil
}
