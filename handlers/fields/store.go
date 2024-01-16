package fields

import (
	"crud/tapcontext"
	"crud/utils"
	"encoding/json"
	"go.elastic.co/apm"
	"net/http"
)

func New() Handler {
	return &fieldHandler{}
}

type fieldHandler struct {
}

func (h fieldHandler) Create(w http.ResponseWriter, r *http.Request) {
	functionDesc := "Create Api"
	var req map[string]string
	ctx := tapcontext.UpgradeCtx(r.Context())
	span, _ := apm.StartSpan(ctx, functionDesc, "Handler")
	defer span.End()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError, err, nil)
		return
	}
	utils.ReturnResponse(w, http.StatusOK, "Created Successful")
}

func (h fieldHandler) Get(w http.ResponseWriter, r *http.Request) {
	//functionDesc := "Get Api"
	//
	//ctx := tapcontext.UpgradeCtx(r.Context())
	//span, _ := apm.StartSpan(ctx, functionDesc, "Handler")
	//defer span.End()
	//req := model.Requests{}
	//err := json.NewDecoder(r.Body).Decode(&req)
	//if err != nil {
	//	utils.ErrorResponse(ctx, w, err.Error(), http.StatusInternalServerError, err, nil)
	//	return
	//}
	//var check bool
	//for _, value := range model.StudentInfo {
	//	if value.Name == req.Name {
	//		utils.ReturnResponse(w, http.StatusOK, value)
	//		check = true
	//	}
	//
	//}
	//if !check {
	//	utils.ErrorResponse(ctx, w, "detail not available", http.StatusInternalServerError, err, nil)
	//
	//}

}

func (h fieldHandler) Update(w http.ResponseWriter, r *http.Request) {

}

func (h fieldHandler) Delete(w http.ResponseWriter, r *http.Request) {

}
