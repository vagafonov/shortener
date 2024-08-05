package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/vagafonov/shortener/internal/container"
	"github.com/vagafonov/shortener/internal/customerror"
	"github.com/vagafonov/shortener/internal/interceptor"
	pb "github.com/vagafonov/shortener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Application Contains routes and starts the server.
type GrpcApplication struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedUsersServer
	cnt *container.Container
}

// Constructor for application.
func NewGrpcApplication(cnt *container.Container) *GrpcApplication {
	return &GrpcApplication{
		cnt: cnt,
	}
}

// Serve run server.
func (a *GrpcApplication) Serve(ctx context.Context) error {
	// TODO move to config
	listen, err := net.Listen("tcp", ":3200") //nolint:gosec
	if err != nil {
		log.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor.AuthInterceptor))
	// регистрируем сервис
	pb.RegisterUsersServer(s, a)

	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (a *GrpcApplication) CreateShortURL(
	ctx context.Context,
	in *pb.CreateShortUrlRequest,
) (*pb.CreateShortUrlResponse, error) {
	var err error
	userID := uuid.Nil
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("userID")
		if len(values) > 0 {
			// ключ содержит слайс строк, получаем первую строку
			userID, err = uuid.Parse(values[0])
			if err != nil {
				a.cnt.GetLogger().Err(err).Msg("cannot covert string to uuid userID from metadata")

				return nil, status.Errorf(codes.InvalidArgument, "cannot convert string %s to uuid", values[0])
			}
		}
	}

	shortURL, err := a.cnt.GetServiceURL().MakeShortURL(
		ctx,
		in.GetUrl(),
		a.cnt.GetConfig().ShortURLLength,
		userID,
	)
	if err != nil {
		if errors.Is(err, customerror.ErrURLAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "url already exists")
		} else {
			a.cnt.GetLogger().Err(err).Msg("cannot make short url")

			return nil, status.Error(codes.Aborted, err.Error())
		}
	}

	var response pb.CreateShortUrlResponse
	response.Result = fmt.Sprintf("%s/%s", a.cnt.GetConfig().ResultURL, shortURL.Short)

	return &response, nil
}
