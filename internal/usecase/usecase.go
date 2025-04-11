package usecase

import "github.com/dinhtatuanlinh/video/internal/usecase/video"

type UseCase struct {
	UseCaseVideo video.VideoRepository
}

func NewUseCase(
	useCaseVideo video.VideoRepository) UseCase {
	return UseCase{
		UseCaseVideo: useCaseVideo,
	}
}
