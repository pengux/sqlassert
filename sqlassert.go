package sqlassert

type testingT interface {
	Errorf(format string, args ...interface{})
}

type nilTestingT struct{}

func (n nilTestingT) Errorf(format string, args ...interface{}) {}

var nilT = new(nilTestingT)
