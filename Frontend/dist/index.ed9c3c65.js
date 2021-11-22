//la variable centers esta en el data.js que linkeamos en el index
let map;
let markers = [];
function handleLocationError(browserHasGeolocation, infoWindow, pos) {
    infoWindow.setPosition(pos);
    infoWindow.setContent(browserHasGeolocation ? "Error: The Geolocation service failed." : "Error: Your browser doesn't support geolocation.");
    infoWindow.open(map);
}
const setListener = ()=>{
    document.querySelectorAll(".center_individualNames").forEach((centerName, index)=>{
        centerName.addEventListener("click", ()=>{
            google.maps.event.trigger(markers[index], "click");
        });
    });
};
const showCentersList = ()=>{
    let centerHTML = "";
    centers.forEach((centro)=>{
        centerHTML += `<h4 class="center_individualNames">${centro.name}</h4>` //para  pasar variable usar ``
        ;
    });
    document.getElementById("centers_names").innerHTML = centerHTML;
};
const createMarker = (coord, name)=>{
    const html = `<h4>${name}</h4>`;
    const marker = new google.maps.Marker({
        position: coord,
        map: map,
        icon: "./icons/vacuna.png"
    });
    google.maps.event.addListener(marker, "click", ()=>{
        infoWindow.setContent(html);
        infoWindow.open(map, marker);
    });
    markers.push(marker);
};
const createLocationMarkers = ()=>{
    centers.forEach((centro)=>{
        let coord = new google.maps.LatLng(centro.lat, centro.lng);
        let name = centro.name;
        createMarker(coord, name);
    });
};
function initMap() {
    let lugar = {
        lat: -12.0706976,
        lng: -77.0454013
    };
    map = new google.maps.Map(document.getElementById("map"), {
        center: lugar,
        zoom: 14
    });
    //var val = document.getElementById("filter").checked;
    //console.log(val);
    //const marker = new google.maps.Marker({
    //    position: lugar,
    //    map:map
    //})
    createLocationMarkers();
    //const marker = new google.maps.Marker({
    //    position: lugar,
    //    map:map,
    //})
    var posi;
    var filtro;
    var valores = {
        lat,
        lng,
        filtro
    };
    infoWindow = new google.maps.InfoWindow();
    const locationButton = document.createElement("button");
    locationButton.textContent = "Click para geolocalizacion";
    locationButton.classList.add("custom-map-control-button");
    map.controls[google.maps.ControlPosition.TOP_CENTER].push(locationButton);
    locationButton.addEventListener("click", ()=>{
        // Try HTML5 geolocation.
        if (navigator.geolocation) navigator.geolocation.getCurrentPosition((position)=>{
            const pos = {
                lat: position.coords.latitude,
                lng: position.coords.longitude
            };
            posi = pos;
            valores.lat = pos.lat;
            valores.lng = pos.lng;
            document.getElementById("lat").value = valores.lat;
            document.getElementById("lng").value = valores.lng;
            //valores.filtro = document.getElementById("filter").value;
            infoWindow.setPosition(pos);
            infoWindow.setContent("Ubicacion encontrada.");
            infoWindow.open(map);
            map.setCenter(pos);
        }, ()=>{
            handleLocationError(true, infoWindow, map.getCenter());
        });
        else // Browser doesn't support Geolocation
        handleLocationError(false, infoWindow, map.getCenter());
    });
    //let html = '<h3>Centro de la Ciudad</h3>'
    //google.maps.event.addListener(marker, "click", ()=>{
    //    infoWindow.setContent(html);
    //    infoWindow.open(map,marker)
    //})
    showCentersList();
    setListener();
}

//# sourceMappingURL=index.ed9c3c65.js.map
