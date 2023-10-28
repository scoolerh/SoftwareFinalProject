$(function(){ 
    for (var i = 0; i < 30; i++){
        create(i);
    }

    //create some confetti
    function create(i){
        var width = Math.random() * 15;
        var height = width * 0.4;
        var colorIndex = Math.ceil(Math.random() * 3);
        var color = "red";
        switch(colorIndex){
        case 1:
            color = "yellow";
            break;
        case 2:
            color = "lightblue";
            break;
        case 3:
            color = "pink";
            break
        case 4:
            color = "purple";
            break
        default:
            color = "red";
        }
        
        $('<div class="confetti-'+i+' '+color+'"></div>').css({
        "width" : width+"px",
        "height" : height+"px",
        "top" : -Math.random()*20+"%",
        "left" : Math.random()*100+"%",
        "opacity" : Math.random()+0.5,
        "transform" : "rotate("+Math.random()*360+"deg)"
        }).appendTo('.wrapper');  
        drop(i);
    }
    
    //drop that fetti
    function drop(x) {
        $('.confetti-'+x).animate({
        top: "100%",
        left: "+="+Math.random()*15+"%"
        }, Math.random()*2000 + 2000, function() {
        reset(x);
        });
    }
    
    //abort mission
    function reset(x) {
        $('.confetti-'+x).css('opacity','1');
        $('.confetti-'+x).animate({
        "top" : -Math.random()*20+"%",
        "left" : "-="+Math.random()*15+"%"
        }, 0, function() {
        drop(x);             
        });
    }
});