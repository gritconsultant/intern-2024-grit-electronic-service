package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	configs "github.com/komkemkku/komkemkku/Back-end_Grit-Electronic/configs"
	"github.com/komkemkku/komkemkku/Back-end_Grit-Electronic/model"
	"github.com/komkemkku/komkemkku/Back-end_Grit-Electronic/requests"
	"github.com/komkemkku/komkemkku/Back-end_Grit-Electronic/response"
)

var db = configs.Database()

func ListOrderService(ctx context.Context, req requests.OrderRequest) ([]response.OrderResponses, int, error) {
	// คำนวณ offset สำหรับ pagination
	var offset int
	if req.Page > 0 {
		offset = int((req.Page - 1) * req.Size)
	}

	// สร้าง slice สำหรับ response
	resp := []response.OrderResponses{}

	// สร้าง query หลัก
	query := db.NewSelect().
		TableExpr("orders AS o").
		Column("o.id", "o.user_id", "u.username", "o.status", "o.created_at", "o.updated_at", "o.total_price", "o.total_amount").
		ColumnExpr("py.system_bank_id, py.price AS payment_price, py.bank_name, py.account_name, py.account_number, py.status AS payment_status").
		ColumnExpr("s.firstname, s.lastname, s.address, s.zip_code, s.sub_district, s.district, s.province, s.status AS shipment_status").
		Join("LEFT JOIN users AS u ON u.id = o.user_id"). 
		Join("LEFT JOIN payments AS py ON py.id = o.payment_id").
		Join("LEFT JOIN shipments AS s ON s.id = o.shipment_id")

	// เงื่อนไขการค้นหา
	if req.Search != "" {
		query.Where("o.status ILIKE ?", "%"+req.Search+"%")
	}

	// สร้าง query สำหรับนับจำนวนทั้งหมด
	countQuery := db.NewSelect().
		TableExpr("orders AS o")
	if req.Search != "" {
		countQuery.Where("o.status ILIKE ?", "%"+req.Search+"%")
	}
	total, err := countQuery.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders: %v", err)
	}

	// ดึงข้อมูลพร้อม pagination
	err = query.Offset(offset).Limit(int(req.Size)).Scan(ctx, &resp)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch orders: %v", err)
	}

	// ส่ง response กลับ
	return resp, total, nil
}

func GetByIdOrderService(ctx context.Context, orderID int64) (*response.OrderResponses, error) {
	// ตรวจสอบว่าผู้ใช้งานมีอยู่ในระบบหรือไม่
	exists, err := db.NewSelect().Table("orders").Where("user_id = ?", orderID).Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	if !exists {
		return nil, errors.New("user not found")
	}

	// สร้าง response object
	order := &response.OrderResponses{}

	err = db.NewSelect().
		TableExpr("orders AS o").
		Column("o.id", "o.user_id", "o.status", "o.created_at", "o.updated_at").
		ColumnExpr("COALESCE(SUM(ci.total_product_amount), 0) AS total_amount").
		ColumnExpr("COALESCE(SUM(p.price * ci.total_product_amount), 0) AS total_price").
		ColumnExpr(`py.system_bank_id, py.price AS payment_price, py.bank_name, py.account_name, py.account_number, py.status AS payment_status`).
		ColumnExpr(`s.firstname, s.lastname, s.address, s.zip_code, s.sub_district, s.district, s.province, s.status AS shipment_status`).
		ColumnExpr(`u.id AS user_id, u.username, u.email, u.firstname, u.lastname, u.phone`).
		Join("LEFT JOIN products AS p ON p.id = ci.product_id").
		Join("LEFT JOIN payments AS py ON py.id = o.payment_id").
		Join("LEFT JOIN shipments AS s ON s.id = o.shipment_id").
		Join("LEFT JOIN users AS u ON u.id = o.user_id").
		Where("o.user_id = ?", orderID).
		Scan(ctx, order)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no orders found for user_id: %d", orderID)
		}
		return nil, fmt.Errorf("failed to fetch order details: %w", err)
	}

	return order, nil
}

func CreateOrderService(ctx context.Context, req requests.OrderCreateRequest) (*model.Orders, error) {
	var cartID int64
	fmt.Printf("Finding cart with user_id: %d\n", req.UserID)
	if err := db.NewSelect().Table("carts").Column("id").Where("user_id = ?", req.UserID).Scan(ctx, &cartID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no cart found for user_id: %d", req.UserID)
		}
		return nil, fmt.Errorf("failed to find cart: %v", err)
	}
	fmt.Printf("Found cart ID: %d\n", cartID)

	var cartItems []struct {
		ProductID   int64   `json:"product_id"`
		ProductName string  `json:"product_name"`
		Amount      int64   `json:"amount"`
		Price       float64 `json:"price"`
		Stock       int64   `json:"stock"`
	}
	if err := db.NewSelect().
		Table("cart_items").
		ColumnExpr("cart_items.product_id, products.name AS product_name, cart_items.total_product_amount AS amount, products.price, products.stock").
		Join("JOIN products ON products.id = cart_items.product_id").
		Where("cart_id = ?", cartID).
		Scan(ctx, &cartItems); err != nil {
		return nil, fmt.Errorf("failed to fetch cart items: %v", err)
	}

	// ตรวจสอบว่าสินค้าในสต็อกเพียงพอหรือไม่
	for _, item := range cartItems {
		if item.Amount > item.Stock {
			return nil, fmt.Errorf("not enough stock for product %s", item.ProductName)
		}
	}

	totalPrice := 0.0
	totalAmount := 0
	for _, item := range cartItems {
		totalPrice += item.Price * float64(item.Amount)
		totalAmount += int(item.Amount)
	}

	order := &model.Orders{
		UserID:       req.UserID,
		ShipmentID:   req.ShipmentID,
		PaymentID:    req.PaymentID,
		Total_price:  totalPrice,
		Total_amount: totalAmount,
		Status:       "pending",
	}
	order.SetCreatedNow()
	order.SetUpdateNow()

	if _, err := db.NewInsert().Model(order).Returning("id").Exec(ctx); err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}

	for _, item := range cartItems {
		// ลดจำนวน stock ของสินค้า
		if _, err := db.NewUpdate().
			Table("products").
			Set("stock = stock - ?", item.Amount).
			Where("id = ?", item.ProductID).
			Exec(ctx); err != nil {
			return nil, fmt.Errorf("failed to update stock for product %s: %v", item.ProductName, err)
		}

		orderDetail := &model.OrderDetail{
			OrderID:            order.ID,
			ProductName:        item.ProductName,
			TotalProductPrice:  item.Price * float64(item.Amount),
			TotalProductAmount: int(item.Amount),
		}
		if _, err := db.NewInsert().Model(orderDetail).Exec(ctx); err != nil {
			return nil, fmt.Errorf("failed to create order detail: %v", err)
		}
	}

	if _, err := db.NewDelete().Table("cart_items").Where("cart_id = ?", cartID).Exec(ctx); err != nil {
		return nil, fmt.Errorf("failed to delete cart items: %v", err)
	}
	if _, err := db.NewDelete().Table("carts").Where("id = ?", cartID).Exec(ctx); err != nil {
		return nil, fmt.Errorf("failed to delete cart: %v", err)
	}

	return order, nil
}

func UpdateOrderService(ctx context.Context, id int64, req requests.OrderUpdateRequest) (*model.Orders, error) {
	// ตรวจสอบว่า order มีอยู่ในฐานข้อมูลหรือไม่
	exists, err := db.NewSelect().
		TableExpr("orders").
		Where("id = ?", id).
		Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if order exists: %v", err)
	}
	if !exists {
		return nil, errors.New("order not found")
	}

	// ดึงข้อมูล order
	order := &model.Orders{}
	err = db.NewSelect().
		Model(order).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order: %v", err)
	}

	// อัปเดตข้อมูล
	if req.Status != "" {
		order.Status = req.Status
	}
	if req.PaymentID != 0 {
		order.PaymentID = req.PaymentID
	}
	if req.ShipmentID != 0 {
		order.ShipmentID = req.ShipmentID
	}
	if req.CartID != 0 {
	}
	order.SetUpdateNow() // ตั้งค่า UpdatedAt

	// บันทึกข้อมูลกลับไปยังฐานข้อมูล
	_, err = db.NewUpdate().
		Model(order).
		Column("status", "payment_id", "shipment_id", "cart_id", "updated_at").
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %v", err)
	}

	return order, nil
}

func DeleteOrderService(ctx context.Context, id int64) error {
	ex, err := db.NewSelect().TableExpr("orders").Where("id=?", id).Exists(ctx)

	if err != nil {
		return err
	}

	if !ex {
		return errors.New("order not found")
	}

	_, err = db.NewDelete().TableExpr("orders").Where("id =?", id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
