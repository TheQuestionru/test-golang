package logger

import "github.com/ivankorobkov/di"

func TestModule(m *di.Module) {
	m.Import(Module)
	m.AddConstructor(NewTestConfig)
	m.AddConstructor(NewTest)
	m.MarkPackageDep(Config{})
}

type TestLogger struct {
	Logger
}

func NewTest(logger Logger) TestLogger {
	return TestLogger{logger}
}

func NewTestConfig() Config {
	return Config{}
}
