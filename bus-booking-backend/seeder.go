// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"time"

// 	"github.com/Go_final_exam/bus-booking-backend/src/config"
// 	"github.com/Go_final_exam/bus-booking-backend/src/models"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// var (
// 	companyCollection *mongo.Collection
// 	vehicleCollection *mongo.Collection
// 	tripCollection    *mongo.Collection
// )

// func main() {
// 	// 1. Tải cấu hình và kết nối DB
// 	cfg, err := config.LoadConfig()
// 	if err != nil {
// 		log.Fatalf("Không thể tải cấu hình: %v", err)
// 	}
// 	config.ConnectDB(cfg)

// 	// Lấy các collection
// 	companyCollection = config.DB.Collection("companies")
// 	vehicleCollection = config.DB.Collection("vehicles")
// 	tripCollection = config.DB.Collection("trips")

// 	log.Println("Bắt đầu quá trình seeding dữ liệu...")

// 	// 2. Xóa dữ liệu cũ (để đảm bảo làm sạch mỗi lần chạy)
// 	clearOldData()

// 	// 3. Seed dữ liệu theo thứ tự
// 	companyIDs := seedCompanies()
// 	vehicleData := seedVehicles(companyIDs)
// 	seedTrips(vehicleData)

// 	log.Println("🎉 Seeding dữ liệu thành công!")
// }

// func clearOldData() {
// 	log.Println("--- Xóa dữ liệu cũ ---")
// 	ctx := context.Background()
// 	companyCollection.Drop(ctx)
// 	vehicleCollection.Drop(ctx)
// 	tripCollection.Drop(ctx)
// 	log.Println("Đã xóa các collection cũ.")
// }

// // seedCompanies tạo dữ liệu cho các nhà xe
// func seedCompanies() []primitive.ObjectID {
// 	log.Println("--- Seeding Companies ---")
// 	ctx := context.Background()

// 	companies := []interface{}{
// 		models.Company{ID: primitive.NewObjectID(), Name: "Phương Trang", Code: "FUTA", Description: "Nhà xe chất lượng cao Phương Trang.", LogoURL: "https://futabus.vn/images/logo/bus-general-futa.png"},
// 		models.Company{ID: primitive.NewObjectID(), Name: "Thành Bưởi", Code: "THANHBUOI", Description: "Dịch vụ vận chuyển hành khách và hàng hóa.", LogoURL: "https://thanhbuoi.com.vn/images/logo.png"},
// 		models.Company{ID: primitive.NewObjectID(), Name: "Kumho Samco", Code: "KUMHO", Description: "Nhà xe uy tín tuyến Sài Gòn - Vũng Tàu.", LogoURL: "https://kumhosamco.com.vn/images/logo-kumho-samco-buslines.png"},
// 	}

// 	result, err := companyCollection.InsertMany(ctx, companies)
// 	if err != nil {
// 		log.Fatalf("Lỗi khi seed companies: %v", err)
// 	}

// 	ids := make([]primitive.ObjectID, len(result.InsertedIDs))
// 	for i, id := range result.InsertedIDs {
// 		ids[i] = id.(primitive.ObjectID)
// 	}

// 	log.Printf("Đã tạo %d companies.", len(ids))
// 	return ids
// }

// // seedVehicles tạo dữ liệu cho các loại xe, liên kết với nhà xe
// func seedVehicles(companyIDs []primitive.ObjectID) map[primitive.ObjectID][]primitive.ObjectID {
// 	log.Println("--- Seeding Vehicles ---")
// 	ctx := context.Background()
// 	vehicleData := make(map[primitive.ObjectID][]primitive.ObjectID)
// 	var vehiclesToInsert []interface{}

// 	for _, companyID := range companyIDs {
// 		// Mỗi nhà xe có 2 loại xe
// 		vehicle1 := models.Vehicle{ID: primitive.NewObjectID(), Type: "Giường nằm 40 chỗ", TotalSeats: 40}
// 		vehicle2 := models.Vehicle{ID: primitive.NewObjectID(), Type: "Limousine 22 chỗ", TotalSeats: 22}

// 		vehiclesToInsert = append(vehiclesToInsert, vehicle1, vehicle2)
// 		vehicleData[companyID] = []primitive.ObjectID{vehicle1.ID, vehicle2.ID}
// 	}

// 	_, err := vehicleCollection.InsertMany(ctx, vehiclesToInsert)
// 	if err != nil {
// 		log.Fatalf("Lỗi khi seed vehicles: %v", err)
// 	}

// 	log.Printf("Đã tạo %d vehicles.", len(vehiclesToInsert))
// 	return vehicleData
// }

// // seedTrips tạo dữ liệu cho các chuyến đi
// func seedTrips(vehicleData map[primitive.ObjectID][]primitive.ObjectID) {
// 	log.Println("--- Seeding Trips ---")
// 	ctx := context.Background()
// 	var tripsToInsert []interface{}

// 	routes := []models.Route{
// 		{From: models.LocationPoint{Name: "TP. Hồ Chí Minh"}, To: models.LocationPoint{Name: "Đà Lạt"}},
// 		{From: models.LocationPoint{Name: "Đà Lạt"}, To: models.LocationPoint{Name: "TP. Hồ Chí Minh"}},
// 		{From: models.LocationPoint{Name: "TP. Hồ Chí Minh"}, To: models.LocationPoint{Name: "Vũng Tàu"}},
// 		{From: models.LocationPoint{Name: "Vũng Tàu"}, To: models.LocationPoint{Name: "TP. Hồ Chí Minh"}},
// 		{From: models.LocationPoint{Name: "Hà Nội"}, To: models.LocationPoint{Name: "Hạ Long"}},
// 		{From: models.LocationPoint{Name: "Hạ Long"}, To: models.LocationPoint{Name: "Hà Nội"}},
// 	}
// 	departureHours := []int{7, 9, 13, 17, 21} // Các giờ khởi hành trong ngày

// 	// Lấy thông tin chi tiết của tất cả các nhà xe và xe
// 	var companies []models.Company
// 	cursor, _ := companyCollection.Find(ctx, bson.M{})
// 	cursor.All(ctx, &companies)

// 	var vehicles []models.Vehicle
// 	cursor, _ = vehicleCollection.Find(ctx, bson.M{})
// 	cursor.All(ctx, &vehicles)

// 	// Tạo map để dễ dàng truy xuất thông tin
// 	companyMap := make(map[primitive.ObjectID]models.Company)
// 	for _, c := range companies {
// 		companyMap[c.ID] = c
// 	}
// 	vehicleMap := make(map[primitive.ObjectID]models.Vehicle)
// 	for _, v := range vehicles {
// 		vehicleMap[v.ID] = v
// 	}

// 	// Tạo chuyến đi cho 7 ngày tới
// 	for day := 0; day < 7; day++ {
// 		for _, companyID := range keys(vehicleData) {
// 			vehicleIDs := vehicleData[companyID]
// 			for _, hour := range departureHours {
// 				// Chọn ngẫu nhiên 1 tuyến và 1 loại xe
// 				route := routes[rand.Intn(len(routes))]
// 				vehicleID := vehicleIDs[rand.Intn(len(vehicleIDs))]
// 				vehicle := vehicleMap[vehicleID]

// 				now := time.Now()
// 				departureTime := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location()).AddDate(0, 0, day)
// 				// Chuyến từ SG -> DL mất khoảng 7 tiếng
// 				arrivalTime := departureTime.Add(7 * time.Hour)

// 				newTrip := models.Trip{
// 					ID:                  primitive.NewObjectID(),
// 					CompanyID:           companyID,
// 					VehicleID:           vehicleID,
// 					Route:               route,
// 					DepartureTime:       departureTime,
// 					ExpectedArrivalTime: arrivalTime,
// 					Price:               float64(250000 + rand.Intn(150000)), // Giá từ 250k -> 400k
// 					Seats:               generateSeats(vehicle.TotalSeats),
// 				}
// 				tripsToInsert = append(tripsToInsert, newTrip)
// 			}
// 		}
// 	}

// 	_, err := tripCollection.InsertMany(ctx, tripsToInsert)
// 	if err != nil {
// 		log.Fatalf("Lỗi khi seed trips: %v", err)
// 	}

// 	log.Printf("Đã tạo %d trips.", len(tripsToInsert))
// }

// // generateSeats tạo danh sách ghế cho một chuyến đi
// func generateSeats(totalSeats int) []models.Seat {
// 	seats := make([]models.Seat, totalSeats)
// 	for i := 0; i < totalSeats; i++ {
// 		seatNumber := fmt.Sprintf("A%02d", i+1)
// 		status := "available"
// 		// Fake một vài ghế đã được đặt trước
// 		if rand.Float32() < 0.15 { // 15% cơ hội ghế đã được đặt
// 			status = "booked"
// 		}
// 		seats[i] = models.Seat{SeatNumber: seatNumber, Status: status}
// 	}
// 	return seats
// }

// // Hàm tiện ích để lấy keys từ map
// func keys[K comparable, V any](m map[K]V) []K {
// 	keys := make([]K, 0, len(m))
// 	for k := range m {
// 		keys = append(keys, k)
// 	}
// 	return keys
// }