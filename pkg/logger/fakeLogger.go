package logger

type FakeLogger struct{}

func (f *FakeLogger) Debug(msg string)   {}
func (f *FakeLogger) Info(msg string)    {}
func (f *FakeLogger) Warning(msg string) {}
func (f *FakeLogger) Error(msg string)   {}
func (f *FakeLogger) Fatal(msg string)   {}
func (f *FakeLogger) Close()             {}

func NewFakeLogger() LoggerInterface {
	return &FakeLogger{}
}
