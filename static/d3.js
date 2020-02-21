(function (d3, fetch) {
    // ==================================================
    // https://medium.com/ninjaconcept/interactive-dynamic-force-directed-graphs-with-d3-da720c6d7811
    // https://github.com/ninjaconcept/d3-force-directed-graph/blob/master/example/4-dynamic-updates.html

    const width = window.innerWidth;
    const height = window.innerHeight;

    const svg = d3.select('svg');
    svg.attr('width', width).attr('height', height);

    const simulation = d3.forceSimulation()
          .force("link", d3.forceLink().id(linkKey))
          .force('charge', d3.forceManyBody().strength(-20))
          .force('center', d3.forceCenter(width / 2, height / 2));

    var data = {
        nodes: [],
        nodeSet: {},
        links: [],
        linkSet: {},
    };

    refreshTask(5000);

    function updateGraph(svg, sim, data, resp) {
        // links point to node objects
        var key, x, del;
        var add = setwiseAdds(data.nodeSet, nodeKey, resp.nodes);
        for (key in add) {
            x = add[key];
            x.x = randInt(width);
            x.y = randInt(height);
            data.nodes.push(x);
            data.nodeSet[nodeKey(x)] = x;
        }
        setwiseDel(data.nodes, data.nodeSet, nodeKey, add);

        add = setwiseAdds(data.linkSet, linkKey, resp.links);
        for (key in add) {
            x = add[key];
            x.source = data.nodeSet[x.source];
            x.target = data.nodeSet[x.target];
        }

        const ns = updateNodes(svg, data.nodes);
        const ls = updateLinks(svg, data.links);
        updateSimulation(sim, ns, ls);
    }

    function updateSimulation(sim, node, link) {
        simulation.on("tick", () => {
            link
                .attr("x1", d => d.source.x)
                .attr("y1", d => d.source.y)
                .attr("x2", d => d.target.x)
                .attr("y2", d => d.target.y);

            node
                .attr("cx", d => d.x)
                .attr("cy", d => d.y);
        });
        sim.restart();
    }

    function updateNodes(svg, nodes) {
        let ns = svg.selectAll("circle").data(nodes, nodeKey);
        ns.exit().remove();

        const enter = ns.enter()
              .append("circle")
              .attr("stroke", "#fff")
              .attr("stroke-width", 1.5)
              .attr("r", 5)
              .attr("fill", "#000");
        return ns;
    }

    function updateLinks(svg, links) {
        let ls = svg.selectAll("line").data(links, linkKey);
        ls.exit().remove();

        const enter = ls.enter()
              .append("line")
              .attr("stroke", "#999")
              .attr("stroke-opacity", 0.6);
        ls = enter.merge(ls);
        return ls;
    }

    // ==================================================
    // Fetch a bunch of stuff

    function randInt(n, m) {
        return Math.floor(Math.random() * (m - n + 1)) + n;
    }

    function handle(text) {
        var resp = JSON.parse(text);
        updateGraph(svg, simulation, data, resp);
    }

    function refresh() {
        var xhr = new XMLHttpRequest();
        xhr.addEventListener("load", function () {
            handle(this.responseText);
        });
        xhr.open("GET", "d3");
        xhr.send();
    }

    function refreshTask(n) {
        refresh();
        setTimeout(function () {refreshTask(n);}, n);
    }

    // ==================================================
    // Only add new items, remove missing ones

    function setwiseAdds(set, keyf, xs) {
        var xss = {};
        // add new elements, populate xss
        xs.forEach(function (x) {
            var k = keyf(x);
            xss[k] = true;
            if (!set[k]) {
                set[k] = true;
            }
        });

        return xss;
    }

    function setwiseDel(data, set, keyf, xss) {
        var del = {};

        for (let k in set) {
            if (!xss[k]) {
                del[k] = true;
            }
        }

        data.filter(function (d) {
            return del[keyf(d)];
        });
    }

    function nodeKey(n) {
        return n.id;
    }

    function linkKey(l) {
        var s = (typeof l.source == "string") ? l.source : l.source.id,
            t = (typeof l.target == "string") ? l.target : l.target.id;
        return s + t;
    }
})(d3, fetch);
