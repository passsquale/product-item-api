package api

import (
	"context"
	"github.com/passsquale/product-item-api/internal/repo"
	product_item_api "github.com/passsquale/product-item-api/pkg/product-item-api"
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

type itemAPI struct {
	product_item_api.UnimplementedProductItemApiServiceServer
	repo repo.Repo
}

// NewTemplateAPI returns api of omp-template-api service
func NewTemplateAPI(r repo.Repo) product_item_api.ProductItemApiServiceServer {
	return &itemAPI{repo: r}
}

func (o *itemAPI) DescribeTemplateV1(
	ctx context.Context,
	req *product_item_api.DescribeItemV1Request,
) (*product_item_api.DescribeItemV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribeTemplateV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	item, err := o.repo.DescribeTemplate(ctx, req.ItemID)
	if err != nil {
		log.Error().Err(err).Msg("DescribeTemplateV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if item == nil {
		log.Debug().Uint64("templateId", req.ItemID).Msg("item not found")
		totalTemplateNotFound.Inc()

		return nil, status.Error(codes.NotFound, "item not found")
	}

	log.Debug().Msg("DescribeItemV1 - success")

	return &product_item_api.DescribeItemV1Response{
		Value: &product_item_api.Item{
			ID: item.ID,
		},
	}, nil
}
