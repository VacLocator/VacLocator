
function loadData(){
    //

    let latVal = document.getElementById('lat').value;
    let lngVal = document.getElementById('lng').value;
    let filterVal = document.getElementById('filter').checked ? '1' : '0';

    console.log('latVal', latVal);
    console.log('lngVal', lngVal);
    console.log('filterVal', filterVal);

    fetch('http://190.117.53.166:9594/lat/'+ latVal + '/lng/' + lngVal + '/filtro/' + filterVal)
    .then((resp) => resp.json())
    .then(data => {

        var msg = '<table><thead><tr><th> Nombre Centro </th>';
        msg+='<th> Latitud </th><th> Longitud </th><th> Cantidad Personas </th></tr></thead>';
        msg+='<tbody>';


        for(var i = 0; i < data.length; i++){
            let centro = data[i];

            msg+='<tr>';
            msg+='<td>'+centro.nombre_centro+'</td>';
            msg+='<td>'+centro.latitud+'</td>';
            msg+='<td>'+centro.longitud+'</td>';
            msg+='<td>'+centro.cantidad+'</td>';
            msg+= "</tr>";
        }
        msg+= "</tbody></table>";
        console.log(msg);

        var masCercano = data[0];
        //TODO
        //mas cercano es el que quiero mostrar en el front :u
        console.log('masCercano : ', masCercano);


      //  alert('Se encontraron los centros de salud más cercanos : ' + msg);
      Swal.fire({
        html: msg,
        width: 900,
        padding: '3em',
        background: '#fff url(/images/trees.png)',
        backdrop: `
          rgba(0,0,123,0.4)
          url("nyan-cat.gif")
          left top
          no-repeat
        `
      })


    })
    .catch(error => {
        alert('Ocurrió un error al realizar el calculo');
    });

}
