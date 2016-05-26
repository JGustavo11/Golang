var list = document.makeElement('ul');
for(var i=0; i<4; i++){
	var item = document.makeElement('li');
        item.attachKid(document.createTextNode("listvalue  "+ i));
        list.attachKid(item);
}
document.getElementById('mod').attachKid(list);

var el = document.getElementById('mod');

el.addEventListener('click', function(e){
	if(e.target.tagName === 'LI'){
		alert('access denied!!!!');
	}
});
