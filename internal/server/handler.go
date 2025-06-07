package server

import (
	"github.com/gofiber/fiber/v2"
	attendanceV1 "github.com/vnnyx/employee-management/internal/attendance/delivery/http/v1"
	attendanceRepo "github.com/vnnyx/employee-management/internal/attendance/repository"
	attendanceUseCase "github.com/vnnyx/employee-management/internal/attendance/usecase"
	authV1 "github.com/vnnyx/employee-management/internal/auth/delivery/http/v1"
	authRepo "github.com/vnnyx/employee-management/internal/auth/repository"
	authUseCase "github.com/vnnyx/employee-management/internal/auth/usecase"
	"github.com/vnnyx/employee-management/internal/middleware"
	overtimeV1 "github.com/vnnyx/employee-management/internal/overtime/delivery/http/v1"
	overtimeRepo "github.com/vnnyx/employee-management/internal/overtime/repository"
	overtimeUseCase "github.com/vnnyx/employee-management/internal/overtime/usecase"
	payrollV1 "github.com/vnnyx/employee-management/internal/payroll/delivery/http/v1"
	payrollRepo "github.com/vnnyx/employee-management/internal/payroll/repository"
	payrollUseCase "github.com/vnnyx/employee-management/internal/payroll/usecase"
	reimbursementV1 "github.com/vnnyx/employee-management/internal/reimbursement/delivery/http/v1"
	reimbursementRepo "github.com/vnnyx/employee-management/internal/reimbursement/repository"
	reimbursementUseCase "github.com/vnnyx/employee-management/internal/reimbursement/usecase"
	userRepo "github.com/vnnyx/employee-management/internal/users/repository"
)

func (s *Server) MapHandlers() error {
	health := s.Fiber.Group("/health")
	health.Get("/check", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(
			map[string]any{
				"status": "ok",
			},
		)
	})

	authRepo := authRepo.NewAuthRepository(s.DB)
	attendanceRepo := attendanceRepo.NewAttendanceRepository(s.DB)
	overtimeRepo := overtimeRepo.NewOvertimeRepository(s.DB)
	reimbursementRepo := reimbursementRepo.NewReimbursementRepository(s.DB)
	payrollRepo := payrollRepo.NewPayrollRepository(s.DB)
	userRepo := userRepo.NewUserRepository(s.DB)

	authUC := authUseCase.NewAuthUseCase(authRepo, authUseCase.AuthConfig{
		Key: s.Config.App.Key,
	})
	attendanceUC := attendanceUseCase.NewAttendanceUseCase(attendanceRepo)
	overtimeUC := overtimeUseCase.NewOvertimeUseCase(overtimeRepo)
	reimbursementUC := reimbursementUseCase.NewReimbursementUseCase(reimbursementRepo)
	payrollUC := payrollUseCase.NewPayrollUseCase(
		payrollRepo,
		userRepo,
		attendanceRepo,
		overtimeRepo,
		reimbursementRepo,
	)

	authHandler := authV1.NewAuthHandler(authUC)
	attendanceHandler := attendanceV1.NewAttendanceHandler(attendanceUC)
	overtimeHandler := overtimeV1.NewOvertimeHandler(overtimeUC)
	reimbursementHandler := reimbursementV1.NewReimbursementHandler(reimbursementUC)
	payrollHandler := payrollV1.NewPayrollHandler(payrollUC)

	externalV1 := s.Fiber.Group("/external/api/v1")

	noGuardRoutes := externalV1.Group("")
	authV1.MapAuth(noGuardRoutes, authHandler)

	externalV1.Use(middleware.Auth(s.Config))

	attendanceV1.MapAttendance(externalV1, attendanceHandler)
	overtimeV1.MapOvertime(externalV1, overtimeHandler)
	reimbursementV1.MapReimbursement(externalV1, reimbursementHandler)
	payrollV1.MapPayroll(externalV1, payrollHandler)

	return nil
}
