package lsp

type Unimplemented struct{}

func (*Unimplemented) UnmarshalJSON([]byte) error {
	panic("Unimplemented!")
}
