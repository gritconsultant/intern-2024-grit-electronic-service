package users

import (
	"context"
	"errors"
	//"net/http"

	"github.com/gin-gonic/gin"
	config "github.com/komkemkku/komkemkku/Back-end_Grit-Electronic/configs"
	"github.com/komkemkku/komkemkku/Back-end_Grit-Electronic/model"
	"github.com/komkemkku/komkemkku/Back-end_Grit-Electronic/requests"
	"github.com/komkemkku/komkemkku/Back-end_Grit-Electronic/response"
	"github.com/komkemkku/komkemkku/Back-end_Grit-Electronic/utils"
	"golang.org/x/crypto/bcrypt"
)

var db = config.Database()

func GetUserByID(c *gin.Context) {
	id := requests.UserIdRequest{}
	if err := c.BindUri(&id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	data, err := GetByIdUserService(c, int64(id.ID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func CreateUsersService(ctx context.Context, req requests.UserCreateRequest) (*model.Users, error) {
	// ตรวจสอบ Role ID = 3
	ex, err := db.NewSelect().TableExpr("roles").Where("id = ?", 3).Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !ex {
		return nil, errors.New("default role not found (role_id = 3)")
	}

	// แฮชรหัสผ่าน
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password: " + err.Error())
	}

	// สร้างผู้ใช้
	user := &model.Users{
		FirstName: req.Firstname,
		LastName:  req.Lastname,
		Username:  req.Username,
		Password:  hashedPassword, // เก็บ Password ที่ถูกแฮช
		Email:     req.Email,
		Phone:     req.Phone,
	}
	user.SetCreatedNow()
	user.SetUpdateNow()

	_, err = db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return nil, err
	}

	// กำหนด Role ID = 3
	userRole := &model.UserRole{
		UserID: user.ID,
		RoleID: 3,
	}

	_, err = db.NewInsert().Model(userRole).Exec(ctx)
	if err != nil {
		return nil, errors.New("failed to assign role to user: " + err.Error())
	}

	return user, nil
}

// func RegisterUser(c *gin.Context) {
// 	var req requests.UserCreateRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	user, err := CreateUsersService(c.Request.Context(), req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "User registered successfully",
// 		"user":    user, // Password จะถูกข้ามเพราะมี json:"-" ใน Model
// 	})
// }

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func UpdateUserService(ctx context.Context, id int64, req requests.UserUpdateRequest) (*model.Users, error) {

	ex, err := db.NewSelect().TableExpr("roles").Where("id = ?", 3).Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !ex {
		return nil, errors.New("default role not found (role_id = 3)")
	}

	user := &model.Users{}

	err = db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	user.FirstName = req.Firstname
	user.LastName = req.Lastname
	user.Username = req.Username
	user.Email = req.Email
	user.Phone = req.Phone
	user.SetUpdateNow()

	_, err = db.NewUpdate().Model(user).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func DeleteUserService(ctx context.Context, id int64) error {
	ex, err := db.NewSelect().TableExpr("users").Where("id=?", id).Exists(ctx)

	if err != nil {
		return err
	}

	if !ex {
		return errors.New("user not found")
	}

	_, err = db.NewDelete().TableExpr("users").Where("id =?", id).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
