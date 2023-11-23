package api

import (
	"context"
	"github.com/passsquale/product-item-api/internal/repo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	totalTemplateNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "product_item_api_item_not_found_total",
		Help: "Total number of templates that were not found",
	})
)

type templateAPI struct {
	pb.UnimplementedOmpTemplateApiServiceServer
	repo repo.Repo
}

// NewTemplateAPI returns api of omp-template-api service
func NewTemplateAPI(r repo.Repo) pb.OmpTemplateApiServiceServer {
	return &templateAPI{repo: r}
}

func (o *templateAPI) DescribeTemplateV1(
	ctx context.Context,
	req *pb.DescribeTemplateV1Request,
) (*pb.DescribeTemplateV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribeTemplateV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	template, err := o.repo.DescribeTemplate(ctx, req.TemplateId)
	if err != nil {
		log.Error().Err(err).Msg("DescribeTemplateV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if template == nil {
		log.Debug().Uint64("templateId", req.TemplateId).Msg("template not found")
		totalTemplateNotFound.Inc()

		return nil, status.Error(codes.NotFound, "template not found")
	}

	log.Debug().Msg("DescribeTemplateV1 - success")

	return &pb.DescribeTemplateV1Response{
		Value: &pb.Template{
			Id:  template.ID,
			Foo: template.Foo,
		},
	}, nil
}
