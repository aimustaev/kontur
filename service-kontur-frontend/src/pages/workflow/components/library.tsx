import React from 'react';
import { Typography } from 'antd';
import { ACTIVITY_NODES, SIGNAL_NODES, TIMER_NODES } from '../config';

interface ComponentLibraryProps {
  onDragStart: (event: React.DragEvent, nodeData: any) => void;
}

export const ComponentLibrary: React.FC<ComponentLibraryProps> = ({ onDragStart }) => {
  const renderNodeGroup = (title: string, nodes: any[]) => (
    <div style={{ marginTop: 18, width: "100%" }}>
      <Typography.Text strong>{title}</Typography.Text>
      {nodes.map((node) => (
        <div
          key={node.label}
          className="draggable-node"
          draggable
          onDragStart={(e) => onDragStart(e, node)}
          style={{
            border: "1px solid #d9d9d9",
            borderRadius: 6,
            padding: 12,
            background: "#fff",
            cursor: "grab",
            marginBottom: 4,
            boxShadow: "0 1px 2px rgba(0,0,0,0.04)",
          }}
        >
          <Typography.Text>{node.label}</Typography.Text>
          <div style={{ fontSize: 12, color: "#888" }}>
            {node.description}
          </div>
        </div>
      ))}
    </div>
  );

  return (
    <div className="side-panel">
      <Typography.Title level={5} style={{ marginTop: 16 }}>
        Библиотека шагов
      </Typography.Title>
      {renderNodeGroup("Активити", ACTIVITY_NODES)}
      {renderNodeGroup("Таймер", TIMER_NODES)}
      {renderNodeGroup("Сигнал", SIGNAL_NODES)}
    </div>
  );
}; 