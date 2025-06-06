package seed

import (
	"bus_booking_api/internal/config"
	"bus_booking_api/internal/models"
	"bus_booking_api/internal/repositories"
	"bus_booking_api/internal/services"
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RunSeeders là hàm chính để thực thi tất cả các seeder.
func RunSeeders(
	cfg *config.Config,
	userRepo repositories.UserRepository,
	companyRepo repositories.CompanyRepository,
	vehicleRepo repositories.VehicleRepository,
	tripRepo repositories.TripRepository,
	bookingRepo repositories.BookingRepository,
	tripService services.TripService,
) {
	log.Println("Bắt đầu seeding dữ liệu (nếu cần thiết)...")
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano()) // Khởi tạo seed cho random

	// Seed Companies
	seededCompanies := seedCompanies(ctx, companyRepo)
	if len(seededCompanies) == 0 {
		log.Println("Không có nhà xe nào được seed hoặc đã tồn tại. Dừng seeding.")
		return
	}

	// Seed Vehicles for each company
	var allSeededVehicles []*models.Vehicle
	for _, company := range seededCompanies {
		vehicles := seedVehiclesForCompany(ctx, vehicleRepo, company)
		allSeededVehicles = append(allSeededVehicles, vehicles...)
	}
	if len(allSeededVehicles) == 0 {
		log.Println("Không có xe nào được seed hoặc đã tồn tại. Dừng seeding trips.")
		return
	}

	// Seed Trips
	seedMultipleTrips(ctx, tripRepo, tripService, seededCompanies, allSeededVehicles)

	log.Println("Hoàn thành seeding dữ liệu.")
}

func seedCompanies(ctx context.Context, companyRepo repositories.CompanyRepository) []*models.Company {
	var seededCompanies []*models.Company
	companiesToSeed := []models.Company{
		{
			Name: "Phương Trang - FUTA Bus Lines", Code: "PHUONGTRANG", Address: "TP. Hồ Chí Minh", Phone: "19006067", Email: "hotro@futabus.vn",
			Description: "Nhà xe Phương Trang uy tín chất lượng.", LogoURL: "https://storage.googleapis.com/vex-config/production-masstransit-operator-logo/1652931009762/futa.png", IsActive: true,
		},
		{
			Name: "Thành Bưởi", Code: "THANHBUOI", Address: "TP. Hồ Chí Minh", Phone: "19006079", Email: "hotro@thanhbuoi.com.vn",
			Description: "Nhà xe Thành Bưởi.", LogoURL: "https://vcdn1-vnexpress.vnecdn.net/2023/11/03/thanh-buoi-jpeg-1698992969-3819-1698993057.jpg?w=0&h=0&q=100&dpr=2&fit=crop&s=pRhD4gvyTn4Xw87dJk4G8A", IsActive: true,
		},
		{
			Name: "Kumho Samco", Code: "KUMHO", Address: "TP. Hồ Chí Minh", Phone: "19006089", Email: "hotro@kumhosamco.com.vn",
			Description: "Dịch vụ vận chuyển hành khách chất lượng cao.", LogoURL: "https://static.vexere.com/production/images/1649649998959.jpeg", IsActive: true,
		},
		{
			Name: "Hoàng Long Limousine", Code: "HOANGLONG", Address: "Hà Nội", Phone: "19001198", Email: "cskh@hoanglongasia.com",
			Description: "Xe limousine cao cấp tuyến Hà Nội - Sài Gòn và các tỉnh.", LogoURL: "https://hoanglongasia.vn/wp-content/uploads/2022/04/logo-hoang-long-bus-2.png", IsActive: true,
		},
	}

	for _, compData := range companiesToSeed {
		existing, err := companyRepo.GetCompanyByCode(ctx, compData.Code)
		if err != nil && err != mongo.ErrNoDocuments {
			log.Printf("Lỗi khi kiểm tra company %s: %v", compData.Code, err)
			continue
		}
		if existing == nil {
			created, errCr := companyRepo.CreateCompany(ctx, &compData)
			if errCr != nil {
				log.Printf("Lỗi khi tạo company %s: %v", compData.Name, errCr)
			} else {
				log.Printf("Đã tạo company: %s (ID: %s)", created.Name, created.ID.Hex())
				seededCompanies = append(seededCompanies, created)
			}
		} else {
			log.Printf("Company %s đã tồn tại.", existing.Name)
			seededCompanies = append(seededCompanies, existing)
		}
	}
	return seededCompanies
}

func seedVehiclesForCompany(ctx context.Context, vehicleRepo repositories.VehicleRepository, company *models.Company) []*models.Vehicle {
	var seededVehicles []*models.Vehicle
	if company == nil || company.ID == primitive.NilObjectID {
		return seededVehicles
	}

	vehicleTypesToSeed := []struct {
		Type        string
		Description string
		SeatMap     models.VehicleSeatMap
		TotalSeats  int
	}{
		{
			Type: "Giường nằm 40 chỗ", Description: "Wifi, Nước uống, Điều hòa, TV", TotalSeats: 40,
			SeatMap: models.VehicleSeatMap{
				Rows: 10, Cols: 4, Layout: layoutGenerator(10, 4, []string{"A", "B", "", "C", "D"}), Legend: map[string]string{"": "Lối đi"},
			},
		},
		{
			Type: "Ghế ngồi 29 chỗ", Description: "Điều hòa, Nước uống", TotalSeats: 29,
			SeatMap: models.VehicleSeatMap{ // Ví dụ 2-2
				Rows: 8, Cols: 4, Layout: layoutGenerator(8, 4, []string{"A", "B", "", "C", "D"}, 29), Legend: map[string]string{"": "Lối đi"},
			},
		},
		{
			Type: "Limousine 9 chỗ", Description: "Ghế massage, Wifi, Sạc USB, TV màn hình lớn", TotalSeats: 9,
			SeatMap: models.VehicleSeatMap{ // Ví dụ 2-1-2 (3 hàng đầu), 0-3-0 (hàng cuối)
				Rows: 4, Cols: 3, Layout: layoutGenerator(4, 3, []string{"A", "B", "C"}, 9, true), Legend: map[string]string{"": "Lối đi"},
			},
		},
	}

	for _, vt := range vehicleTypesToSeed {
		// Kiểm tra xem vehicle type này đã tồn tại cho company chưa
		existingVehicles, _ := vehicleRepo.GetVehiclesByCompanyID(ctx, company.ID)
		alreadyExists := false
		var existingVehicle *models.Vehicle
		for i := range existingVehicles { // Phải dùng i để lấy con trỏ
			if existingVehicles[i].Type == vt.Type {
				alreadyExists = true
				existingVehicle = &existingVehicles[i]
				break
			}
		}

		if !alreadyExists {
			vehicleData := models.Vehicle{
				CompanyID:   company.ID,
				Type:        vt.Type,
				Description: vt.Description,
				SeatMap:     vt.SeatMap,
				TotalSeats:  vt.TotalSeats,
			}
			createdVehicle, err := vehicleRepo.CreateVehicle(ctx, &vehicleData)
			if err != nil {
				log.Printf("Lỗi khi tạo vehicle %s cho %s: %v", vt.Type, company.Name, err)
			} else {
				log.Printf("Đã tạo vehicle: %s (ID: %s) cho %s", createdVehicle.Type, createdVehicle.ID.Hex(), company.Name)
				seededVehicles = append(seededVehicles, createdVehicle)
			}
		} else {
			log.Printf("Vehicle: %s của %s đã tồn tại.", vt.Type, company.Name)
			seededVehicles = append(seededVehicles, existingVehicle)
		}
	}
	return seededVehicles
}

// layoutGenerator tạo sơ đồ ghế tự động
// maxSeats: giới hạn số ghế cho các loại xe ít ghế như 29 chỗ
// isLimousine: cờ đặc biệt cho layout limousine
func layoutGenerator(rows, colsPerLogicalRow int, colLabels []string, maxSeats ...interface{}) [][]interface{} {
	layout := make([][]interface{}, rows)
	seatCount := 0
	actualMaxSeats := 0
	if len(maxSeats) > 0 {
		if ms, ok := maxSeats[0].(int); ok {
			actualMaxSeats = ms
		}
	}
    
    isLimo := false
    if len(maxSeats) > 1 {
        if limoFlag, ok := maxSeats[1].(bool); ok {
            isLimo = limoFlag
        }
    }


	for r := 0; r < rows; r++ {
		layout[r] = make([]interface{}, len(colLabels))
		for c, label := range colLabels {
			if label == "" { // Lối đi
				layout[r][c] = nil
			} else {
				if actualMaxSeats > 0 && seatCount >= actualMaxSeats {
					layout[r][c] = nil // Đã đủ ghế
					continue
				}
                
                seatNumStr := ""
                if isLimo { // Logic đặc biệt cho Limousine 9 chỗ (ví dụ)
                    // Hàng 1, 2, 3: A, B, C (ghế đơn)
                    // Hàng 4: A, B, C (ghế băng)
                    // Điều chỉnh logic này cho phù hợp với sơ đồ limousine 9 chỗ thực tế
                    // Ví dụ:
                    // [A1, B1]
                    // [A2, B2]
                    // [A3, B3, C3]
                    // [A4, B4, C4]
                    // Sơ đồ của bạn là 2-1-2 -> A, B | C | D, E
                    // hoặc 3 hàng đầu 1-1-1, hàng cuối 3 ghế
                    // Mã hiện tại đang là A, B, C cho mỗi hàng nếu colsPerLogicalRow = 3
                    seatNumStr = fmt.Sprintf("%s%d", label, r+1)

                } else {
				    seatNumStr = fmt.Sprintf("%s%d", label, r+1)
                }

				layout[r][c] = seatNumStr
				seatCount++
			}
		}
	}
	return layout
}


func seedMultipleTrips(ctx context.Context, tripRepo repositories.TripRepository, tripService services.TripService,
	companies []*models.Company, vehicles []*models.Vehicle) {

	if len(companies) == 0 || len(vehicles) == 0 {
		log.Println("Không thể seed trips: thiếu thông tin công ty hoặc xe.")
		return
	}

	routes := []struct {
		From models.StopPoint
		To   models.StopPoint
		DurationHours int
	}{
		{From: models.StopPoint{Name: "Sài Gòn", Location: models.GeoJSONPoint{Type: "Point", Coordinates: []float64{106.660172, 10.762622}}},
			To:   models.StopPoint{Name: "Đà Lạt", Location: models.GeoJSONPoint{Type: "Point", Coordinates: []float64{108.436390, 11.940419}}}, DurationHours: 7},
		{From: models.StopPoint{Name: "Đà Lạt", Location: models.GeoJSONPoint{Type: "Point", Coordinates: []float64{108.436390, 11.940419}}},
			To:   models.StopPoint{Name: "Sài Gòn", Location: models.GeoJSONPoint{Type: "Point", Coordinates: []float64{106.660172, 10.762622}}}, DurationHours: 7},
		{From: models.StopPoint{Name: "Sài Gòn", Location: models.GeoJSONPoint{Type: "Point", Coordinates: []float64{106.660172, 10.762622}}},
			To:   models.StopPoint{Name: "Nha Trang", Location: models.GeoJSONPoint{Type: "Point", Coordinates: []float64{109.190996, 12.245833}}}, DurationHours: 8},
		{From: models.StopPoint{Name: "Nha Trang", Location: models.GeoJSONPoint{Type: "Point", Coordinates: []float64{109.190996, 12.245833}}},
			To:   models.StopPoint{Name: "Sài Gòn", Location: models.GeoJSONPoint{Type: "Point", Coordinates: []float64{106.660172, 10.762622}}}, DurationHours: 8},
		{From: models.StopPoint{Name: "Hà Nội", Location: models.GeoJSONPoint{Type: "Point", Coordinates: []float64{105.854444, 21.027778}}},
			To:   models.StopPoint{Name: "Sapa", Location: models.GeoJSONPoint{Type: "Point", Coordinates: []float64{103.843611, 22.336389}}}, DurationHours: 5},
		{From: models.StopPoint{Name: "Sapa", Location: models.GeoJSONPoint{Type: "Point", Coordinates: []float64{103.843611, 22.336389}}},
			To:   models.StopPoint{Name: "Hà Nội", Location: models.GeoJSONPoint{Type: "Point", Coordinates: []float64{105.854444, 21.027778}}}, DurationHours: 5},
	}

	basePrices := []float64{280000, 320000, 350000, 400000, 250000}
	
	// Seed trips cho 7 ngày tới
	numDaysToSeed := 7
	tripsPerDayPerRoute := 2 // Số chuyến mỗi ngày cho mỗi tuyến (có thể random)

	totalTripsCreated := 0

	for dayOffset := 0; dayOffset < numDaysToSeed; dayOffset++ {
		currentDate := time.Now().AddDate(0, 0, dayOffset)
		
		for _, route := range routes {
			for i := 0; i < tripsPerDayPerRoute; i++ {
				// Chọn ngẫu nhiên một nhà xe
				company := companies[rand.Intn(len(companies))]
				
				// Chọn ngẫu nhiên một xe của nhà xe đó
				var companyVehicles []*models.Vehicle
				for _, v := range vehicles {
					if v.CompanyID == company.ID {
						companyVehicles = append(companyVehicles, v)
					}
				}
				if len(companyVehicles) == 0 {
					continue // Nhà xe này không có xe nào được seed
				}
				vehicle := companyVehicles[rand.Intn(len(companyVehicles))]

				// Giờ khởi hành ngẫu nhiên (ví dụ từ 6h đến 22h)
				hour := 6 + rand.Intn(17) // 6 -> 22
				minute := []int{0, 15, 30, 45}[rand.Intn(4)]
				departureTime := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), hour, minute, 0, 0, time.Local)
				expectedArrivalTime := departureTime.Add(time.Duration(route.DurationHours) * time.Hour).Add(time.Duration(rand.Intn(60)-30) * time.Minute) // Thêm chút ngẫu nhiên cho thời gian đến
				
				price := basePrices[rand.Intn(len(basePrices))] + float64(rand.Intn(5)-2)*10000 // Biến động giá nhẹ

				tripData := models.Trip{
					CompanyID: company.ID,
					VehicleID: vehicle.ID,
					Route: models.Route{
						From: route.From,
						To:   route.To,
					},
					DepartureTime:       departureTime,
					ExpectedArrivalTime: expectedArrivalTime,
					Price:               price,
					Status:              models.TripStatusScheduled,
					Seats:               []models.TripSeat{},
				}

				for _, seatRow := range vehicle.SeatMap.Layout {
					for _, seatLabel := range seatRow {
						if seatStr, ok := seatLabel.(string); ok && seatStr != "" {
							tripData.Seats = append(tripData.Seats, models.TripSeat{
								SeatNumber: seatStr,
								Status:     models.SeatStatusAvailable,
							})
						}
					}
				}
				
				_, err := tripService.AddTrip(ctx, tripData)
				if err != nil {
					if !strings.Contains(err.Error(), "đã tồn tại") && !strings.Contains(err.Error(), "duplicate key") {
						log.Printf("Lỗi khi tạo trip %s-%s (%s) vào %s: %v", route.From.Name, route.To.Name, company.Name, departureTime.Format("2006-01-02 15:04"), err)
					}
				} else {
					totalTripsCreated++
					// log.Printf("Đã tạo trip: %s-%s (%s) (ID: %s) vào %s", createdTrip.Route.From.Name, createdTrip.Route.To.Name, company.Name, createdTrip.ID.Hex(), createdTrip.DepartureTime.Format("2006-01-02 15:04"))
				}
				if totalTripsCreated >= 100 { // Giới hạn tổng số trip được tạo
					goto EndSeedTrips
				}
			}
		}
	}
EndSeedTrips:
	log.Printf("Tổng số chuyến đi đã được tạo/xác nhận: %d", totalTripsCreated)
}