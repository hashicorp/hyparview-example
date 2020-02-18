(function (d3, fetch) {
    // https://observablehq.com/@d3/force-directed-graph
    // https://medium.com/ninjaconcept/interactive-dynamic-force-directed-graphs-with-d3-da720c6d7811
    function chart(config, data) {
        const width = config.width;
        const height = config.height;

        var links = data.links.map(d => Object.create(d));
        var nodes = data.nodes.map(d => Object.create(d));

        const simulation = d3.forceSimulation(nodes)
              .force("link", d3.forceLink(links).id(d => d.id))
              .force("charge", d3.forceManyBody())
              .force("center", d3.forceCenter(width / 2, height / 2));

        const svg = d3.create("svg")
              .attr("viewBox", [0, 0, width, height]);

        const link = svg.append("g")
              .attr("stroke", "#999")
              .attr("stroke-opacity", 0.6)
              .selectAll("line")
              .data(links)
              .join("line")
              .attr("stroke-width", d => Math.sqrt(d.value));

        const node = svg.append("g")
              .attr("stroke", "#fff")
              .attr("stroke-width", 1.5)
              .selectAll("circle")
              .data(nodes)
              .join("circle")
              .attr("r", 5)
              .attr("fill", color);
              // .call(drag(simulation));

        node.append("title")
            .text(d => d.id);

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

    var data = {
        nodes: [],
        links: []
    };

    var svg = chart({height: 600, width: 800}, data);

    fetch.every = 5000;
    fetch.then = function (resp) {
        resp.nodes.forEach(function (n) {
            data.nodes.push(n);
        });
        resp.edges.forEach(function (n) {
            data.links.push(n);
        });
    };

    document.getElementById("container").appendChild(svg);
    fetch.loop();

})(d3, fetch);
