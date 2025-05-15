export function transformMockToWorkflow(mockWorkflow: any) {
    // First, create a map of all nodes for easy lookup
    const nodesMap = new Map(mockWorkflow.nodes.map((node: any) => [node.id, node]));

    // Find the first node (node that has no incoming edges)
    const firstNodeId = mockWorkflow.nodes.find((node: any) =>
        !mockWorkflow.edges.some((edge: any) => edge.target === node.id)
    )?.id;

    if (!firstNodeId) {
        throw new Error("Could not find starting node");
    }

    // Build ordered array of nodes following the edges
    const orderedNodes: any[] = [];
    let currentNodeId = firstNodeId;

    while (currentNodeId) {
        const currentNode = nodesMap.get(currentNodeId);
        if (!currentNode) break;

        orderedNodes.push(currentNode);

        // Find the next node by looking for an edge where current node is the source
        const nextEdge = mockWorkflow.edges.find((edge: any) => edge.source === currentNodeId);
        currentNodeId = nextEdge?.target;
    }

    // Transform the ordered nodes into states
    const states = orderedNodes.map((node: any) => {
        const baseState = {
            name: node.data.label,
            type: node.type,
            activityName: node.data.activityName || "",
            input: node.data.input ? Object.values(node.data.input) : [],
            outputSchema: {
                type: "object"
            }
        };

        // Add specific properties based on node type
        switch (node.type) {
            case "signal":
                return {
                    "name": "MessageListener",
                    "type": "signal",
                    "concurrent": true,
                    "signalName": "NewMessage",
                    "actions": [
                        {
                            "type": "activity",
                            "input": [
                                "$.signalPayload",
                                "$.ticket.Id"
                            ],
                            "activityName": "AddMassageToTicketActivity",
                            "outputSchema": {
                                "type": "object"
                            }
                        }
                    ],
                }
            case "timer":
                return {
                    name: `WaitForResponse ${new Date().valueOf()}`,
                    type: "timer",
                    timerDuration: node.data.input.timerDuration
                };
            case "activity":
                return {
                    ...baseState,
                    type: "activity",
                    output: node.data.argOut?.[0] || undefined
                };
            default:
                return baseState;
        }
    });

    return { states };
}
