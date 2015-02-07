package shared

type AsDoormanUpdatePayloader interface {
	AsDoormanUpdatePayload() ([]byte, error)
}
