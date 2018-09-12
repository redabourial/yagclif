package argument

type argument struct {
	text string
	used bool
}
type Argument interface {
	Text() string
	Equals(string) bool
	IsUsed() bool
	IsNotUsed() bool
	Use()
}

func (a *argument) Text() string {
	return a.text
}
func (a *argument) Equals(s string) bool {
	return a.text == s
}
func (a *argument) IsUsed() bool {
	return a.used
}
func (a *argument) IsNotUsed() bool {
	return !a.used
}
func (a *argument) Use() {
	a.used = true
}

func New(s string) Argument {
	return &argument{
		text: s,
		used: false,
	}
}
