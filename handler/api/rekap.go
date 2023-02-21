package api

import (
	"net/http"
	"vandesar/entity"
	"vandesar/repository"
)

type RekapAPI struct {
	rekapRepo *repository.RekapRepository
}

func NewRekapAPI(rekapRepo *repository.RekapRepository) *RekapAPI {
	return &RekapAPI{
		rekapRepo:rekapRepo,
	}
}

func (p *RekapAPI) ListRekapPerMonth(w http.ResponseWriter, r *http.Request) {
	adminId := r.Context().Value("id").(uint)
	listRekap, err := p.rekapRepo.ListRekap(adminId)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	WriteJSON(w, http.StatusOK, listRekap)
}

func (p *RekapAPI) ListRekapPerDays(w http.ResponseWriter, r *http.Request) {
	adminId := r.Context().Value("id").(uint)
	listRekap, err := p.rekapRepo.ListRekapPerDays(adminId)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	WriteJSON(w, http.StatusOK, listRekap)
}

// func (p *RekapAPI) GetRekap(w http.ResponseWriter, r *http.Request) {
// 	adminId := r.Context().Value("id").(uint)

// 	product := r.URL.Query()
// 	date, foundDate := product["date"] // bulan-tahun, ex: 01-2021

// 	if foundDate {
// 		fileName := "rekap-" + strconv.Itoa(int(adminId)) + "-" + date[0]

// 		fmt.Println("---" + fileName)

// 		w.Header().Set("Content-Type", "application/pdf")
// 		w.Header().Set("Content-Disposition", "attachment; filename="+fileName+".pdf")
// 		http.ServeFile(w, r, "rekap-pdf/"+fileName+".pdf")
// 	}

// 	// if foundProductSearch {
// 	// 	ProductBySearch, err := p.productService.SearchProducts(r.Context(), productSearch[0])
// 	// 	if err != nil {
// 	// 		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
// 	// 		return
// 	// 	}

// 	// 	var productsFiltered []entity.Product
// 	// 	for _, v := range ProductBySearch {
// 	// 		if v.UserID == adminId {
// 	// 			productsFiltered = append(productsFiltered, v)
// 	// 		}
// 	// 	}

// 	// 	WriteJSON(w, http.StatusOK, productsFiltered)
// 	// 	return
// 	// }

// 	// list, err := p.productService.GetProducts(r.Context(), adminId)
// 	// if err != nil {
// 	// 	WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
// 	// 	return
// 	// }

// 	WriteJSON(w, http.StatusOK, nil)
// }
