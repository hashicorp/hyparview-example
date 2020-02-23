const initData = {
    nodes: [ {id: 0 } ],
    links: []
};

const elem = document.getElementById("graph");

const Graph = ForceGraph()(elem)
      .onNodeHover(node => elem.style.cursor = node ? 'pointer' : null)
      .nodeAutoColorBy('app')
      .d3AlphaDecay(0)
      .d3VelocityDecay(0.08)
      .cooldownTime(60000)
      .linkColor(() => 'rgba(0,0,0,0.05)')
      .zoom(0.05)
      .enablePointerInteraction(false)
      .graphData(initData);

setInterval(refresh, 3000);

function refresh() {
    var xhr = new XMLHttpRequest();
    xhr.addEventListener("load", () => {
        var resp = JSON.parse(xhr.responseText);
        updateGraph(resp);
    });
    xhr.open("GET", "d3");
    xhr.send();
}

function updateGraph(resp) {
    const data = mergeData(Graph.graphData(), resp);
    Graph.graphData(data);
}

// create new arrays of nodes, links. We're using object identity to
// see of a thing has changed, so we want to preserve the old object
// between refreshes
function mergeData(data, resp) {
    var node_ids = {};
    var nodes = data.nodes.map(n => {
        if (! resp.nodes[n.id]) return null; // missing
        node_ids[n.id] = true;               // mark the idp
        if (nodeEq(n, resp.nodes[n.id])) return n; // keep old obj
        return resp.nodes[n.id];
    })
        .filter(n => { return n; });

    for (var k in resp.nodes) {
        if (node_ids[k]) continue;
        nodes.push(resp.nodes[k]);
    }

    var link_ids = {};
    const links = data.links.map(l => {
        const id = l.source + l.target;
        if (! resp.links[id]) return null; // missing
        link_ids[id] = true;
        return l;               // updates don't matter
    }).filter(l => { return l; });

    for (k in resp.links) {
        if (link_ids[k]) continue;
        links.push(resp.links[k]);
    }

    function nodeEq(n, o) {
        return n && o &&
            n.id == o.id &&
            n.app == o.app;
    }

    return { nodes, links };
}
