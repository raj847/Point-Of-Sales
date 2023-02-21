package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"vandesar/entity"
	"vandesar/repository"

	"github.com/Rhymond/go-money"
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/minio/minio-go/v7"
)

func DoRekapEachMonth(service *RekapService) {
	wib, _ := time.LoadLocation("Asia/Jakarta")
	s := gocron.NewScheduler(wib)

	// each last month at 22:00
    _, err := s.Every(1).MonthLastDay().At("22:00").Do(service.Rekap)
	// _, err := s.Every(5).Second().Do(service.Rekap) // simulate for testing purpose
    if err != nil {
        log.Fatalf("Gagal menjadwalkan tugas: %v", err)
        return
    }

	fmt.Println("CRON JOB STARTED...")
	s.StartAsync()
}

func DoRekapEveryDay(service *RekapService) {
	wib, _ := time.LoadLocation("Asia/Jakarta")
	s := gocron.NewScheduler(wib)

	// each last month at 22:00
    // _, err := s.Every(1).MonthLastDay().At("22:00").Do(service.Rekap)
	_, err := s.Every(1).Day().At("21:09").Do(service.RekapPerDay) // simulate for testing purpose
    if err != nil {
        log.Fatalf("Gagal menjadwalkan tugas: %v", err)
        return
    }

	fmt.Println("CRON JOB STARTED...")
	s.StartAsync()
}

type RekapService struct {
	rekapRepo *repository.RekapRepository
	transactionRepo *repository.TransactionRepository
	userRepo *repository.UserRepository

	minioClient *minio.Client
}

func NewRekapService(
	rekapRepo *repository.RekapRepository,
	transactRepo *repository.TransactionRepository,
	userRepo *repository.UserRepository,
	minioClient *minio.Client) *RekapService {

	return &RekapService{
		rekapRepo: rekapRepo,
		transactionRepo: transactRepo,
		userRepo: userRepo,
		minioClient: minioClient,
	}
}

type RekapFetcher struct {
	AdminId uint
	Rekap entity.Rekap
}

func (r *RekapService) Rekap() {
	now := time.Now()
	wib, _ := time.LoadLocation("Asia/Jakarta")

	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, wib)
	currentMonth := now.Month()

    // Mengambil tanggal akhir bulan
    tahun, _, _ := now.Date()
    bulanBerikutnya := currentMonth + 1
    if bulanBerikutnya > 12 {
        bulanBerikutnya = 1
        tahun++
    }

    endOfMonth := time.Date(tahun, bulanBerikutnya, 0, 0, 0, 0, 0, wib)

	admins, _ := r.userRepo.GetAllAdmins()
	// rekapFetcher := make([]RekapFetcher, 0, len(admins))

	for _,admin := range admins {
		tmpRekap := RekapFetcher{
			AdminId: admin.ID,
		}

		result, _ := r.transactionRepo.ReadTransByDateRange(startOfMonth, endOfMonth, admin.ID)

		totalPrice := 0.0
		totalProfit := 0.0
		totalDebt := 0.0
		totalPeopleDebt := 0

		for _,v := range result {
			totalPrice += v.TotalPrice
			totalProfit += v.TotalProfit
			totalDebt += v.Debt
			if v.Status == "Belum Lunas" {
				totalPeopleDebt++
			}
		}

		tmpRekap.Rekap = entity.Rekap{
			AdminID: admin.ID,
			TotalPrice: totalPrice,
			TotalProfit: totalProfit,
			TotalDebt: totalDebt,
			TotalPeopleDebt: totalPeopleDebt,
			StartDate: startOfMonth,
			EndDate: endOfMonth,
		}

		// rekapFetcher = append(rekapFetcher, tmpRekap)

		linkPdf, err := GenerateRekapPDF(
			tmpRekap.Rekap,
			r.minioClient,
		)
		if err != nil {
			fmt.Println("Gagal membuat PDF")
		}

		tmpRekap.Rekap.LinkPdf = linkPdf

		err = r.rekapRepo.AddRekap(tmpRekap.Rekap)
		if err != nil {
			fmt.Println("Gagal menyimpan rekap ke database, ", err)
		}
	}
}

func (r *RekapService) RekapPerDay() {
	now := time.Now()
	wib, _ := time.LoadLocation("Asia/Jakarta")

	// 06:00 WIB - WAKTU SEKARANG WIB
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, wib)
	endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, wib)

	admins, _ := r.userRepo.GetAllAdmins()
	// rekapFetcher := make([]RekapFetcher, 0, len(admins))

	for _,admin := range admins {
		tmpRekap := RekapFetcher{
			AdminId: admin.ID,
		}

		result, _ := r.transactionRepo.ReadTransByDateRange(startTime, endTime, admin.ID)

		totalPrice := 0.0
		totalProfit := 0.0
		totalDebt := 0.0
		totalPeopleDebt := 0

		for _,v := range result {
			totalPrice += v.TotalPrice
			totalProfit += v.TotalProfit
			totalDebt += v.Debt
			if v.Status == "Belum Lunas" {
				totalPeopleDebt++
			}
		}

		tmpRekap.Rekap = entity.Rekap{
			AdminID: admin.ID,
			TotalPrice: totalPrice,
			TotalProfit: totalProfit,
			TotalDebt: totalDebt,
			TotalPeopleDebt: totalPeopleDebt,
			StartDate: startTime,
			EndDate: endTime,
		}

		// rekapFetcher = append(rekapFetcher, tmpRekap)

		linkPdf, err := GenerateRekapPDF(
			tmpRekap.Rekap,
			r.minioClient,
		)
		if err != nil {
			fmt.Println("Gagal membuat PDF")
		}

		tmpRekap.Rekap.LinkPdf = linkPdf

		// conert into recap per day
		rekapPerDay := entity.RekapPerDay{
			AdminID: tmpRekap.Rekap.AdminID,
			TotalPrice: tmpRekap.Rekap.TotalPrice,
			TotalProfit: tmpRekap.Rekap.TotalProfit,
			TotalDebt: tmpRekap.Rekap.TotalDebt,
			TotalPeopleDebt: tmpRekap.Rekap.TotalPeopleDebt,
			StartDate: tmpRekap.Rekap.StartDate,
			EndDate: tmpRekap.Rekap.EndDate,
			LinkPdf: tmpRekap.Rekap.LinkPdf,
		}

		err = r.rekapRepo.AddRekapPerDay(rekapPerDay)
		if err != nil {
			fmt.Println("Gagal menyimpan rekap per day ke database, ", err)
		}
	}
}

const rekapFolderUploadPath = "./rekap-pdf"

func GenerateRekapPDF(rekap entity.Rekap, minioClient *minio.Client) (string, error ){
	fileName := fmt.Sprintf("rekap-%d-%d-%d.pdf", rekap.AdminID, rekap.StartDate.Month(), rekap.StartDate.Year())

	fmt.Println("filename yg akan di generate = ", fileName)

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Menambahkan header tabel
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(190, 10, fileName, "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 12)
    pdf.CellFormat(20, 10, "Admin ID", "1", 0, "", false, 0, "")
    pdf.CellFormat(50, 10, "Total Price", "1", 0, "", false, 0, "")
    pdf.CellFormat(50, 10, "Total Profit", "1", 0, "", false, 0, "")
    pdf.CellFormat(50, 10, "Total Debt", "1", 0, "", false, 0, "")
    pdf.CellFormat(40, 10, "Debt", "1", 0, "", false, 0, "")
    pdf.CellFormat(40, 10, "Start Date", "1", 0, "", false, 0, "")
    pdf.CellFormat(40, 10, "End Date", "1", 1, "", false, 0, "")

	startDate := fmt.Sprintf("%d-%d-%d", rekap.StartDate.Year(), rekap.StartDate.Month(), rekap.StartDate.Day())
	endDate := fmt.Sprintf("%d-%d-%d", rekap.EndDate.Year(), rekap.EndDate.Month(), rekap.EndDate.Day())

	totalPriceFormat := money.NewFromFloat(rekap.TotalPrice, money.IDR).Display()
	totalProfitFormat := money.NewFromFloat(rekap.TotalProfit, money.IDR).Display()
	totalDebtFormat := money.NewFromFloat(rekap.TotalDebt, money.IDR).Display()

	pdf.CellFormat(20, 10, fmt.Sprintf("%d", rekap.AdminID), "1", 0, "", false, 0, "")
	pdf.CellFormat(50, 10, totalPriceFormat, "1", 0, "", false, 0, "")
	pdf.CellFormat(50, 10, totalProfitFormat, "1", 0, "", false, 0, "")
	pdf.CellFormat(50, 10, totalDebtFormat, "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%d Orang", rekap.TotalPeopleDebt), "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 10, startDate, "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 10, endDate, "1", 1, "", false, 0, "")

	err := pdf.OutputFileAndClose(fmt.Sprintf("%s/%s", rekapFolderUploadPath, fileName))
	if err != nil {
		fmt.Println("error happen = " , err)
		return "", err
	}

	filePath := fmt.Sprintf("%s/%s", rekapFolderUploadPath, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	fileName = uuid.New().String() + "-" + fileName
	fmt.Println(fileName)
	fileLinkUploaded, err := UploadToCloud(context.Background(),  minioClient, file, fileName)
	if err != nil {
		log.Println("gagal upload ke cloud")
	}

	// remove from filePath
	err = os.Remove(filePath)
	if err != nil {
		log.Print("gagal remove file")
	}

	return fileLinkUploaded, nil
}
