package handler

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/MultivendorEcom/serviceutil/logger"
	"github.com/MultivendorEcom/storage"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

type DeliveryChargeTempData struct {
	CSRFField    template.HTML
	Form         DeliveryChargeForm
	FormErrors   map[string]string
	Data         []DeliveryChargeForm
	DistrictData []DistrictForm
	CountryData  []CountryForm
	StationData  []StationForm
}

func (s DeliveryChargeForm) Validate(srv *Server, id string) error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.DeliveryChargeStatus,
			validation.Required.Error("The status is required"),
			validation.Min(1).Error("Status is Invalid"),
			validation.Max(2).Error("Status is Invalid"),
		),
		validation.Field(&s.DistrictID,
			validation.Required.Error("The District name is required"),
			validation.By(checkDistrictExists(srv, s.DistrictID)),
		),
		validation.Field(&s.CountryID,
			validation.Required.Error("The country name is required"),
			validation.By(checkCountryExists(srv, s.CountryID)),
		),
		validation.Field(&s.WeightMin,
			validation.Required.Error("The minimum weight is required"),
		),
		validation.Field(&s.WeightMax,
			validation.Required.Error("The maximum weight is required"),
		),
		validation.Field(&s.DeliveryCharge,
			validation.Required.Error("The delivery charge is required"),
		),
		validation.Field(&s.StationID,
			validation.Required.Error("The station name is required"),
			validation.By(checkStationExists(srv, s.StationID)),
		),
	)
}

func (s *Server) deliveryChargeListHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("deliveryCharge list")
	tmpl := s.lookupTemplate("delivery-charge.html")
	if tmpl == nil {
		logger.Error(ult)
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}

	dcdata := s.deliveryChargeList(r, w)
	data := DeliveryChargeTempData{
		Data: dcdata,
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Error(ewte + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
}

func (s *Server) deliveryChargeFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("deliveryCharge submit")
	disdata := s.districtList(r, w, true)
	cntrydata := s.countryList(r, w, true)
	stndata := s.stationList(r, w, true)
	data := DeliveryChargeTempData{
		CSRFField:    csrf.TemplateField(r),
		DistrictData: disdata,
		CountryData:  cntrydata,
		StationData:  stndata,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) submitDeliveryChargeHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("deliveryCharge submit")
	uid, _ := s.GetSetSessionValue(r, w)
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	var form DeliveryChargeForm
	err := s.decoder.Decode(&form, r.PostForm)
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	if err := form.Validate(s, ""); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErrs[key] = value.Error()
				}
			}
		}
		data := DeliveryChargeTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	_, err = s.st.CreateDeliveryCharge(r.Context(), storage.DeliveryCharge{
		CountryID:      form.CountryID,
		DistrictID:     form.DistrictID,
		StationID:      form.StationID,
		Status:         form.DeliveryChargeStatus,
		WeightMin:      form.WeightMin,
		WeightMax:      form.WeightMax,
		DeliveryCharge: form.DeliveryCharge,
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedBy: uid,
			UpdatedBy: uid,
		},
	})
	if err != nil {
		logger.Error(err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) updateDeliveryChargeFormHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view update deliveryCharge form")
	params := mux.Vars(r)
	id := params["id"]
	disdata := s.districtList(r, w, true)
	cntrydata := s.countryList(r, w, true)
	stndata := s.stationList(r, w, true)
	brnFrm := s.getDeliveryChargeInfo(r, id, w)
	data := DeliveryChargeTempData{
		Form:         brnFrm,
		DistrictData: disdata,
		CountryData:  cntrydata,
		StationData:  stndata,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateDeliveryChargeHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("update deliveryCharge")
	uid, _ := s.GetSetSessionValue(r, w)
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form DeliveryChargeForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logger.Error(deErr + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusInternalServerError)
		return
	}
	if err := form.Validate(s, id); err != nil {
		vErrs := map[string]string{}
		if e, ok := err.(validation.Errors); ok {
			if len(e) > 0 {
				for key, value := range e {
					vErrs[key] = value.Error()
				}
			}
		}
		data := DeliveryChargeTempData{
			CSRFField:  csrf.TemplateField(r),
			FormErrors: vErrs,
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	dbdata := storage.DeliveryCharge{
		ID:             id,
		CountryID:      form.CountryID,
		DistrictID:     form.DistrictID,
		StationID:      form.StationID,
		Status:         form.DeliveryChargeStatus,
		WeightMin:      form.WeightMin,
		WeightMax:      form.WeightMax,
		DeliveryCharge: form.DeliveryCharge,
		CRUDTimeDate: storage.CRUDTimeDate{
			UpdatedBy: uid,
		},
	}
	_, err := s.st.UpdateDeliveryCharge(r.Context(), dbdata)
	if err != nil {
		logger.Error("error while update deliveryCharge data ." + err.Error())
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) viewDeliveryChargeHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("view deliveryCharge form")
	params := mux.Vars(r)
	id := params["id"]
	brnFrm := s.getDeliveryChargeInfo(r, id, w)
	data := DeliveryChargeTempData{
		Form: brnFrm,
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) updateDeliveryChargeStatusHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Update delivery Charge status")
	uid, _ := s.GetSetSessionValue(r, w)
	params := mux.Vars(r)
	id := params["id"]
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg + err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	res, err := s.st.GetDeliveryChargeBy(r.Context(), id)
	if err != nil {
		logger.Error("unable to get delivery Charge info " + err.Error())
	}
	if res.Status == 1 {
		_, err := s.st.UpdateDeliveryChargeStatus(r.Context(), storage.DeliveryCharge{
			ID:     id,
			Status: 2,
			CRUDTimeDate: storage.CRUDTimeDate{
				UpdatedBy: uid,
			},
		})
		if err != nil {
			logger.Error("unable to update status" + err.Error())
		}
	} else {
		_, err := s.st.UpdateDeliveryChargeStatus(r.Context(), storage.DeliveryCharge{
			ID:     id,
			Status: 1,
			CRUDTimeDate: storage.CRUDTimeDate{
				UpdatedBy: uid,
			},
		})
		if err != nil {
			logger.Error("unable to update status" + err.Error())
		}
	}
	http.Redirect(w, r, "/admin/"+deliveryChargeListPath, http.StatusSeeOther)
}

func (s *Server) deleteDeliveryChargeHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("delete deliveryCharge")
	if err := r.ParseForm(); err != nil {
		logger.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id := params["id"]
	err := s.st.DeleteDeliveryCharge(r.Context(), id, "1")
	if err != nil {
		logger.Error("error while delete deliveryCharge" + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	json.NewEncoder(w).Encode(msg)
}

func (s *Server) getDeliveryChargeInfo(r *http.Request, id string, w http.ResponseWriter) DeliveryChargeForm {
	res, err := s.st.GetDeliveryChargeBy(r.Context(), id)
	if err != nil {
		logger.Error("error while get deliveryCharge " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	Form := DeliveryChargeForm{
		ID:                   id,
		CountryID:            res.CountryID,
		CountryName:          res.CountryName.String,
		DistrictID:           res.DistrictID,
		DistrictName:         res.DistrictName.String,
		StationID:            res.StationID,
		StationName:          res.StationName.String,
		DeliveryChargeStatus: res.Status,
		WeightMin:            res.WeightMin,
		WeightMax:            res.WeightMax,
		DeliveryCharge:       res.DeliveryCharge,
		CreatedAt:            res.CreatedAt,
		CreatedBy:            res.CreatedBy,
		UpdatedAt:            res.UpdatedAt,
		UpdatedBy:            res.UpdatedBy,
	}
	return Form
}

func (s *Server) deliveryChargeList(r *http.Request, w http.ResponseWriter) []DeliveryChargeForm {
	dcList, err := s.st.GetDeliveryChargeList(r.Context())
	if err != nil {
		logger.Error("error while get delivery charge : " + err.Error())
		http.Redirect(w, r, ErrorPath, http.StatusSeeOther)
	}
	dcdata := make([]DeliveryChargeForm, 0)
	for _, item := range dcList {
		dcApnd := DeliveryChargeForm{
			ID:             item.ID,
			CountryID:      item.CountryID,
			CountryName:    item.CountryName.String,
			DistrictID:     item.DistrictID,
			DistrictName:   item.DistrictName.String,
			StationID:      item.StationID,
			StationName:    item.StationName.String,
			WeightMin:      item.WeightMin,
			WeightMax:      item.WeightMax,
			DeliveryCharge: item.DeliveryCharge,
			CreatedAt:      item.CreatedAt,
			CreatedBy:      item.CreatedBy,
			UpdatedAt:      item.UpdatedAt,
			UpdatedBy:      item.UpdatedBy,
		}
		dcdata = append(dcdata, dcApnd)
	}
	return dcdata
}
