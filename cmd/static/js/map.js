
//Initialise with default location
let selectedLocation = [-1.1680703527302385, 52.95592495232059]
let marker;
const latitudeInput = document.getElementById("latitude")
const longitudeInput = document.getElementById("longitude")

const addBusinessMarkerToMap = (coords)=>{
    if(!marker){
        const markerWrapper = document.createElement("div");
        markerWrapper.className = "marker-wrapper";
        const markerEl = document.createElement("div");
        markerEl.className = "marker business";
        markerWrapper.appendChild(markerEl);

        marker = new mapboxgl.Marker(markerWrapper).setLngLat(coords);
        marker.addTo(map)
    }else{
        marker.setLngLat(coords)
    }
    console.log(coords)
    longitudeInput.value = coords.lng
    latitudeInput.value = coords.lat
}



const handleMapClick = (e)=>{
    addBusinessMarkerToMap(e.lngLat) 
}

mapboxgl.accessToken = 'pk.eyJ1IjoiYXJkZWEtYXBwcyIsImEiOiJjbGhibm56MTAwdGdjM3JudTFlYThnZ3RtIn0.lFLUhfJCaqJUFXOk6CsERw';
const map = new mapboxgl.Map({
    container: 'location-picker', 
    center: selectedLocation, 
    zoom: 9,
    style: 'mapbox://styles/mapbox/satellite-streets-v12'
});
map.on('click', handleMapClick);

const postcodeInput = document.getElementById("postcode")
if(postcodeInput){
    postcodeInput.addEventListener("blur", (e)=>{
       const v = e.target.value;
       moveToPostcode(v)
    })
}

const locationCache = {}
async function moveToPostcode(postcode){
    if(!postcode) return; 
    if(locationCache[postcode]){
        if(selectedLocation != locationCache[postcode]){
            selectedLocation = locationCache[postcode]
            map.jumpTo({center: locationCache[postcode]})
        }
        return
    }

    const url = `/api/geocoding?postcode=${postcode}`
    const res = await fetch(url)
    if(!res.ok){
        if(res.status === 404){
            setErrorToast("Postcode not found")
        }
        return;
        }

    const location = await res.json()
    locationCache[postcode] = [location.lng, location.lat]
    map.jumpTo({center: locationCache[postcode], zoom: 12})
}