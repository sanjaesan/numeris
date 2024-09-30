package server

import (
	"context"

	pb "github.com/numeris/pkg/proto"

	"github.com/numeris/internal/domain"
	"github.com/numeris/internal/repository"
	"github.com/numeris/internal/service"
)

type InvoiceServer struct {
	pb.UnimplementedInvoiceServiceServer
	inv *service.InvoiceService
}

func NewInvoiceServer(inv *service.InvoiceService) *InvoiceServer {
	return &InvoiceServer{
		inv: inv,
	}
}

func (s *InvoiceServer) CreateInvoice(ctx context.Context, req *pb.CreateInvoiceRequest) (*pb.CreateInvoiceResponse, error) {
	resp := &pb.CreateInvoiceResponse{}
	invoice := &domain.Invoice{
		IssueDate:       req.GetIssueDate().AsTime(),
		DueDate:         req.GetDueDate().AsTime(),
		BillingCurrency: req.GetBillingCurrency(),
		Discount:        req.GetDiscount(),
	}
	err := s.inv.Create(ctx, invoice)
	if err != nil {
		switch err {
		case repository.ErrDateRequired:
			return nil, repository.ErrDateRequired
		case repository.ErrCurrencyRequired:
			return nil, repository.ErrCurrencyRequired
		default:
			return nil, repository.ErrInvalidData
		}
	}
	return &pb.CreateInvoiceResponse{
		Code:    resp.GetCode(),
		Message: resp.Message,
	}, nil
}

func (s *InvoiceServer) UpdateInvoice(ctx context.Context, req *pb.UpdateInvoiceRequest) (*pb.UpdateInvoiceResponse, error) {
	resp := &pb.UpdateInvoiceResponse{}

	invoice, err := s.inv.ByID(ctx, uint(req.GetId()))
	if err != nil {
		return nil, err
	}
	invoice.IssueDate = req.GetIssueDate().AsTime()
	invoice.DueDate = req.GetDueDate().AsTime()
	invoice.BillingCurrency = req.GetBillingCurrency()
	invoice.Discount = req.GetDiscount()

	err = s.inv.Update(ctx, invoice)
	if err != nil {
		switch err {
		case repository.ErrDateRequired:
			return nil, repository.ErrDateRequired
		case repository.ErrCurrencyRequired:
			return nil, repository.ErrCurrencyRequired
		default:
			return nil, repository.ErrInvalidData
		}
	}
	return &pb.UpdateInvoiceResponse{
		Success: true,
		Message: resp.Message,
	}, nil
}

func (s *InvoiceServer) DeleteInvoice(ctx context.Context, req *pb.DeleteInvoiceRequest) (*pb.DeleteInvoiceResponse, error) {
	return &pb.DeleteInvoiceResponse{}, nil
}

func (s *InvoiceServer) GetInvoiceByID(ctx context.Context, req *pb.GetInvoiceByIDRequest) (*pb.GetInvoiceByIDResponse, error) {
	return &pb.GetInvoiceByIDResponse{}, nil
}

func (s *InvoiceServer) GetInvoiceByInvoiceNo(ctx context.Context, req *pb.GetInvoiceByInvoiceNoRequest) (*pb.GetInvoiceByInvoiceNoResponse, error) {
	return &pb.GetInvoiceByInvoiceNoResponse{}, nil
}

func (s *InvoiceServer) ListInvoices(ctx context.Context, req *pb.ListInvoicesRequest) (*pb.ListInvoicesResponse, error) {
	return &pb.ListInvoicesResponse{}, nil
}

func (s *InvoiceServer) CreateInvoiceItem(ctx context.Context, req *pb.CreateInvoiceItemRequest) (*pb.CreateInvoiceItemResponse, error) {
	return &pb.CreateInvoiceItemResponse{}, nil
}

func (s *InvoiceServer) GetPaymentInfoByInvoiceID(ctx context.Context, req *pb.GetPaymentInfoByInvoiceIDRequest) (*pb.GetPaymentInfoByInvoiceIDResponse, error) {
	return &pb.GetPaymentInfoByInvoiceIDResponse{}, nil
}
