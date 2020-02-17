(function () {

    var graph = {
            nodes: [],
            edges: [],
    };

    function refresh() {
        var xhr = new XMLHttpRequest();
        xhr.addEventListener("load", function () {
            var resp = JSON.parse(this.responseText);
            graph = resp;
        });
        xhr.open("GET", "sigma");
        xhr.send();
    }

    function refreshLp() {
        refresh();
        setTimeout(200, refreshLp);
    }

    s = new sigma({
            graph: graph,
            container: 'graph-container',
            settings: {
            drawEdges: false
            }
    });

    // Start the ForceAtlas2 algorithm:
    s.startForceAtlas2({worker: true, barnesHutOptimize: false});

    refreshLp();
})();
