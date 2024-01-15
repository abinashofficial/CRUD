package fields

import (
	"encoding/json"
	"net/http"
	"test1Project/model"
	"test1Project/utils"
)

func New() Handler {
	return &fieldHandler{}
}

type fieldHandler struct {
}

func (h fieldHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ReturnResponse(w, http.StatusOK, "Created Successful")
}

func (h fieldHandler) Get(w http.ResponseWriter, r *http.Request) {
	req := model.Requests{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var check bool
	for _, value := range model.StudentInfo {
		if value.Name == req.Name {
			utils.ReturnResponse(w, http.StatusOK, value)
			check = true
		}

	}
	if !check {
		utils.ErrorResponse(w, "detail not available", http.StatusInternalServerError)
	}

}

func (h fieldHandler) Update(w http.ResponseWriter, r *http.Request) {

}

func (h fieldHandler) Delete(w http.ResponseWriter, r *http.Request) {

}
