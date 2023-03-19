package ratings

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	localMdl "semesta-ban/internal/api/middleware"
	"semesta-ban/internal/api/products"
	"semesta-ban/internal/api/response"
	"semesta-ban/pkg/constants"
	cn "semesta-ban/pkg/constants"
	"semesta-ban/pkg/crashy"
	"semesta-ban/pkg/helper"
	"semesta-ban/pkg/log"
	"semesta-ban/repository/repo_products"
	rateRepo "semesta-ban/repository/repo_ratings"

	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
)

type RatingsHandler struct {
	db             *sqlx.DB
	rateRepository rateRepo.RatingsRepository
	prodRepo       repo_products.ProductsRepository
	baseAssetUrl   string
	uploadPath     string
	imgMaxSize     int
}

//todo REMEMBER 30 May gmail tidak support lagi less secure app find solution

func NewRatingsHandler(db *sqlx.DB, rr rateRepo.RatingsRepository, pr repo_products.ProductsRepository, baseAssetUrl, uploadPath string,
	imgMaxSize int) *RatingsHandler {
	return &RatingsHandler{db: db, rateRepository: rr, prodRepo: pr, baseAssetUrl: baseAssetUrl, uploadPath: uploadPath, imgMaxSize: imgMaxSize}
}

func (rh *RatingsHandler) SubmitRatingProduct(w http.ResponseWriter, r *http.Request) {
	var (
		ctx          = r.Context()
		authData     = ctx.Value(localMdl.CtxKey).(localMdl.Token)
		prodId       = r.FormValue("product_id")
		comment      = r.FormValue("comment")
		rate         = r.FormValue("rate")
		fileNameList = []string{}
	)

	//validate input
	if len(rate) == 0 || len(prodId) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, crashy.Message(crashy.ErrCodeValidation)), http.StatusBadRequest)
		return
	}
	// return

	custId, errCode, err := rh.prodRepo.GetCustomerId(ctx, authData.Uid)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	//process image if exists
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		// handle error
		response.Nay(w, r, crashy.New(err, crashy.ErrFileNotFound, crashy.Message(crashy.ErrCode(crashy.ErrFileNotFound))), http.StatusBadRequest)
		return
	}
	//check all file size before uploading
	for _, fh := range r.MultipartForm.File["photos"] {

		if fh.Size > int64(helper.ConvertFileSizeToMb(rh.imgMaxSize)) {
			errMsg := fmt.Sprintf("%s%v mb", crashy.Message(crashy.ErrCode(crashy.ErrExceededFileSize)), rh.imgMaxSize)
			response.Nay(w, r, crashy.New(errors.New(crashy.ErrExceededFileSize), crashy.ErrExceededFileSize, errMsg), http.StatusBadRequest)
			return
		}

	}
	for _, fh := range r.MultipartForm.File["photos"] {
		f, err := fh.Open()
		if err != nil {
			// Handle error
			response.Nay(w, r, crashy.New(err, crashy.ErrFileNotFound, crashy.Message(crashy.ErrCode(crashy.ErrFileNotFound))), http.StatusBadRequest)
			return
		}

		tempFile, err := ioutil.TempFile(rh.uploadPath+cn.RatingsProductDir, "pic-*.png")
		if err != nil {
			response.Nay(w, r, crashy.New(err, crashy.ErrUploadFile, crashy.Message(crashy.ErrCode(crashy.ErrUploadFile))), http.StatusBadRequest)
			return
		}
		defer tempFile.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(f)
		if err != nil {
			response.Nay(w, r, crashy.New(err, crashy.ErrUploadFile, crashy.Message(crashy.ErrCode(crashy.ErrUploadFile))), http.StatusBadRequest)
			return
		}
		// write this byte array to our temporary file
		fileName := helper.GetUploadedFileName(tempFile.Name())

		tempFile.Write(fileBytes)
		tempFile.Chmod(0604)
		log.Infof("success upload %s to the server x \n", fileName)
		fileNameList = append(fileNameList, fileName)

		// Read data from f
		f.Close()
	}

	//submit data to db
	errCode, err = rh.rateRepository.SubmitRatingProduct(ctx, custId, prodId, comment, rate, fileNameList)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)

}

func (rh *RatingsHandler) SubmitRatingOutlet(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		authData = ctx.Value(localMdl.CtxKey).(localMdl.Token)
		outletId = r.FormValue("outlet_id")
		comment  = r.FormValue("comment")
		rate     = r.FormValue("rate")
		invoiceId         = r.FormValue("invoice_id")
		fileNameList = []string{}
	)

	//validate input
	if len(rate) == 0 || len(outletId) == 0 {
		response.Nay(w, r, crashy.New(errors.New(crashy.ErrCodeValidation), crashy.ErrCodeValidation, crashy.Message(crashy.ErrCodeValidation)), http.StatusBadRequest)
		return
	}
	// return

	custId, errCode, err := rh.prodRepo.GetCustomerId(ctx, authData.Uid)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	//process image if exists
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		// handle error
		response.Nay(w, r, crashy.New(err, crashy.ErrFileNotFound, crashy.Message(crashy.ErrCode(crashy.ErrFileNotFound))), http.StatusBadRequest)
		return
	}
	//check all file size before uploading
	for _, fh := range r.MultipartForm.File["photos"] {

		if fh.Size > int64(helper.ConvertFileSizeToMb(rh.imgMaxSize)) {
			errMsg := fmt.Sprintf("%s%v mb", crashy.Message(crashy.ErrCode(crashy.ErrExceededFileSize)), rh.imgMaxSize)
			response.Nay(w, r, crashy.New(errors.New(crashy.ErrExceededFileSize), crashy.ErrExceededFileSize, errMsg), http.StatusBadRequest)
			return
		}

	}
	for _, fh := range r.MultipartForm.File["photos"] {
		f, err := fh.Open()
		if err != nil {
			// Handle error
			response.Nay(w, r, crashy.New(err, crashy.ErrFileNotFound, crashy.Message(crashy.ErrCode(crashy.ErrFileNotFound))), http.StatusBadRequest)
			return
		}

		tempFile, err := ioutil.TempFile(rh.uploadPath+cn.RatingOutletDir, "pic-*.png")
		if err != nil {
			response.Nay(w, r, crashy.New(err, crashy.ErrUploadFile, crashy.Message(crashy.ErrCode(crashy.ErrUploadFile))), http.StatusBadRequest)
			return
		}
		defer tempFile.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(f)
		if err != nil {
			response.Nay(w, r, crashy.New(err, crashy.ErrUploadFile, crashy.Message(crashy.ErrCode(crashy.ErrUploadFile))), http.StatusBadRequest)
			return
		}
		// write this byte array to our temporary file
		fileName := helper.GetUploadedFileName(tempFile.Name())

		tempFile.Write(fileBytes)
		tempFile.Chmod(0604)
		log.Infof("success upload %s to the server x \n", fileName)
		fileNameList = append(fileNameList, fileName)

		// Read data from f
		f.Close()
	}

	//submit data to db
	errCode, err = rh.rateRepository.SubmitRatingOutlet(ctx, custId, outletId, comment, rate, fileNameList, invoiceId)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}

	response.Yay(w, r, "success", http.StatusOK)

}

func (rh *RatingsHandler) GetListRatingOutler(w http.ResponseWriter, r *http.Request) {
	var (
		p              GetListRatingOutletRequest
		ctx            = r.Context()
		listRatingTemp = []ListRating{}
	)

	if err := render.Bind(r, &p); err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCodeValidation, err.Error()), http.StatusBadRequest)
		return
	}

	limit := p.Limit
	if limit < 1 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}
	page := p.Page
	if page < 1 {
		page = 1
	}

	ratingSummary, errCode, err := rh.rateRepository.GetRatingSummary(ctx, p.OutletId)
	if err != nil {
		response.Nay(w, r, crashy.New(err, crashy.ErrCode(errCode), crashy.Message(crashy.ErrCode(errCode))), http.StatusInternalServerError)
		return
	}
	avgRate := float32(ratingSummary.SumRating) / float32(ratingSummary.All)
	avgRateString := fmt.Sprintf("%.1f", avgRate)
	satisfyRate := avgRate / float32(5) * 100
	satisfyRateString := fmt.Sprintf("%.2f", satisfyRate)

	listRating, totalData, listIdOutlet, errCode, err := rh.rateRepository.GetListRatingOutlet(ctx, rateRepo.GetListRatingOutletRequestParam{
		Limit:       limit,
		Page:        page,
		WithMedia:   p.WithMedia,
		WithComment: p.WithComment,
		Ratings:     p.Ratings,
	}, p.OutletId)

	listImageRating, errCode, err := rh.rateRepository.GetListRatingImage(ctx, listIdOutlet)
	mappedListImageByRatingId := make(map[int][]ListRatingImage)
	for _, m := range listImageRating {
		mappedListImageByRatingId[m.IdRating] = append(mappedListImageByRatingId[m.IdRating], ListRatingImage{
			ImageUrl: rh.baseAssetUrl + constants.RatingOutletDir + m.ImageName,
		})
	}

	for _, v := range listRating {
		avatar := ""

		if len(v.CustomerAvatar.String) > 0 && v.CustomerAvatar.String[:3] == "pic" {
			avatar = rh.baseAssetUrl + constants.UserDir + v.CustomerAvatar.String
		} else if len(v.CustomerAvatar.String) > 0 && v.CustomerAvatar.String[:3] != "pic" {
			avatar = v.CustomerAvatar.String
		}
		if p.WithMedia {
			if len(mappedListImageByRatingId[v.IdRating]) > 0 {
				listRatingTemp = append(listRatingTemp, ListRating{
					IdRating:       v.IdRating,
					CustomerName:   v.CustomerName,
					CustomerAvatar: avatar,
					Rating:         v.Rating,
					OutletName:     v.OutletName,
					Comment:        v.Comment,
					CreatedAt:      v.CreatedAt,
					ListImages:     mappedListImageByRatingId[v.IdRating],
				},
				)
			}
		} else {
			listRatingTemp = append(listRatingTemp, ListRating{
				IdRating:       v.IdRating,
				CustomerName:   v.CustomerName,
				CustomerAvatar: avatar,
				Rating:         v.Rating,
				OutletName:     v.OutletName,
				Comment:        v.Comment,
				CreatedAt:      v.CreatedAt,
				ListImages:     mappedListImageByRatingId[v.IdRating],
			},
			)
		}
	}
	if p.WithMedia {
		totalData = len(mappedListImageByRatingId)
	}

	response.Yay(w, r, GetListRatingOutletResponse{
		DataInfo: products.DataInfo{
			CurrentPage: page,
			MaxPage: func() int {
				maxPage := float64(totalData) / float64(limit)
				if helper.IsFloatNoDecimal(maxPage) {
					return int(maxPage)
				}
				return int(maxPage) + 1
			}(),
			Limit:       limit,
			TotalRecord: totalData,
		},
		ListRatingData: listRatingTemp,

		SummaryRatingData: SummaryRating{
			AvgRating:    avgRateString,
			SatisfyLevel: satisfyRateString + "%",
			All:          ratingSummary.All,
			WithMedia:    ratingSummary.WithMedia,
			WithComment:  ratingSummary.WithComment,
			RateFive:     ratingSummary.RateFive,
			RateFour:     ratingSummary.RateFour,
			RateThree:    ratingSummary.RateThree,
			RateTwo:      ratingSummary.RateTwo,
			RateOne:      ratingSummary.RateOne,
		},
	}, http.StatusOK)
	return
}
