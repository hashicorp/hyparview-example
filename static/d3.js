(function (d3, fetch) {
    // ==================================================
    // https://medium.com/ninjaconcept/interactive-dynamic-force-directed-graphs-with-d3-da720c6d7811
    // https://github.com/ninjaconcept/d3-force-directed-graph/blob/master/example/4-dynamic-updates.html

    const width = window.innerWidth - 20;
    const height = window.innerHeight - 20;

    const svg = d3.select('svg');
    svg.attr('width', width).attr('height', height);

    const simulation = d3.forceSimulation()
          .force("link", d3.forceLink().id(linkKey))
          .force('charge', d3.forceManyBody())
          .force('x', d3.forceX())
          .force('y', d3.forceY())
          .force('center', d3.forceCenter(width / 2, height / 2));

    var data = {
        nodes: [],
        nodeSet: {},
        links: [],
        linkSet: {},
    };

    refresh();
    refreshTask(5000);

    function updateData(resp) {
        var add = {};

        // add nodes
        for (var k in resp.nodes) {
            var node = resp.nodes[k],
                old = data.nodeSet[k];

            if (nodeEq(node, old)) {
                continue;
            }

            if (old) {
                node.x = old.x;
                node.y = old.y;
            } else {
                node.x = randInt(0, width);
                node.y = randInt(0, height);
            }

            add[k] = node;
            data.nodes.push(node);
            data.nodeSet[k] = node;
        }

        data.nodes.filter(function (n) {
            // node not in resp, it's actually deleted
            if (! resp.nodes[n.id]) {
                delete data.nodeSet[n.id];
                return false;
            }
            // keep only unchanged or the copy that was in the resp
            var a = add[n.id];
            return !a || nodeEq(a, n);
        });

        for (k in resp.links) {
            if (data.linkSet[k]) continue;
            var link = resp.links[k];
            link.source = data.nodeSet[link.source];
            link.target = data.nodeSet[link.target];
            link.value = 2;
            data.links.push(link);
            data.linkSet[k] = link;
        }

        data.links.filter(function (l) {
            var k = linkKey(l);
            if (! resp.links[k]) {
                delete data.linkSet[k];
                return false;
            }
            return true;
        });

        function nodeEq(n, o) {
            return n && o &&
                n.id == o.id &&
                n.app == o.app;
        }
    }

    function updateGraph(resp) {
        updateData(resp);
        const ns = updateNodes(svg, data.nodes);
        const ls = updateLinks(svg, data.links);
        updateSimulation(simulation, ns, ls);
    }

    function updateSimulation(sim, node, link) {
        sim.on("tick", () => {
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
        let ns = svg.selectAll("circle")
            .data(nodes, nodeKey)
            .join("circle");
        ns.exit().remove();

        const enter = ns.enter()
              .append("circle")
              .attr("stroke", "#fff")
              .attr("stroke-width", 1.5)
              .attr("r", 6)
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
        updateGraph(resp);
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
