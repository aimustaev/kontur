import React, { useState, useCallback, useEffect } from "react";
import {
  useNodesState,
  useEdgesState,
  Node,
  Edge,
  Connection,
  addEdge,
} from "reactflow";
import { WorkflowCanvas } from "./components/WorkflowCanvas";
import { ComponentLibrary } from "./components/library";
import { ConfigSelector } from "./components/ConfigSelector";
import "./styles.css";
import { Button } from "antd";
import { transformMockToWorkflow } from "./utils";

interface WorkflowConfig {
  name: string;
  version: string;
  is_active: boolean;
  id: string;
}

interface WorkflowStep {
  id: string;
  name: string;
  type: "activity" | "signal" | "timer" | "decision";
  activityName?: string;
  input: (string | { value: any; isTemplate: boolean })[];
  output?: string;
  outputSchema?: any;
  timerDuration?: string;
  signalName?: string;
  actions?: WorkflowStep[];
  concurrent?: boolean;
}

interface Config {
  id: string;
  name: string;
  version: string;
  content: string;
  schema: { nodes: Node[]; edges: Edge[] };
  created_by: "aimustaev";
  is_active: boolean;
}

const WorkflowBuilder: React.FC = () => {
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [workflow, setWorkflow] = useState<Config | null>(null);

  const [selectedConfig, setSelectedConfig] = useState<WorkflowConfig | null>(
    null
  );

  useEffect(() => {
    const fetchConfigs = async () => {
      try {
        const response = await fetch(`/api-w/config/${selectedConfig?.id}`);
        if (!response.ok) {
          throw new Error("Failed to fetch configs");
        }
        const data = await response.json();
        setWorkflow(data?.[0]);
      } catch (error) {
        console.error("Error fetching configs:", error);
      }
    };

    fetchConfigs();
  }, [selectedConfig]);

  useEffect(() => {
    setNodes(
      workflow?.schema?.nodes.map((node) => {
        const data = node.data;
        return { ...node, data: { ...data, onInputChange, onOutputChange } };
      }) || []
    );
    setEdges(workflow?.schema?.edges || []);
  }, [workflow]);

  const onConnect = useCallback(
    (params: Connection) => {
      setEdges((eds) => addEdge(params, eds));
    },
    [setEdges]
  );

  const onNodeDragStop = useCallback(
    (_: React.MouseEvent, node: Node) => {
      setNodes((nds) =>
        nds.map((n) => {
          if (n.id === node.id) {
            return {
              ...n,
              position: node.position,
            };
          }
          return n;
        })
      );
    },
    [setNodes]
  );

  const onDragStart = useCallback((event: React.DragEvent, nodeData: any) => {
    event.dataTransfer.setData("application/reactflow/type", nodeData.type);
    event.dataTransfer.setData("application/reactflow/label", nodeData.label);
    event.dataTransfer.setData(
      "application/reactflow/description",
      nodeData.description
    );
    event.dataTransfer.setData(
      "application/reactflow/argIn",
      JSON.stringify(nodeData.argIn)
    );
    event.dataTransfer.setData(
      "application/reactflow/argOut",
      JSON.stringify(nodeData.argOut)
    );
    if (nodeData.activityName)
      event.dataTransfer.setData(
        "application/reactflow/activityName",
        nodeData.activityName
      );
    if (nodeData.timerDuration)
      event.dataTransfer.setData(
        "application/reactflow/timerDuration",
        nodeData.timerDuration
      );
    if (nodeData.signalName)
      event.dataTransfer.setData(
        "application/reactflow/signalName",
        nodeData.signalName
      );
    event.dataTransfer.effectAllowed = "move";
  }, []);

  useEffect(() => {
    console.log(nodes);
  }, [nodes]);

  const onInputChange = (nodeId: string, paramName: string, value: string) => {
    setNodes((nds) =>
      nds.map((node) =>
        node.id === nodeId
          ? {
              ...node,
              data: {
                ...node.data,
                input: { ...node.data.input, [paramName]: value },
                onInputChange,
              },
            }
          : node
      )
    );
  };

  const onOutputChange = (nodeId: string, paramName: string, value: string) => {
    setNodes((nds) =>
      nds.map((node) =>
        node.id === nodeId
          ? {
              ...node,
              data: {
                ...node.data,
                output: { ...node.data.output, [paramName]: value },
                onOutputChange,
              },
            }
          : node
      )
    );
  };

  const onDrop = useCallback(
    (event: React.DragEvent) => {
      event.preventDefault();
      const reactFlowBounds = (
        event.target as HTMLDivElement
      ).getBoundingClientRect();

      const type = event.dataTransfer.getData("application/reactflow/type");
      const label = event.dataTransfer.getData("application/reactflow/label");
      const description = event.dataTransfer.getData(
        "application/reactflow/description"
      );
      const argIn = JSON.parse(
        event.dataTransfer.getData("application/reactflow/argIn") || "[]"
      );
      const argOut = JSON.parse(
        event.dataTransfer.getData("application/reactflow/argOut") || "[]"
      );
      const activityName = event.dataTransfer.getData(
        "application/reactflow/activityName"
      );
      const timerDuration = event.dataTransfer.getData(
        "application/reactflow/timerDuration"
      );
      const signalName = event.dataTransfer.getData(
        "application/reactflow/signalName"
      );

      if (!type) return;

      const position = {
        x: event.clientX - reactFlowBounds.left,
        y: event.clientY - reactFlowBounds.top,
      };

      const newNode: WorkflowStep = {
        id: `step-${Date.now()}`,
        name: label,
        type: type as WorkflowStep["type"],
        input: [],
      };

      if (type === "activity") {
        newNode.activityName = activityName;
        newNode.output = "";
        newNode.outputSchema = { type: "object" };
      }
      if (type === "signal") {
        newNode.signalName = signalName;
      }
      if (type === "timer") {
        newNode.timerDuration = timerDuration || "5s";
      }

      setNodes((nds) => [
        ...nds,
        {
          id: newNode.id,
          type,
          data: {
            label: newNode.name,
            activityName,
            timerDuration,
            signalName,
            argIn,
            argOut,
            description,
            input: {},
            onInputChange,
            onOutputChange,
            output: {},
          },
          position,
        },
      ]);
    },
    [setNodes]
  );

  const onDragOver = useCallback((event: React.DragEvent) => {
    event.preventDefault();
    event.dataTransfer.dropEffect = "move";
  }, []);

  const handleConfigSelect = (config: WorkflowConfig | null) => {
    if (config) {
      setSelectedConfig(config);
    }
  };

  const handleSave = () => {
    const { states } = transformMockToWorkflow({ nodes, edges });
    const version = (Number(workflow?.version) + 0.1).toFixed(1).toString();
    const name = workflow?.name || "New Workflow";
    const content = JSON.stringify({ states, version, name });

    const newConfig: Config = {
      name: workflow?.name ?? "config1",
      version: version,
      content: content,
      created_by: "aimustaev",
      is_active: true,
      id: workflow?.id ?? "1",
      schema: { nodes, edges },
    };

    const fetchConfigs = async () => {
      try {
        const response = await fetch(`/api-w/config`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(newConfig),
        });
        if (!response.ok) {
          throw new Error("Failed to fetch configs");
        }
        const data = await response.json();
        console.log(data);
      } catch (error) {
        console.error("Error fetching configs:", error);
      }
    };

    fetchConfigs();
    // const workflow = {
    //   nodes: nodes.map((node) => ({
    //     id: node.id,
    //     type: node.type,
    //     data: node.data,
    //     position: node.position,
    //   })),
    //   edges: edges.map((edge) => ({
    //     id: edge.id,
    //     source: edge.source,
    //     target: edge.target,
    //   })),
    // };
    // console.log(workflow);
  };

  return (
    <div className="workflow-builder" style={{ position: "relative" }}>
      <div
        style={{
          padding: "16px",
          position: "absolute",
          zIndex: 1000,
          gap: 10,
          display: "flex",
          flexDirection: "row",
          alignItems: "center",
        }}
      >
        <ConfigSelector onConfigSelect={handleConfigSelect} />
        <Button type="primary" onClick={handleSave}>
          Сохранить
        </Button>
      </div>
      <WorkflowCanvas
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        onNodeDragStop={onNodeDragStop}
        onDrop={onDrop}
        onDragOver={onDragOver}
      />
      <ComponentLibrary onDragStart={onDragStart} />
    </div>
  );
};

export default WorkflowBuilder;
