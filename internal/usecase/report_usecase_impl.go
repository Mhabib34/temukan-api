package usecase

import (
	"context"
	"mime/multipart"
	"temukan-api/internal/dto"
	"temukan-api/internal/exception"
	"temukan-api/internal/helper"
	"temukan-api/internal/model"
	"temukan-api/internal/repository"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportUsecaseImpl struct {
	repo     repository.ReportRepository
	Validate *validator.Validate
	Cld      *cloudinary.Cloudinary
}

func NewReportUsecase(repo repository.ReportRepository, validate *validator.Validate, cld *cloudinary.Cloudinary) ReportUsecase {
	return &ReportUsecaseImpl{
		repo:     repo,
		Validate: validate,
		Cld:      cld,
	}
}

func (u *ReportUsecaseImpl) Create(ctx context.Context, reporterID uuid.UUID, request *dto.CreateReportRequest) (*dto.ReportResponse, error) {
	if err := u.Validate.Struct(request); err != nil {
		return nil, err
	}

	report := &model.Report{
		ReporterID:       reporterID,
		Type:             request.Type,
		Name:             request.Name,
		Gender:           request.Gender,
		EstimatedAge:     request.EstimatedAge,
		Description:      request.Description,
		LastSeenLocation: request.LastSeenLocation,
		City:             request.City,
		Province:         request.Province,
		Latitude:         request.Latitude,
		Longitude:        request.Longitude,
		Status:           model.ReportStatusActive,
	}

	result, err := u.repo.Create(ctx, report)
	if err != nil {
		return nil, err
	}

	response := helper.ToReportResponse(*result)
	return &response, nil
}

func (u *ReportUsecaseImpl) GetAll(ctx context.Context, query dto.GetReportsQuery) (*dto.ReportListData, error) {
	reports, total, err := u.repo.FindAll(ctx, query)
	if err != nil {
		return nil, err
	}

	page := query.Page
	if page < 1 {
		page = 1
	}
	limit := query.Limit
	if limit < 1 {
		limit = 10
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &dto.ReportListData{
		Reports: helper.ToReportResponseList(reports),
		Meta: dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (u *ReportUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*dto.ReportResponse, error) {
	report, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundError("laporan tidak ditemukan")
		}
		return nil, err
	}

	response := helper.ToReportResponse(*report)
	return &response, nil
}

func (u *ReportUsecaseImpl) GetMyReports(ctx context.Context, reporterID uuid.UUID, page, limit int) (*dto.ReportListData, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	reports, total, err := u.repo.FindByReporterID(ctx, reporterID, page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return &dto.ReportListData{
		Reports: helper.ToReportResponseList(reports),
		Meta: dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (u *ReportUsecaseImpl) Update(ctx context.Context, id uuid.UUID, reporterID uuid.UUID, request *dto.UpdateReportRequest) (*dto.ReportResponse, error) {
	report, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundError("laporan tidak ditemukan")
		}
		return nil, err
	}

	// Hanya pelapor yang boleh update
	if report.ReporterID != reporterID {
		return nil, exception.NewForbiddenError("anda tidak memiliki akses ke laporan ini")
	}

	// Patch fields yang dikirim saja
	if request.Name != nil {
		report.Name = request.Name
	}
	if request.Gender != nil {
		report.Gender = *request.Gender
	}
	if request.EstimatedAge != nil {
		report.EstimatedAge = request.EstimatedAge
	}
	if request.Description != nil {
		report.Description = *request.Description
	}
	if request.LastSeenLocation != nil {
		report.LastSeenLocation = *request.LastSeenLocation
	}
	if request.City != nil {
		report.City = *request.City
	}
	if request.Province != nil {
		report.Province = *request.Province
	}
	if request.Latitude != nil {
		report.Latitude = request.Latitude
	}
	if request.Longitude != nil {
		report.Longitude = request.Longitude
	}
	if request.Status != nil {
		report.Status = *request.Status
	}

	result, err := u.repo.Update(ctx, report)
	if err != nil {
		return nil, err
	}

	response := helper.ToReportResponse(*result)
	return &response, nil
}

func (u *ReportUsecaseImpl) Delete(ctx context.Context, id uuid.UUID, reporterID uuid.UUID) error {
	report, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return exception.NewNotFoundError("laporan tidak ditemukan")
		}
		return err
	}

	if report.ReporterID != reporterID {
		return exception.NewForbiddenError("anda tidak memiliki akses ke laporan ini")
	}

	return u.repo.Delete(ctx, id)
}

func (u *ReportUsecaseImpl) UploadPhoto(ctx context.Context, id uuid.UUID, reporterID uuid.UUID, file *multipart.FileHeader) (*dto.PhotoUploadResponse, error) {
	report, err := u.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundError("laporan tidak ditemukan")
		}
		return nil, err
	}

	if report.ReporterID != reporterID {
		return nil, exception.NewForbiddenError("anda tidak memiliki akses ke laporan ini")
	}

	if err := helper.ValidatePhoto(file); err != nil {
		return nil, exception.NewBadRequestError(err.Error())
	}

	photoURL, err := helper.UploadReportPhoto(ctx, u.Cld, file, id.String())
	if err != nil {
		return nil, err
	}

	if err := u.repo.UpdatePhotoURL(ctx, id, photoURL); err != nil {
		return nil, err
	}

	return &dto.PhotoUploadResponse{PhotoURL: photoURL}, nil
}

func (u *ReportUsecaseImpl) GetMapPins(ctx context.Context, query dto.GetMapPinsQuery) (*dto.MapPinsData, error) {
	reports, err := u.repo.FindMapPins(ctx, query)
	if err != nil {
		return nil, err
	}

	pins := helper.ToMapPinResponseList(reports)

	return &dto.MapPinsData{
		Pins:  pins,
		Total: len(pins),
	}, nil
}
