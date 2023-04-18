package testutil

/* Test Utility to store a variable's value in a stack and later
restore it with a single command. Good for storing globals that
we're not sure what are, defer a revert.

Single-threading tests only
*/

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

func RestoreGlobals() {
	for i := len(globalStoreStack) - 1; i >= 0; i-- {
		globalStoreStack[i]()
	}
	globalStoreStack = nil
}
