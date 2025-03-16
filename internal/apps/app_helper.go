package apps

func MockTestsApp(mockApp *FizzbuzzApp) {
	App = func() *FizzbuzzApp { return mockApp }
}
