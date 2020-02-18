var fetch = (function () {
    function randInt(n, m) {
        return Math.floor(Math.random() * (m - n + 1)) + n;
    }

    function refresh() {
        var then = this.then;
        var xhr = new XMLHttpRequest();
        xhr.addEventListener("load", function () {
            var resp = JSON.parse(this.responseText);

            resp.nodes.forEach(function (n) {
                n.x = randInt(0, 255);
                n.y = randInt(0, 255);
                n.label = n.id;
                n.size = 1;
                delete n.color;
            });

            // resp.edges.forEach(function (e) {
            //     s.addEdge(e);
            // });

            resp.links = resp.edges;

            then(resp);
        });
        xhr.open("GET", "sigma");
        xhr.send();
    }

    var module = {
        every: 1000,
        then: function () {},
        refresh: refresh,
        loop: function () {
            module.refresh();
            setTimeout(module.loop, module.every);
        }
    };

    return module;
})();
