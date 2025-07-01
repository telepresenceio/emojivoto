package api

import (
	"context"

	"google.golang.org/grpc"

	"github.com/telepresenceio/emojivoto/emojivoto-emoji-svc/emoji"
	pb "github.com/telepresenceio/emojivoto/emojivoto-emoji-svc/gen/proto"
)

type EmojiServiceServer struct {
	allEmoji emoji.AllEmoji
	pb.UnimplementedEmojiServiceServer
}

func (svc *EmojiServiceServer) ListAll(context.Context, *pb.ListAllEmojiRequest) (*pb.ListAllEmojiResponse, error) {
	es := svc.allEmoji.List()

	list := make([]*pb.Emoji, 0)
	for _, e := range es {
		pbE := pb.Emoji{
			Unicode:   e.Unicode,
			Shortcode: e.Shortcode,
		}
		list = append(list, &pbE)
	}

	return &pb.ListAllEmojiResponse{List: list}, nil
}

func (svc *EmojiServiceServer) FindByShortcode(_ context.Context, req *pb.FindByShortcodeRequest) (*pb.FindByShortcodeResponse, error) {
	var pbE *pb.Emoji
	foundEmoji := svc.allEmoji.WithShortcode(req.Shortcode)
	if foundEmoji != nil {
		pbE = &pb.Emoji{
			Unicode:   foundEmoji.Unicode,
			Shortcode: foundEmoji.Shortcode,
		}
	}
	return &pb.FindByShortcodeResponse{
		Emoji: pbE,
	}, nil
}

func NewGrpServer(grpcServer *grpc.Server, allEmoji emoji.AllEmoji) {
	pb.RegisterEmojiServiceServer(grpcServer, &EmojiServiceServer{
		allEmoji,
		pb.UnimplementedEmojiServiceServer{},
	})
}
