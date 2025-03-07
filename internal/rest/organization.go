package rest

import (
	"time"

	"github.com/Comvoca-AI/comvoca-admin-back/internal/service"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrganizationHandler struct {
	OrganizationService service.OrganizationService
}

func NewOrganizationHandler(organizationService service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		OrganizationService: organizationService,
	}
}

func (Organization *OrganizationHandler) Register(app *fiber.App) {
	app.Put("/api/v1/organizations/:id", Organization.addOrUpdateOrganization)
	app.Get("/api/v1/organizations/:id", Organization.getOrganizationById)
	app.Get("/api/v1/organizations/specialty", Organization.getOrganizationSpeciality)
	app.Post("/api/v1/organizations/:organizationId/users", Organization.registerUserRelatedByOrganization)
	app.Get("/api/v1/organizations/:organizationId/users", Organization.getUsersByOrganization)
	app.Get("/api/v1/organizations/user", Organization.getOrginizationByUserId)
}

// UpdateOrganization updates an existing organization by ID.
//
// @Summary      Update an organization
// @Description  Update all fields and relationships of an organization by its ID.
// @Tags         organizations
// @Accept       json
// @Produce      json
// @Param        id path string true "Organization ID"
// @Param        body body types.OrganizationRequest true "Organization Update Data"
// @Success      200 {object} entity.Organization
// @Failure      400 {object} map[string]string "Bad Request"
// @Failure      404 {object} map[string]string "Organization not found"
// @Failure      500 {object} map[string]string "Internal Server Error"
// @Router       /organizations/{id} [put]
func (Organization *OrganizationHandler) addOrUpdateOrganization(c *fiber.Ctx) error {
	orgID := c.Params("id")

	// Parse Request Body
	var updateRequest types.OrganizationRequest
	if err := c.BodyParser(updateRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Call service to update
	updatedOrg, err := Organization.OrganizationService.UpdateOrganization(orgID, updateRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Organization updated successfully", "organization": updatedOrg})

}

// Get Organization by id
// @Summary Get Organization by id
// @Description Get Organization by id
// @Tags Organization
// @Produce  json
// @Param   OrganizationId path string true "Organization ID"
// @Success 200 {object} types.OrganizationResponse
// @Router /api/v1/organizations/{OrganizationId} [get]
func (Organization *OrganizationHandler) getOrganizationById(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	org, err := Organization.OrganizationService.GetOrganizationById(id.String())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(org)
}

// Get Organization speciality
// @Summary Get Organization speciality
// @Description Get Organization speciality
// @Tags Organization
// @Produce  json
// @Success 200 {object} entity.Speciality
// @Router /api/v1/organizations/speciality [get]
func (Organization *OrganizationHandler) getOrganizationSpeciality(c *fiber.Ctx) error {

	/*entity.Speciality.GetFlatSpecialityList()*/
	return nil
}

// registerUserRelatedByOrganization registers a user related to an organization.
// @Summary Register User Related to Organization
// @Description Register a user related to an organization
// @Tags Organization
// @Param organizationId path string true "Organization ID"
// @Produce json
// @Success 201 {object} entity.Speciality
// @Router /api/v1/organizations/{organizationId}/users [post]
// todo implement the logic to register the organization by user
func (Organization *OrganizationHandler) registerUserRelatedByOrganization(c *fiber.Ctx) error {
	// Assuming some logic to register the organization by user
	// Always return status 201
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Organization registered successfully by user",
	})
}

// @Summary Get Users by Organization
// @Description Retrieve users associated with a specific organization
// @Tags Organization
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Success 200 {object} []entity.User
// @Router /api/v1/organizations/{organizationId}/users [get]
func (h *OrganizationHandler) getUsersByOrganization(c *fiber.Ctx) error {
	// Parse organization ID from the URL
	_, err := uuid.Parse(c.Params("organizationId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid organization ID",
		})
	}

	// Return users as JSON
	return c.Status(fiber.StatusOK).JSON(nil)
}

// Helper functions to handle pointers
func uuidPtr(u uuid.UUID) uuid.UUID {
	return u
}

func strPtr(s string) string {
	return s
}

func timePtr(t time.Time) time.Time {
	return t
}



// @Summary Get Organization by User ID
// @Description Retrieve the organization associated with the authenticated user
// @Tags Organization
// @Produce json
// @Security BearerAuth  // Make sure this matches the defined security scheme
// @Success 200 {object} types.OrganizationResponse
// @Failure 400 {object} errors.ErrorResponse "Invalid User ID"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} errors.ErrorResponse "Internal Server Error"
// @Router /api/v1/organizations/user [get]
func (h *OrganizationHandler) getOrginizationByUserId(c *fiber.Ctx) error {
	// Extract user ID from the Bearer token
	_, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid User ID",
		})
	}

	// Return organization as JSON
	return c.Status(fiber.StatusOK).JSON(nil)
}
