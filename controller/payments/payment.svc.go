package payments

import (
	"context"
	"errors"

	configs "github.com/komkemkku/komkemkku/Back-end_Grit-Electronic/configs"
	"github.com/komkemkku/komkemkku/Back-end_Grit-Electronic/controller/image"
	"github.com/komkemkku/komkemkku/Back-end_Grit-Electronic/model"
	"github.com/komkemkku/komkemkku/Back-end_Grit-Electronic/requests"
	"github.com/komkemkku/komkemkku/Back-end_Grit-Electronic/response"
)

var db = configs.Database()

func ListPaymentService(ctx context.Context, req requests.PaymentRequest) ([]response.PaymentResponses, int, error) {

	var Offset int64
	if req.Page > 0 {
		Offset = (req.Page - 1) * req.Size
	}

	resp := []response.PaymentResponses{}

	// สร้าง query
	query := db.NewSelect().
		TableExpr("payments AS p").
		Column("p.id", "p.price", "p.status", "p.updated_by", "p.bank_name", "p.account_name", "p.account_number", "p.created_at", "p.updated_at").
		ColumnExpr("a.id AS admin__id").
		ColumnExpr("a.name AS admin__name").
		ColumnExpr("sb.id AS systembank__id").
		ColumnExpr("sb.bank_name AS systembank__bank_name").
		ColumnExpr("sb.account_name AS systembank__account_name").
		ColumnExpr("sb.account_number AS systembank__account_number").
		ColumnExpr("sb.description AS systembank__description").
		ColumnExpr("i.description AS image").
		Join("LEFT JOIN system_banks AS sb ON sb.id = p.system_bank_id").
		Join("LEFT JOIN admins AS a ON a.id = p.admin_id").
		Join("LEFT JOIN images AS i ON i.ref_id = p.id AND i.type = 'payment_slip'")

	if req.Search != "" {
		query.Where("p.status ILIKE ?", "%"+req.Search+"%")
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Execute query
	err = query.Offset(int(Offset)).Limit(int(req.Size)).Scan(ctx, &resp)
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}

func GetByIdPaymentService(ctx context.Context, id int64) (*response.PaymentResponses, error) {
	ex, err := db.NewSelect().TableExpr("payments").Where("id = ?", id).Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !ex {
		return nil, errors.New("payment not found")
	}
	payment := &response.PaymentResponses{}

	err = db.NewSelect().TableExpr("payments AS p").
		Column("p.id", "p.price", "p.status", "p.updated_by", "p.bank_name", "p.account_name", "p.account_number", "p.created_at", "p.updated_at").
		ColumnExpr("a.id AS admin__id").
		ColumnExpr("a.name AS admin__name").
		ColumnExpr("sb.id AS systembank__id").
		ColumnExpr("sb.bank_name AS systembank__bank_name").
		ColumnExpr("sb.account_name AS systembank__account_name").
		ColumnExpr("sb.account_number AS systembank__account_number").
		ColumnExpr("sb.description AS systembank__description").
		ColumnExpr("i.description AS image").
		Join("LEFT JOIN system_banks AS sb ON sb.id = p.system_bank_id").
		Join("LEFT JOIN admins AS a ON a.id = p.admin_id").
		Join("LEFT JOIN images AS i ON i.ref_id = p.id AND i.type = 'payment_slip'").
		Where("p.id = ?", id).Scan(ctx, payment)
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func CreatePaymentService(ctx context.Context, req requests.PaymentCreateRequest) (*model.Payments, error) {

	payment := &model.Payments{
		Price:         float64(req.Price),
		UpdatedBy:     req.UpdatedBy,
		AdminID:       req.AdminID,
		SystemBankID:  req.SystemBankID,
		Status:        req.Status,
		BankName:      req.BankName,
		AccountName:   req.AccountName,
		AccountNumber: req.AccountNumber,
	}
	payment.SetCreatedNow()
	payment.SetUpdateNow()

	_, err := db.NewInsert().Model(payment).Exec(ctx)
	if err != nil {
		return nil, err
	}

	img := requests.ImageCreateRequest{
		RefID:       payment.ID,
		Type:        "payment_slip",
		Description: req.PaymentSlip,
	}

	_, err = image.CreateImageService(ctx, img)
	if err != nil {
		return nil, err
	}

	return payment, nil

}

func UpdatePaymentService(ctx context.Context, id int64, req requests.PaymentUpdateRequest) (*model.Payments, error) {
	ex, err := db.NewSelect().TableExpr("payments").Where("id=?", id).Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !ex {
		return nil, errors.New("payment not found")
	}

	payment := &model.Payments{}

	err = db.NewSelect().Model(payment).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	payment.Price = float64(req.Price)
	payment.UpdatedBy = req.UpdatedBy
	payment.AdminID = req.AdminID
	payment.SystemBankID = req.SystemBankID
	payment.Status = req.Status
	payment.BankName = req.BankName
	payment.AccountName = req.AccountName
	payment.AccountNumber = req.AccountNumber
	payment.SetUpdateNow()

	_, err = db.NewUpdate().Model(payment).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return nil, err
	}

	img := requests.ImageCreateRequest{
		RefID:       payment.ID,
		Type:        "payment_slip",
		Description: req.PaymentSlip,
	}

	_, err = image.CreateImageService(ctx, img)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func DeletePaymentService(ctx context.Context, id int64) error {
	ex, err := db.NewSelect().TableExpr("payments").Where("id=?", id).Exists(ctx)

	if err != nil {
		return err
	}

	if !ex {
		return errors.New("payment not found")
	}

	_, err = db.NewDelete().TableExpr("payments").Where("id =?", id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
