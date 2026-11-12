package api

// Tests driven by a fake git or exec have no real checkout for the pre-push
// check to make, so by default it is a no-op there; the tests of the check
// itself (TestShipPrePushCheck) use real git and restore the real one.
func init() {
	prePushCheck = func(*Ctx, string) error { return nil }
}
