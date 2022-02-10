// view : DeliveryCompany Form with Dropdown
function viewDeliveryCompanyForm() {
    $.ajax({
        url: '/admin/delivery-company/create',
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            dcDropdown(obj)
        }
    });
}


// dropdown
function dcDropdown(obj) {
    var districts = obj.DistrictData
    var dis = $("#districtdd");
    $(districts).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.DistrictID);
        dis.append(option);
    });
    var cntrys = obj.CountryData
    var cntry = $("#countrydd");
    $(cntrys).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.CountryID);
        cntry.append(option);
    });
    var stations = obj.StationData
    var stn = $("#stationdd");
    $(stations).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.StationID);
        stn.append(option);
    });
}

// reset dropdown data
function resetData() {
    // reset all form data after close modal
    $('#modal-add').on('hidden.bs.modal', function () {
        $(this).find('form').trigger('reset');
        $("#countrydd").empty();
        $("#districtdd").empty();
        $("#stationdd").empty();
    });
}
// create
$(document).ready(function () {
    $('#saveForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        $.ajax({
            url: "/admin/delivery-company/create",
            method: 'post',
            data: $('form.tagForm').serialize(),
            success: function (data) {
                var obj = jQuery.parseJSON(data);
                var Toast = Swal.mixin({
                    toast: true,
                    position: 'top-end',
                    showConfirmButton: false,
                    timer: 3000
                });
                if (obj.Status) {
                    $('#saveForm').trigger("reset");
                    Toast.fire({
                        icon: 'success',
                        title: obj.Message
                    })
                    setTimeout(function () { window.location.reload(true); }, 1000);
                    $('#modal-add').modal('hide');
                } else {
                    $("#CompanyName").empty().append(obj.FormErrors.Name);
                    $("#CountryID").empty().append(obj.FormErrors.CountryID);
                    $("#DistrictID").empty().append(obj.FormErrors.DistrictID);
                    $("#StationID").empty().append(obj.FormErrors.StationID);
                    $("#Phone").empty().append(obj.FormErrors.Phone);
                    $("#Email").empty().append(obj.FormErrors.Email);
                    $("#CompanyAddress").empty().append(obj.FormErrors.DeliveryCompanyAddress);
                    $("#Position").empty().append(obj.FormErrors.Position);
                    $("#CompanyStatus").empty().append(obj.FormErrors.CompanyStatus);
                    Toast.fire({
                        icon: 'error',
                        title: "Please Insert All Data Carefully."
                    })
                }
            },
        });
        // reset all form data after close modal
        $('#modal-add').on('hidden.bs.modal', function () {
            $(this).find('form').trigger('reset');
            $("#CompanyName").empty();
            $("#CountryID").empty();
            $("#DistrictID").empty();
            $("#StationID").empty();
            $("#Phone").empty();
            $("#Email").empty();
            $("#CompanyAddress").empty();
            $("#Position").empty();
            $("#CompanyStatus").empty();
        });
    });
});

// view : DeliveryCompany
function viewDeliveryCompany(id) {
    $.ajax({
        url: "/admin/delivery-company/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#VCompanyName").empty().append(obj.Form.Name);
            $("#VCountryName").empty().append(obj.Form.CountryName);
            $("#VDistrictName").empty().append(obj.Form.DistrictName);
            $("#VStationName").empty().append(obj.Form.StationName);
            $("#VPhone").empty().append(obj.Form.Phone);
            $("#VEmail").empty().append(obj.Form.Email);
            $("#VCompanyAddress").empty().append(obj.Form.CompanyAddress);
            $("#VPosition").empty().append(obj.Form.Position);
            $("#VCompanyStatus").empty().append(obj.Form.CompanyStatus);
        }
    })
}

// Update : DeliveryCompany View
function viewDeliveryCompanyUpdateData(id) {
    $.ajax({
        url: "/admin/delivery-company/update/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#UdID").empty().val(obj.Form.ID);
            $("#UdCompanyName").empty().val(obj.Form.CompanyName);
            $("#UdPhone").empty().val(obj.Form.Phone);
            $("#UdEmail").empty().val(obj.Form.Email);
            $("#UdCompanyAddress").empty().val(obj.Form.CompanyAddress);
            $("#UdPosition").empty().val(obj.Form.Position);
            $("#UdCompanyStatus").empty().val(obj.Form.CompanyStatus);
            dcDropdownUpdate(obj)
        }
    });
    resetDataUpdate()
}

// dropdown Update
function dcDropdownUpdate(obj) {
    var districts = obj.DistrictData
    var dis = $("#districtdd-update");
    $(districts).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.DistrictID);
        dis.append(option);
    });
    var cntrys = obj.CountryData
    var cntry = $("#countrydd-update");
    $(cntrys).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.CountryID);
        cntry.append(option);
    });
    var stations = obj.StationData
    var stn = $("#stationdd-update");
    $(stations).each(function () {
        var option = $("<option />");
        option.html(this.Name);
        option.val(this.StationID);
        stn.append(option);
    });
}
// Update : DeliveryCompany Submit
$(document).ready(function () {
    $('#updateForm').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#UdID').val();
        $.ajax({
            url: "/admin/delivery-company/update/" + id,
            method: 'post',
            data: $('form.tagUpForm').serialize(),
            success: function (data) {
                var obj = jQuery.parseJSON(data);
                var Toast = Swal.mixin({
                    toast: true,
                    position: 'top-end',
                    showConfirmButton: false,
                    timer: 3000
                });
                if (obj.Status) {
                    $('updateForm').trigger("reset");
                    Toast.fire({
                        icon: 'success',
                        title: obj.Message
                    })
                    setTimeout(function () { window.location.reload(true); }, 1000);
                    $('#modal-update').modal('hide');
                } else {
                    $("#NameErr").empty().append(obj.FormErrors.Name);
                    $("#CountryIDErr").empty().append(obj.FormErrors.CountryID);
                    $("#DistrictIDErr").empty().append(obj.FormErrors.DistrictID);
                    $("#StationIDErr").empty().append(obj.FormErrors.StationID);
                    $("#Phone1Err").empty().append(obj.FormErrors.Phone);
                    $("#EmailErr").empty().append(obj.FormErrors.Email);
                    $("#AddressErr").empty().append(obj.FormErrors.CompanyAddress);
                    $("#PositionErr").empty().append(obj.FormErrors.Position);
                    $("#StatusErr").empty().append(obj.FormErrors.Status);
                    Toast.fire({
                        icon: 'error',
                        title: "Please Insert All Data Carefully."
                    })
                }
            }
        });
        // reset all form data after close modal
        $('#modal-update').on('hidden.bs.modal', function () {
            $(this).find('form').trigger('reset');
            $("#NameErr").empty();
            $("#CountryIDErr").empty();
            $("#DistrictIDErr").empty();
            $("#StationIDErr").empty();
            $("#Phone1Err").empty();
            $("#Phone2Err").empty();
            $("#EmailErr").empty();
            $("#AddressErr").empty();
            $("#PositionErr").empty();
            $("#StatusErr").empty();
        });
    });
});
// Delete : Country View
function deleteDeliveryCompanyData(id) {
    $.ajax({
        url: "/admin/delivery-company/view/" + id,
        method: 'get',
        success: function (data) {
            var obj = jQuery.parseJSON(data);
            $("#dID").empty().val(obj.Form.ID);
            $("#dCompanyName").empty().append(obj.Form.Name);
            $("#dCountryName").empty().append(obj.Form.CountryName);
            $("#dDistrictName").empty().append(obj.Form.DistrictName);
            $("#dStationName").empty().append(obj.Form.StationName);
            $("#dPosition").empty().append(obj.Form.Position);
            $("#dStatus").empty().append(obj.Form.Status);
        }
    })
}
// delete country
$(document).ready(function () {
    $('#deleteDeliveryCompany').submit(function (e) {
        e.preventDefault();
        $.ajaxSetup({
            headers: { 'X-CSRF-TOKEN': $('meta[name="_token"]').attr('content') }
        });
        var id = jQuery('#dID').val();
        $.ajax({
            url: "/admin/delivery-company/delete/" + id,
            method: 'get',
            success: function (data) {
                var obj = jQuery.parseJSON(data);
                var Toast = Swal.mixin({
                    toast: true,
                    position: 'top-end',
                    showConfirmButton: false,
                    timer: 3000
                });
                if (obj.Status) {
                    $('deleteCountry').trigger("reset");
                    Toast.fire({
                        icon: 'success',
                        title: obj.Message
                    })
                    setTimeout(function () { window.location.reload(true); }, 1000);
                    $('#modal-delete').modal('hide');
                }
            }
        });
    });
});

// reset dropdown data update
function resetDataUpdate() {
    // reset all form data after close modal
    $('#modal-update').on('hidden.bs.modal', function () {
        $(this).find('form').trigger('reset');
        $("#countrydd-update").empty();
        $("#districtdd-update").empty();
        $("#stationdd-update").empty();
    });
}