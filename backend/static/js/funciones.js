var servidor = "http://localhost:8080/api/login";
var indexUrl = "http://localhost:8080/backend/index";

function ingresar() {
    var _correo = $("#usuario").val();
    var _clave = $("#clave").val();
    console.log(_correo + _clave);
    $.ajax({
        type: "post",
        url: servidor,
        data: {
            username: _correo,
            password: _clave,
        }

    })
        .then(function (responseData) {
            // console.info("response", responseData.data);
            if (responseData.code !== 0) {
                alert(responseData.message);
            } else {
                if (typeof(Storage) !== "undefined") {
                    // Code for localStorage/sessionStorage.
                    localStorage.setItem('tokenUTNFRA', responseData.data.token);
                    window.location.href = indexUrl + "?uid=" + responseData.data.uid + "&token=" + responseData.data.token;
                } else {
                    console.log("Sorry! No Web Storage support..");
                }
            }

        }, function (error) {
            alert(error.responseText);
            console.info("error", error);
        });


}

function enviarToken() {
    $.ajax({
        url: servidor + "tomarToken/",
        type: 'GET',

        headers: {"miTokenUTNFRA": localStorage.getItem('tokenUTNFRA')}
    }).then(function (itemResponse) {

            console.info("bien -->", itemResponse);
        },
        function (error) {

            alert(error.responseText);
            console.info("error", error);
        });
}