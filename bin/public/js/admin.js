

var popUpToast;
(function() {
    var toastID=0;
    popUpToast=function(word, length) {
        toastID++;
        if (length==undefined) {
            length=3000;
        }
        var newThing=$("<div class='toast toastUndeclared' id='toast_"+toastID+"'>").text(word);

        $(".toastContainer").prepend(newThing);
        setTimeout(function(){newThing.removeClass("toastUndeclared");}, 100);
        setTimeout(function(){
            newThing.addClass("toastUndeclared");
            setTimeout(function(){newThing.remove();}, 300);
        }, length);
    }
})()

//////////////////////////////////

var mapping={};
function makeElement(k) {
    var kele=$("<textarea class='showAndEdit stext' placeholder='source'>");
    kele[0].value=k;
    var vele=$("<textarea class='showAndEdit vtext' placeholder='destination'>");
    vele[0].value=mapping[k];
    var thisk=k;
    kele.blur(function() {
        if (kele[0].value=="") {
            if (vele[0].value=="" && thisk!="") {
                modify(thisk, "", "");
            }
            return;
        }
        if (kele[0].value==thisk)
            return;
        if (kele[0].value in mapping) {
            popUpToast("Duplicate key detected: "+kele[0].value);
            return;
        }
        if (vele[0].value!="") {
            modify(thisk, kele[0].value, vele[0].value);
        }
    });
    vele.blur(function() {
        if (vele[0].value=="") {
            if (kele[0].value=="" && thisk!="") {
                modify(thisk, "", "");
            }
            return;
        }
        if (vele[0].value==mapping[thisk])
            return;
        if (kele[0].value!="") {
            modify(thisk, kele[0].value, vele[0].value);
        }
    });
    var tot=$("<div class='tcontainer' id='tentry_"+k+"'>").append(kele).append(vele);
    $("body").append(tot);
}
function mapSync() {
    $(".tcontainer").remove();
    mapping[""]="";
    for (var k in mapping) {
        makeElement(k);
    }
}
function getList() {
    $.get("./fakead", function(data) {
        mapping=JSON.parse(data);
        mapSync();
        popUpToast("Successfully refreshed latest list.");
    });
}
function modify(oldk, k, v) {
    $.post("./fakead",
    {
        oldk: oldk,
        k: k,
        v: v
    }, function() {
        getList();
    });
}
$(document).ready(function() {
    mapSync();
    getList();
});
