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

    fetch.every = 5000;
    fetch.then = function (resp) {
        updateGraph(svg, simulation, data, resp);
    };

    fetch.loop();

    function updateGraph(svg, sim, data, resp) {
        setwiseAdd(data.nodes, data.nodeSet, nodeKey, resp.nodes);
        setwiseAdd(data.links, data.linkSet, linkKey, resp.links);
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
        const ng = svg.append("g").attr("class", "node");
        let ns = ng.selectAll("circle").data(nodes, nodeKey);
        ns.exit().remove();
        const enter = ns.enter().append()
              .join("circle")
              .attr("r", 5)
              .attr("fill", color);
        return ns;
    }

    function updateLinks(svg, links) {
        const lg = svg.append("g").attr("class", "link");
        let ls = lg.selectAll("line").data(links, linkKey);
        ls.exit().remove();
        const enter = ls.enter().append()
              .attr("stroke", "#999")
              .attr("stroke-opacity", 0.6);
        ls = enter.merge(ls);
        return ls;
    }

    // https://observablehq.com/@d3/force-directed-graph
    // https://medium.com/ninjaconcept/interactive-dynamic-force-directed-graphs-with-d3-da720c6d7811
    function chart(config, data) {
        const width = config.width;
        const height = config.height;

        var simulation = d3.forceSimulation(data.nodes)
            .force("link", d3.forceLink(data.links).id(d => d.id))
            .force("charge", d3.forceManyBody())
            .force("center", d3.forceCenter(width / 2, height / 2));

        var svg = d3.create("svg")
            .attr("viewBox", [0, 0, width, height]);

        var link = svg.append("g")
            .attr("stroke", "#999")
            .attr("stroke-opacity", 0.6)
            .selectAll("line")
            .data(data.links)
            .enter().append('line')
            .join("line")
            .attr("stroke-width", d => Math.sqrt(d.value));
        link.exit().remove();

        var node = svg.append("g")
            .attr("stroke", "#fff")
            .attr("stroke-width", 1.5)
            .selectAll("circle")
            .data(data.nodes)
            .enter().append("circle")
            .join("circle")
            .attr("r", 5)
            .attr("fill", color);
        node.exit().remove();
        // node.append("title").text(d => d.id);

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

        // invalidation.then(() => simulation.stop());
        return svg.node();
    }

    function color(node) {
        const scale = d3.scaleOrdinal(d3.schemeCategory10);
        return d => scale(d.group);
    }

    // ==================================================
    // Only add new items, remove missing ones

    function setwiseAdd(data, set, keyf, xs) {
        var xss = {},
            del = {};

        // add new elements, populate xss
        xs.forEach(function (x) {
            var k = keyf(x);
            xss[k] = true;
            if (!set[k]) {
                data.push(x);
                set[k] = true;
            }
        });

        // use xss to find delete keys
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

    function linkKey(x) {
        return x.source + x.target;
    }
})(d3, fetch);
