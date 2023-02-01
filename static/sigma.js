/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

(function () {
    function createSigma(graph) {
        return new sigma({
            graph: graph,
            container: 'graph-container',
            settings: {
                drawEdges: false
            }
        });
    }

    var s;
    var graph = {
        nodes: [],
        edges: [],
    };
    s = createSigma(graph);
    // s.startForceAtlas2({worker: true, barnesHutOptimize: false});
})();
