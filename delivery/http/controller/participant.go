package controller

import (
	"io"
	"net/http"
	"strconv"

	"iam-service/delivery/http/middleware"
	"iam-service/delivery/http/presenter"
	"iam-service/pkg/errors"
	"iam-service/saving/participant"
	"iam-service/saving/participant/participantdto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ParticipantController struct {
	usecase participant.Usecase
}

func NewParticipantController(uc participant.Usecase) *ParticipantController {
	return &ParticipantController{
		usecase: uc,
	}
}

func (ctrl *ParticipantController) Create(c *fiber.Ctx) error {
	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	var req participantdto.CreateParticipantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if err := validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	appIDStr := c.Get("X-Application-ID")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid X-Application-ID header",
		})
	}

	req.TenantID = tenantID
	req.ApplicationID = appID
	req.UserID = userClaims.UserID

	result, err := ctrl.usecase.CreateParticipant(c.UserContext(), &req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapParticipantResponse(result),
	})
}

func (ctrl *ParticipantController) Get(c *fiber.Ctx) error {
	participantID := c.Params("id")
	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	result, err := ctrl.usecase.GetParticipant(c.UserContext(), participantID, tenantID.String())
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapParticipantResponse(result),
	})
}

func (ctrl *ParticipantController) List(c *fiber.Ctx) error {
	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	req := &participantdto.ListParticipantsRequest{
		TenantID:  tenantID,
		Search:    c.Query("search"),
		Status:    nil,
		Page:      page,
		PerPage:   perPage,
		SortBy:    c.Query("sort_by", "created_at"),
		SortOrder: c.Query("sort_order", "desc"),
	}

	if status := c.Query("status"); status != "" {
		req.Status = &status
	}

	if appIDStr := c.Query("application_id"); appIDStr != "" {
		appID, err := uuid.Parse(appIDStr)
		if err == nil {
			req.ApplicationID = &appID
		}
	}

	if err := validate.Struct(req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	result, err := ctrl.usecase.ListParticipants(c.UserContext(), req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) UpdatePersonalData(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid participant ID",
		})
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	var req participantdto.UpdatePersonalDataRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if err := validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.TenantID = tenantID
	req.ParticipantID = pID
	req.UserID = userClaims.UserID

	result, err := ctrl.usecase.UpdatePersonalData(c.UserContext(), &req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapParticipantResponse(result),
	})
}

func (ctrl *ParticipantController) SaveIdentity(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid participant ID",
		})
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	var req participantdto.SaveIdentityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if err := validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.TenantID = tenantID
	req.ParticipantID = pID

	result, err := ctrl.usecase.SaveIdentity(c.UserContext(), &req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) DeleteIdentity(c *fiber.Ctx) error {
	participantID := c.Params("id")
	identityID := c.Params("identityId")

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	if err := ctrl.usecase.DeleteIdentity(c.UserContext(), identityID, participantID, tenantID.String()); err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (ctrl *ParticipantController) SaveAddress(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid participant ID",
		})
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	var req participantdto.SaveAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if err := validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.TenantID = tenantID
	req.ParticipantID = pID

	result, err := ctrl.usecase.SaveAddress(c.UserContext(), &req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) DeleteAddress(c *fiber.Ctx) error {
	participantID := c.Params("id")
	addressID := c.Params("addressId")

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	if err := ctrl.usecase.DeleteAddress(c.UserContext(), addressID, participantID, tenantID.String()); err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (ctrl *ParticipantController) SaveBankAccount(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid participant ID",
		})
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	var req participantdto.SaveBankAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if err := validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.TenantID = tenantID
	req.ParticipantID = pID

	result, err := ctrl.usecase.SaveBankAccount(c.UserContext(), &req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) DeleteBankAccount(c *fiber.Ctx) error {
	participantID := c.Params("id")
	accountID := c.Params("accountId")

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	if err := ctrl.usecase.DeleteBankAccount(c.UserContext(), accountID, participantID, tenantID.String()); err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (ctrl *ParticipantController) SaveFamilyMember(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid participant ID",
		})
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	var req participantdto.SaveFamilyMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if err := validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.TenantID = tenantID
	req.ParticipantID = pID

	result, err := ctrl.usecase.SaveFamilyMember(c.UserContext(), &req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) DeleteFamilyMember(c *fiber.Ctx) error {
	participantID := c.Params("id")
	memberID := c.Params("memberId")

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	if err := ctrl.usecase.DeleteFamilyMember(c.UserContext(), memberID, participantID, tenantID.String()); err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (ctrl *ParticipantController) SaveEmployment(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid participant ID",
		})
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	var req participantdto.SaveEmploymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if err := validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.TenantID = tenantID
	req.ParticipantID = pID

	result, err := ctrl.usecase.SaveEmployment(c.UserContext(), &req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) SaveBeneficiary(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid participant ID",
		})
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	var req participantdto.SaveBeneficiaryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if err := validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req.TenantID = tenantID
	req.ParticipantID = pID

	result, err := ctrl.usecase.SaveBeneficiary(c.UserContext(), &req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) DeleteBeneficiary(c *fiber.Ctx) error {
	participantID := c.Params("id")
	beneficiaryID := c.Params("beneficiaryId")

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	if err := ctrl.usecase.DeleteBeneficiary(c.UserContext(), beneficiaryID, participantID, tenantID.String()); err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (ctrl *ParticipantController) UploadFile(c *fiber.Ctx) error {
	const maxFileSize = 5 * 1024 * 1024

	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid participant ID",
		})
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	fieldName := c.FormValue("field_name")
	if fieldName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "field_name is required",
		})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "file is required",
		})
	}

	if fileHeader.Size > maxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "file size exceeds 5MB limit",
		})
	}

	allowedContentTypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/gif":       true,
		"application/pdf": true,
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedContentTypes[contentType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "unsupported file type; allowed: jpeg, png, gif, pdf",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to open uploaded file",
		})
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to read uploaded file",
		})
	}
	detectedType := http.DetectContentType(buf[:n])
	if !allowedContentTypes[detectedType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "file content does not match an allowed type; allowed: jpeg, png, gif, pdf",
		})
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to process uploaded file",
		})
	}

	req := &participantdto.UploadFileRequest{
		TenantID:      tenantID,
		ParticipantID: pID,
		FieldName:     fieldName,
	}

	result, err := ctrl.usecase.UploadFile(c.UserContext(), req, file, fileHeader.Size, detectedType, fileHeader.Filename)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) GetStatusHistory(c *fiber.Ctx) error {
	participantID := c.Params("id")

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	result, err := ctrl.usecase.GetStatusHistory(c.UserContext(), participantID, tenantID.String())
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (ctrl *ParticipantController) Submit(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid participant ID",
		})
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	req := &participantdto.SubmitParticipantRequest{
		TenantID:      tenantID,
		ParticipantID: pID,
		UserID:        userClaims.UserID,
	}

	result, err := ctrl.usecase.SubmitParticipant(c.UserContext(), req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapParticipantResponse(result),
	})
}

func (ctrl *ParticipantController) Approve(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid participant ID",
		})
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	req := &participantdto.ApproveParticipantRequest{
		TenantID:      tenantID,
		ParticipantID: pID,
		UserID:        userClaims.UserID,
	}

	result, err := ctrl.usecase.ApproveParticipant(c.UserContext(), req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapParticipantResponse(result),
	})
}

func (ctrl *ParticipantController) Reject(c *fiber.Ctx) error {
	pID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid participant ID",
		})
	}

	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	var body struct {
		Reason string `json:"reason" validate:"required,min=10,max=500"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "invalid request body",
		})
	}

	if err := validate.Struct(&body); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	req := &participantdto.RejectParticipantRequest{
		TenantID:      tenantID,
		ParticipantID: pID,
		UserID:        userClaims.UserID,
		Reason:        body.Reason,
	}

	result, err := ctrl.usecase.RejectParticipant(c.UserContext(), req)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    presenter.MapParticipantResponse(result),
	})
}

func (ctrl *ParticipantController) Delete(c *fiber.Ctx) error {
	participantID := c.Params("id")
	tenantID, err := middleware.GetTenantIDFromContext(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	userClaims, err := middleware.GetMultiTenantClaims(c)
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	err = ctrl.usecase.DeleteParticipant(c.UserContext(), participantID, tenantID.String(), userClaims.UserID.String())
	if err != nil {
		appErr := errors.GetAppError(err)
		return c.Status(appErr.HTTPStatus).JSON(fiber.Map{
			"success": false,
			"error":   appErr.Message,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
