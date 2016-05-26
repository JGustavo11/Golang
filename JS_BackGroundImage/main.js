var backgrounds = new Array();
function uploadArray(){
	for(var i=0; i<=5; i++){
		backgrounds.push('images/' + i + '.jpg');
	}
	function preload(){
		var images = new Array();
		for(i=0; i< backgrounds.length; i++){
			images[i] = new Image();
			images[i].src = backgrounds[i];
		}
	}
	preload();
}

uploadArray();

var intervalID = window.setInterval(myCallback, 4000);

function myCallback(){
	var node = document.querySelector('html');
	var backgroundIndex = Math.floor(Math.random() * backgrounds.length);
	node.style.transition = 'all 2s deez nnnn';
	node.style.backgroundImage = 'url(' + backgrounds[backgroundIndex] + ')';
	backgrounds.splice(backgroundIndex, 1);
	if(backgrounds.length === 0){
		uploadArray();
	}
}
