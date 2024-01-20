package router

import (
	userapp "gdsc/baro/app/user/pb"
	videoapp "gdsc/baro/app/video/pb"

	userrepo "gdsc/baro/app/user/repositories"
	videorepo "gdsc/baro/app/video/repositories"

	userpb "gdsc/baro/protos/user"
	videopb "gdsc/baro/protos/video"

	"gdsc/baro/global/auth"
	"gdsc/baro/global/utils"

	"google.golang.org/grpc"
)

func NewInitApp(userRepository userrepo.UserRepositoryInterface, userUtil utils.UserUtilInterface, videoRepository videorepo.VideoRepositoryInterface) *grpc.Server {
	server := grpc.NewServer(grpc.UnaryInterceptor(auth.UnaryAuthInterceptor))

	userPbApp := userapp.NewUserPbApp(userRepository, userUtil)
	videoPbApp := videoapp.NewVideoPbApp(videoRepository)

	userpb.RegisterUserServiceServer(server, userPbApp)
	videopb.RegisterVideoServiceServer(server, videoPbApp)

	return server
}
