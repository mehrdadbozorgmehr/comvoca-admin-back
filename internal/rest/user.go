package rest

import (
	"github.com/Comvoca-AI/comvoca-admin-back/internal/service"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/types"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (user *UserHandler) Register(app *fiber.App) {
	app.Put("/api/v1/users", user.addOrUpdateUser)
	app.Get("/api/v1/users/:id", user.getUserById)
	app.Get("/api/v1/users/role/:role", user.getUsersByRole)
	app.Get("/api/v1/users/role/:role", user.getUsersByRole)
}

// Add or Update a User such as Admin, Staff, Support
// @Summary Add new User
// @Description Add or Update a User
// @Tags  User
// @Accept  json
// @Produce  json
// @Param   property body types.UserRequest true "UserRequest"
// @Success 200 {object} types.UserResponse
// @Router /api/v1/users [put]
func (u *UserHandler) addOrUpdateUser(c *fiber.Ctx) error {

	return nil

}

// Get User by id
// @Summary Get User by id
// @Description Get User by id
// @Tags User
// @Produce  json
// @Param   user_id path string true "User ID"
// @Success 200 {object} types.UserResponse
// @Router /api/v1/users/{user_id} [get]
func (u *UserHandler) getUserById(c *fiber.Ctx) error {

	return c.JSON(types.UserResponse{})
}

// Get Users by role
// @Summary Get Users by role
// @Description Get Users by role
// @Tags User
// @Produce  json
// @Param   role path string true "Role"
// @Success 200 {array} types.UserResponse
// @Router /api/v1/users/role/{role} [get]
func (u *UserHandler) getUsersByRole(c *fiber.Ctx) error {

	return nil
}

// Add or Update a User such as General User
// @Summary Add new User
// @Description Add or Update a General User
// @Tags  User
// @Accept  json
// @Produce  json
// @Param   property body types.UserRequest true "UserRequest"
// @Success 200 {object} types.UserResponse
// @Router /api/v1/orginization/:orginizationId [post]
func (u *UserHandler) addOrUpdateGenralUser(c *fiber.Ctx) error {

	return nil

}
