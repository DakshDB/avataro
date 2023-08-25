package domain

type AvataroUsecase interface {
	GetAvataro(text string, set *string, backgroundSet *string) []byte
}
