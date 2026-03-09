package repositories

import (
	"attendance-system/backend/models"
	"time"

	"gorm.io/gorm"
)

type AttendanceFilter struct {
	StartDate  *time.Time
	EndDate    *time.Time
	CenterID   *uint
	EmployeeID *uint
}

type DashboardCounts struct {
	TotalEmployees int64 `json:"total_employees"`
	PresentToday   int64 `json:"present_today"`
	AbsentToday    int64 `json:"absent_today"`
	LateArrivals   int64 `json:"late_arrivals"`
}

type CenterSummary struct {
	CenterID       uint   `json:"center_id"`
	CenterName     string `json:"center_name"`
	TotalEmployees int64  `json:"total_employees"`
	PresentToday   int64  `json:"present_today"`
}

type DailyTrend struct {
	Date    string `json:"date"`
	Present int64  `json:"present"`
	Partial int64  `json:"partial"`
}

type AttendanceRepository interface {
	GetByEmployeeAndDate(employeeID uint, date time.Time) (*models.Attendance, error)
	Create(attendance *models.Attendance) error
	Update(attendance *models.Attendance) error
	List(filter AttendanceFilter) ([]models.Attendance, error)
	GetDashboardCounts(centerID *uint, date time.Time, lateAfter time.Time) (DashboardCounts, error)
	GetCenterSummary(date time.Time) ([]CenterSummary, error)
	GetDailyTrend(centerID *uint, startDate, endDate time.Time) ([]DailyTrend, error)
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) GetByEmployeeAndDate(employeeID uint, date time.Time) (*models.Attendance, error) {
	var attendance models.Attendance
	err := r.db.Preload("Employee").Preload("Center").
		Where("employee_id = ? AND date = ?", employeeID, date.Format("2006-01-02")).
		First(&attendance).Error
	if err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *attendanceRepository) Create(attendance *models.Attendance) error {
	return r.db.Create(attendance).Error
}

func (r *attendanceRepository) Update(attendance *models.Attendance) error {
	return r.db.Save(attendance).Error
}

func (r *attendanceRepository) List(filter AttendanceFilter) ([]models.Attendance, error) {
	var records []models.Attendance
	query := r.db.Preload("Employee").Preload("Center").Order("date desc, created_at desc")
	if filter.StartDate != nil {
		query = query.Where("date >= ?", filter.StartDate.Format("2006-01-02"))
	}
	if filter.EndDate != nil {
		query = query.Where("date <= ?", filter.EndDate.Format("2006-01-02"))
	}
	if filter.CenterID != nil {
		query = query.Where("center_id = ?", *filter.CenterID)
	}
	if filter.EmployeeID != nil {
		query = query.Where("employee_id = ?", *filter.EmployeeID)
	}
	err := query.Find(&records).Error
	return records, err
}

func (r *attendanceRepository) GetDashboardCounts(centerID *uint, date time.Time, lateAfter time.Time) (DashboardCounts, error) {
	counts := DashboardCounts{}

	employeeQuery := r.db.Model(&models.Employee{})
	if centerID != nil {
		employeeQuery = employeeQuery.Where("center_id = ?", *centerID)
	}
	if err := employeeQuery.Count(&counts.TotalEmployees).Error; err != nil {
		return counts, err
	}

	attendanceQuery := r.db.Model(&models.Attendance{}).Where("date = ?", date.Format("2006-01-02"))
	if centerID != nil {
		attendanceQuery = attendanceQuery.Where("center_id = ?", *centerID)
	}
	if err := attendanceQuery.Where("status = ?", models.AttendancePresent).Count(&counts.PresentToday).Error; err != nil {
		return counts, err
	}
	counts.AbsentToday = counts.TotalEmployees - counts.PresentToday

	lateQuery := r.db.Model(&models.Attendance{}).
		Where("date = ? AND time_in IS NOT NULL AND time_in > ?", date.Format("2006-01-02"), lateAfter)
	if centerID != nil {
		lateQuery = lateQuery.Where("center_id = ?", *centerID)
	}
	if err := lateQuery.Count(&counts.LateArrivals).Error; err != nil {
		return counts, err
	}

	return counts, nil
}

func (r *attendanceRepository) GetCenterSummary(date time.Time) ([]CenterSummary, error) {
	var rows []CenterSummary
	err := r.db.Raw(`
		SELECT
			c.id AS center_id,
			c.name AS center_name,
			COUNT(DISTINCT e.id) AS total_employees,
			COUNT(DISTINCT CASE WHEN a.status IN (?, ?) THEN a.employee_id END) AS present_today
		FROM centers c
		LEFT JOIN employees e ON e.center_id = c.id
		LEFT JOIN attendances a ON a.center_id = c.id AND a.date = ?
		GROUP BY c.id, c.name
		ORDER BY c.name ASC
	`, models.AttendancePresent, models.AttendancePartial, date.Format("2006-01-02")).Scan(&rows).Error
	return rows, err
}

func (r *attendanceRepository) GetDailyTrend(centerID *uint, startDate, endDate time.Time) ([]DailyTrend, error) {
	var rows []DailyTrend
	query := `
		SELECT
			TO_CHAR(date, 'YYYY-MM-DD') AS date,
			COUNT(CASE WHEN status = ? THEN 1 END) AS present,
			COUNT(CASE WHEN status = ? THEN 1 END) AS partial
		FROM attendances
		WHERE date BETWEEN ? AND ?
	`
	args := []interface{}{models.AttendancePresent, models.AttendancePartial, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")}
	if centerID != nil {
		query += " AND center_id = ?"
		args = append(args, *centerID)
	}
	query += " GROUP BY date ORDER BY date ASC"
	err := r.db.Raw(query, args...).Scan(&rows).Error
	return rows, err
}
